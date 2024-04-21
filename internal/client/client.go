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
	tunnel, err := c.tunnelFactory.Create(c.cfg.TunnelAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	c.tunnel = tunnel

	n, err := c.WriteToTunnel(protocol.NewTunnelPacket(c.cfg.TunnelAddr(), protocol.NewHeader(protocol.FlagAcknowledge), []byte{}))

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

		p := protocol.NetPacket(buf)
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

		command.WritePacket(buf[:n])

		c.fromTunnel <- protocol.UnmarshalTunnelPacket(addr, buf)
	}
}

func (c *Client) fromTunnelConsumer() error {
	for packet := range c.fromTunnel {
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			if c.state.isConnectedToServer {
				log.Println("Warn! Already connected to server")
				break
			}

			log.Printf("Ack FLAG. Connected to server")

			dedicatedIPBytes := packet.Packet().Payload()[protocol.HeaderSize : net.IPv4len+protocol.HeaderSize]
			c.state.allocatedIP = net.IPv4(dedicatedIPBytes[0], dedicatedIPBytes[1], dedicatedIPBytes[2], dedicatedIPBytes[3])
			log.Printf("Allocated IP: %s", c.state.allocatedIP.String())
			c.state.isConnectedToServer = true

			tun, err := c.tunFactory.Create(
				net.IPNet{
					IP:   c.state.allocatedIP,
					Mask: net.IPv4Mask(255, 255, 255, 0), // TODO:
				},
				c.cfg.Mtu,
			)

			if err != nil {
				return fmt.Errorf("create tun interface error: %w", err)
			}

			c.net = tun

			err = c.trafficRoutingConfigurator.RouteToIface(tun.Name()) // TODO: refactor
			if err != nil {
				return fmt.Errorf("traffic route to iface error: %w", err)
			}

			c.connectedChan <- struct{}{} // TODO

		case protocol.FlagFin:
			c.state.isConnectedToServer = false
			log.Printf("Fin FLAG. Disconnect from server")
			err := c.tunnel.Close()
			if err != nil {
				log.Printf("Error! Tunnel close error: %v", err)
			}
			err = c.tun.Close()

			if err != nil {
				log.Printf("Error! Tun close error: %v", err)
			}

			return nil // TODO: err handler
		case protocol.FlagData:
			// TODO: check is connected
			n, err := c.net.Write(packet.Packet().Payload())
			if err != nil {
				return fmt.Errorf("write to net error: %w", err)
			}

			log.Printf("Write bytes to net: %d", n)
		}
	}

	return nil
}

func (c *Client) fromNetConsumer() error {
	for packet := range c.fromNet {
		tunnelPacket := protocol.NewTunnelPacket(c.cfg.TunnelAddr(), protocol.NewHeader(protocol.FlagData), *packet)

		n, err := c.tunnel.WriteTo(*packet, tunnelPacket.Addr())

		if err != nil {
			return fmt.Errorf("write in tunnel error: %w", err)
		}

		log.Printf("Wrote to TUNNEL %d bytes", n)
	}

	return nil
}

func (c *Client) WriteToTunnel(packet *protocol.TunnelPacket) (int, error) {
	return c.tunnel.WriteTo(packet.Packet().Marshal(), packet.Addr())
}
