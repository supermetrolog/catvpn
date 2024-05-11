package client

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/helpers/ippacket"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"net"
)

func (c *Client) ackHandler(packet *protocol.TunnelPacket) error {
	if c.state.isConnectedToServer {
		return fmt.Errorf("already connected")
	}

	logrus.Debugf("Ack FLAG. Connected to server")

	dedicatedIPBytes := packet.Packet().Payload()
	c.state.allocatedIP = net.IPv4(dedicatedIPBytes[0], dedicatedIPBytes[1], dedicatedIPBytes[2], dedicatedIPBytes[3])
	c.state.allocatedIP = net.IPv4(dedicatedIPBytes[0], dedicatedIPBytes[1], dedicatedIPBytes[2], dedicatedIPBytes[3])

	logrus.Infof("Allocated IP: %s", c.state.allocatedIP.String())

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

	c.connectedChan <- struct{}{}

	logrus.Infof("Connected to server")

	return nil
}

func (c *Client) finHandler() error {
	c.state.isConnectedToServer = false

	logrus.Debugln("Fin FLAG. Disconnect from server")

	err := c.tunnel.Close()

	if err != nil {
		logrus.Errorf("Tunnel close error: %v", err)
	}

	err = c.tun.Close()

	if err != nil {
		logrus.Errorf("Tun close error: %v", err)
	}

	return nil
}

func (c *Client) dataHandler(packet *protocol.TunnelPacket) error {
	if !c.state.isConnectedToServer {
		return fmt.Errorf("client is not connected to server")
	}

	ippacket.LogHeader(packet.Packet().Payload())

	n, err := c.net.Write(packet.Packet().Payload())

	if err != nil {
		return fmt.Errorf("write to net error: %w", err)
	}

	logrus.Debugf("Write bytes to net: %d", n)

	return nil
}
