package client

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/common"
	"github.com/supermetrolog/myvpn/internal/helpers/checkerr"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"log"
	"net"
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
	checkerr.CheckErr("setup error", c.setup())

	go func() {
		checkerr.CheckErr("listen tunnel error", c.listenTunnel())
	}()

	go func() {
		checkerr.CheckErr("consume tunnel error", c.fromTunnelConsumer())
	}()

	<-c.connectedChan // TODOl

	go func() {
		checkerr.CheckErr("listen net error", c.listenNet())
	}()

	checkerr.CheckErr("consume net error", c.fromNetConsumer())
}

func (c *Client) setup() error {
	tunnel, err := c.tunnelFactory.Create(c.cfg.TunnelClientAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	c.tunnel = tunnel

	n, err := c.WriteToTunnel(protocol.NewTunnelPacket(c.cfg.TunnelServerAddr(), protocol.NewHeader(protocol.FlagAcknowledge), []byte{}))

	if err != nil {
		return fmt.Errorf("send ack flag to server error: %w", err)
	}

	log.Printf("Send Ack flag to server. Writed %d bytes", n)

	return nil
}

func (c *Client) listenNet() error {
	for {
		buf := make([]byte, c.cfg.BufferSize)
		n, err := c.net.Read(buf)
		if err != nil {
			return fmt.Errorf("read from net error: %w", err)
		}

		log.Printf("Readed bytes from NET %d", n)

		command.WritePacket(buf[:n])

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

		log.Printf("Readed bytes from TUNNEL: %d", n)

		c.fromTunnel <- protocol.UnmarshalTunnelPacket(addr, buf)
	}
}

func (c *Client) fromTunnelConsumer() error {
	for packet := range c.fromTunnel {
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			err := c.ackHandler(packet)
			if err != nil {
				return fmt.Errorf("flag ACK error: %w", err) // TODO
			}
		case protocol.FlagFin:
			err := c.finHandler()
			if err != nil {
				return fmt.Errorf("flag FIN error: %w", err) // TODO
			}
		case protocol.FlagData:
			err := c.dataHandler(packet)
			if err != nil {
				return fmt.Errorf("flag DATA error: %w", err) // TODO
			}
		}
	}

	return nil
}

func (c *Client) fromNetConsumer() error {
	for packet := range c.fromNet {
		tunnelPacket := protocol.NewTunnelPacket(c.cfg.TunnelServerAddr(), protocol.NewHeader(protocol.FlagData), *packet)

		n, err := c.tunnel.WriteTo(tunnelPacket.Packet().Marshal(), tunnelPacket.Addr())

		if err != nil {
			return fmt.Errorf("write in tunnel error: %w", err)
		}

		log.Printf("Wrote to TUNNEL %d bytes", n)
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
