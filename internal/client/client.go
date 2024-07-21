package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/common"
	"github.com/supermetrolog/myvpn/internal/helpers/ippacket"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"net"
	"time"
)

type State struct {
	isConnectedToServer bool
	allocatedIP         net.IP
}

type Client struct {
	cfg           *Config
	fromTunnel    chan *protocol.TunnelPacket
	fromNet       chan *protocol.NetPacket
	tunnel        common.Tunnel
	net           Net
	tun           common.Tun // TODO: saved for close conn
	connectedChan chan struct{}
	state         State

	tunnelFactory              common.TunnelFactory
	tunFactory                 TunFactory
	trafficRoutingConfigurator TrafficRoutingConfigurator
}

func NewClient(
	cfg *Config,
	tunnelF common.TunnelFactory,
	tunF TunFactory,
	trafficRoutingConfigurator TrafficRoutingConfigurator,
) *Client {
	return &Client{
		cfg:                        cfg,
		fromTunnel:                 make(chan *protocol.TunnelPacket),
		fromNet:                    make(chan *protocol.NetPacket),
		tunnelFactory:              tunnelF,
		tunFactory:                 tunF,
		trafficRoutingConfigurator: trafficRoutingConfigurator,
		connectedChan:              make(chan struct{}),
	}
}

func (c *Client) Serve() {
	err := c.setup()

	if err != nil {
		logrus.Fatalf("Setup error: %v", err)
	}

	go func() {
		err = c.listenTunnel()

		if err != nil {
			logrus.Fatalf("Listen tunnel error: %v", err)
		}
	}()

	go func() {
		err = c.fromTunnelConsumer()

		if err != nil {
			logrus.Fatalf("Consume tunnel error: %v", err)
		}
	}()

	logrus.Debugln("Waiting connect...")

	select {
	case <-c.connectedChan:
	case <-time.After(c.cfg.ServerConnectionTimeout):
		logrus.Fatalf("Conntection timeout")
	}

	go func() {
		err = c.listenNet()

		if err != nil {
			logrus.Fatalf("Listen net error: %v", err)
		}
	}()

	err = c.fromNetConsumer()

	if err != nil {
		logrus.Fatalf("Consume net error: %v", err)
	}
}

func (c *Client) setup() error {
	tunnel, err := c.tunnelFactory.Create(c.cfg.TunnelServerAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	c.tunnel = tunnel

	_, err = c.WriteToTunnel(protocol.NewTunnelPacket(
		c.cfg.TunnelServerAddr(),
		protocol.NewHeader(protocol.FlagAcknowledge),
		[]byte{},
	))

	if err != nil {
		return fmt.Errorf("send ack flag to server error: %w", err)
	}

	logrus.Debugln("Send ACK flag to server")

	return nil
}

func (c *Client) listenNet() error {
	for {
		buf := make([]byte, c.cfg.BufferSize)
		n, err := c.net.Read(buf)
		if err != nil {
			return fmt.Errorf("read from net error: %w", err)
		}

		logrus.Debugf("Readed bytes from NET %d", n)

		ippacket.LogHeader(buf[:n])

		p := protocol.NetPacket(buf[:n])
		c.fromNet <- &p
	}
}

func (c *Client) listenTunnel() error {
	for {
		buf := make([]byte, c.cfg.BufferSize)
		n, addr, err := c.tunnel.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("read from tunnel error: %w", err)
		}

		logrus.Debugf("Readed bytes from TUNNEL: %d", n)

		c.fromTunnel <- protocol.UnmarshalTunnelPacket(addr, buf[:n])
	}
}

func (c *Client) fromTunnelConsumer() error {
	for packet := range c.fromTunnel {
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			err := c.ackHandler(packet)
			if err != nil {
				logrus.Warnf("Flag ACK error: %v", err)
			}
		case protocol.FlagFin:
			err := c.finHandler()
			if err != nil {
				logrus.Warnf("Flag FIN error: %v", err)
			}
		case protocol.FlagData:
			err := c.dataHandler(packet)
			if err != nil {
				logrus.Warnf("Flag DATA error: %v", err)
			}
		}
	}

	return nil
}

func (c *Client) fromNetConsumer() error {
	for packet := range c.fromNet {
		tunnelPacket := protocol.NewTunnelPacket(c.cfg.TunnelServerAddr(), protocol.NewHeader(protocol.FlagData), *packet)

		n, err := c.WriteToTunnel(tunnelPacket)

		if err != nil {
			logrus.Warnf("Write to tunnel error: %v", err)
			continue
		}

		logrus.Debugf("Writed to TUNNEL %d bytes", n)
	}

	return nil
}

func (c *Client) WriteToTunnel(packet *protocol.TunnelPacket) (int, error) {
	n, err := c.tunnel.WriteTo(packet.Packet().Marshal(), packet.Addr())

	if err != nil {
		return n, fmt.Errorf("write to tunnel error: %w", err)
	}

	return n, err
}
