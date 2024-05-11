package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/helpers/ippacket"
	"github.com/supermetrolog/myvpn/internal/helpers/network"
	"github.com/supermetrolog/myvpn/internal/protocol"
)

func (s *Server) ackHandler(packet *protocol.TunnelPacket) error {
	ip, err := network.ResoleIpFromAddr(packet.Addr())
	if err != nil {
		return fmt.Errorf("resolve IP address error: %w", err)
	}

	_, exists, err := s.peersManager.FindByRealIp(ip)
	if err != nil {
		return fmt.Errorf("find by real ip error: %w", err)
	}
	if exists {
		return fmt.Errorf("peer with real ip %s already exists", ip)
	}

	allocatedIp, err := s.ipDistributor.AllocateIP()
	if err != nil {
		return fmt.Errorf("allocate ip error: %w", err)
	}

	peer, err := protocol.NewPeer(allocatedIp, packet.Addr())
	if err != nil {
		return fmt.Errorf("flag ack: new peer error: %w", err)
	}

	logrus.Infof("Created new peer with addr: %s and dedicated ip: %s", peer.Addr(), peer.DedicatedIP())

	err = s.peersManager.Add(peer)
	if err != nil {
		return fmt.Errorf("add new peer error: %w", err)
	}

	_, err = s.WriteToTunnel(protocol.NewTunnelPacket(peer.Addr(), protocol.NewHeader(protocol.FlagAcknowledge), allocatedIp.To4()))

	if err != nil {
		return fmt.Errorf("write ack answer to peer error: %w", err)
	}

	return nil
}

func (s *Server) finHandler(packet *protocol.TunnelPacket) error {
	ip, err := network.ResoleIpFromAddr(packet.Addr())
	if err != nil {
		return fmt.Errorf("resolve IP address error: %w", err)
	}

	peer, exists, err := s.peersManager.FindByRealIp(ip)
	if err != nil {
		return fmt.Errorf("find by real ip error: %w", err)
	}
	if !exists {
		return fmt.Errorf("peer with real ip %s not found", ip)
	}

	err = s.peersManager.Remove(peer)
	if err != nil {
		return fmt.Errorf("remove peer error: %w", err)
	}

	logrus.Infof("Removed peer with addr: %s and dedicated ip: %s", peer.Addr(), peer.DedicatedIP())

	err = s.ipDistributor.ReleaseIP(peer.DedicatedIP())

	if err != nil {
		return fmt.Errorf("release peer dedicated IP error: %w", err)
	}

	return nil
}

func (s *Server) dataHandler(packet *protocol.TunnelPacket) error {
	ippacket.LogHeader(packet.Packet().Payload())

	n, err := s.net.Write(packet.Packet().Payload())
	if err != nil {
		return fmt.Errorf("write to net error: %w", err)
	}

	logrus.Debugf("Write bytes to NET: %d", n)

	return nil
}
