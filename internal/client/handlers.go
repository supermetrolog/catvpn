package client

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"log"
	"net"
)

func (c *Client) ackHandler(packet *protocol.TunnelPacket) error {
	if c.state.isConnectedToServer {
		log.Println("Warn! Already connected to server")
		return fmt.Errorf("already connected")
	}

	log.Printf("Ack FLAG. Connected to server")

	//dedicatedIPBytes := packet.Packet().Payload()[protocol.HeaderSize : net.IPv4len+protocol.HeaderSize]
	dedicatedIPBytes := packet.Packet().Payload()
	c.state.allocatedIP = net.IPv4(dedicatedIPBytes[0], dedicatedIPBytes[1], dedicatedIPBytes[2], dedicatedIPBytes[3])
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

	return nil
}

func (c *Client) finHandler() error {
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
}

func (c *Client) dataHandler(packet *protocol.TunnelPacket) error {
	// TODO: check is connected
	command.WritePacket(packet.Packet().Payload())
	n, err := c.net.Write(packet.Packet().Payload())
	if err != nil {
		return fmt.Errorf("write to net error: %w", err)
	}

	log.Printf("Write bytes to net: %d", n)

	return nil
}
