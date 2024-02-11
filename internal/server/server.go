package server

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/helpers/checkerr"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"golang.org/x/net/ipv4"
	"io"
	"log"
	"net"
)

type TrafficRoutingConfigurator interface {
	RouteToSubnet(subnet net.IPNet) error
}

type Tunnel net.PacketConn

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type Tun io.ReadWriteCloser

type TunFactory interface {
	Create(subnet net.IPNet, mtu int) (Tun, error)
}

type Net io.ReadWriter

type IpDistributor interface {
	AllocateIP() (net.IP, error)
	ReleaseIP(net.IP) error
	LoadAvailableSubnet(net.IPNet) error
}

type PeersManager interface {
	Add(peer *protocol.Peer) error
	Remove(peer *protocol.Peer) error
	FindByDedicatedIp(ip net.IP) (peer *protocol.Peer, exists bool, err error)
	FindByRealIp(ip net.IP) (peer *protocol.Peer, exists bool, err error)
}

type Server struct {
	cfg        *Config
	fromTunnel chan *protocol.TunnelPacket
	fromNet    chan *protocol.NetPacket
	tunnel     Tunnel
	net        Net

	tunnelFactory              TunnelFactory
	tunFactory                 TunFactory
	ipDistributor              IpDistributor
	peersManager               PeersManager
	trafficRoutingConfigurator TrafficRoutingConfigurator
}

func NewServer(
	cfg *Config,
	tunnelF TunnelFactory,
	tunF TunFactory,
	ipDistributor IpDistributor,
	peersManager PeersManager,
	trafficRoutingConfigurator TrafficRoutingConfigurator,
) *Server {
	return &Server{
		cfg:                        cfg,
		tunnelFactory:              tunnelF,
		tunFactory:                 tunF,
		ipDistributor:              ipDistributor,
		peersManager:               peersManager,
		trafficRoutingConfigurator: trafficRoutingConfigurator,
	}
}

func (s *Server) Serve() {
	checkerr.CheckErr("setup error", s.setup())

	go func() {
		checkerr.CheckErr("listen tunnel error", s.listenTunnel())
	}()
	go func() {
		checkerr.CheckErr("listen net error", s.listenNet())
	}()

	go func() {
		checkerr.CheckErr("consume tunnel error", s.fromTunnelConsumer())
	}()

	checkerr.CheckErr("consume net error", s.fromTunnelConsumer())
}

func (s *Server) setup() error {
	tun, err := s.tunFactory.Create(s.cfg.Subnet, s.cfg.Mtu)
	if err != nil {
		return fmt.Errorf("create tun iface error: %w", err)
	}
	s.net = tun

	tunnel, err := s.tunnelFactory.Create(s.cfg.TunnelAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	s.tunnel = tunnel

	err = s.trafficRoutingConfigurator.RouteToSubnet(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("traffic route to subnet error: %w", err)
	}

	return nil
}

func (s *Server) listenNet() error {
	for {
		buf := make([]byte, s.cfg.BufferSize)
		n, err := s.net.Read(buf)
		if err != nil {
			return fmt.Errorf("read from net error: %w", err)
		}

		log.Printf("Readed bytes from NET %d", n)

		command.WritePacket(buf[:n])

		p := protocol.NetPacket(buf)
		s.fromNet <- &p
	}
}

func (s *Server) listenTunnel() error {
	for {
		buf := make([]byte, s.cfg.BufferSize)
		n, addr, err := s.tunnel.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("read from tunnel error: %w", err)
		}

		log.Printf("Readed bytes from TUNNEL: %d", n)

		command.WritePacket(buf[:n])

		s.fromTunnel <- protocol.NewTunnelPacket(addr, protocol.NewHeader(protocol.Flag(buf[0])), buf[1:]) // TODO: from bytes to header struct
	}
}

func (s *Server) fromTunnelConsumer() error {
	for packet := range s.fromTunnel {
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			header, err := ipv4.ParseHeader(packet.Packet().Payload())
			if err != nil {
				return fmt.Errorf("parse from net ip header error: %w", err)
			}

			allocatedIp, err := s.ipDistributor.AllocateIP()
			if err != nil {
				return fmt.Errorf("allocate ip error: %w", err)
			}

			err = s.peersManager.Add(protocol.NewPeer(header.Src, allocatedIp, packet.Addr()))
			if err != nil {
				return fmt.Errorf("add new peer error: %w", err)
			}
		case protocol.FlagFin:
			header, err := ipv4.ParseHeader(packet.Packet().Payload())
			if err != nil {
				return fmt.Errorf("parse from net ip header error: %w", err)
			}

			peer, exists, err := s.peersManager.FindByDedicatedIp(header.Src)
			if err != nil {
				return fmt.Errorf("find by dedicated ip error: %w", err)
			}
			if !exists {
				return fmt.Errorf("peer with dedicated ip %s not found", header.Src.String())
			}

			err = s.peersManager.Remove(peer)
			if err != nil {
				return fmt.Errorf("remove peer error: %w", err)
			}
		case protocol.FlagData:
			n, err := s.net.Write(packet.Packet().Payload())
			if err != nil {
				return fmt.Errorf("write to net error: %w", err)
			}

			log.Printf("Write bytes to net: %d", n)
		}
	}

	return nil
}

func (s *Server) fromNetConsumer() error {
	for packet := range s.fromNet {
		header, err := ipv4.ParseHeader(*packet)
		if err != nil {
			return fmt.Errorf("parse from net ip header error: %w", err)
		}

		peer, exists, err := s.peersManager.FindByDedicatedIp(header.Dst)

		if err != nil {
			return fmt.Errorf("find by dedicated port error: %w", err)
		}

		if !exists {
			return fmt.Errorf("peer not found")
		}

		n, err := s.tunnel.WriteTo(*packet, peer.Addr())

		if err != nil {
			return fmt.Errorf("write in tunnel error: %w", err)
		}

		log.Printf("Wrote to TUNNEL %d bytes", n)
	}

	return nil
}
