package server

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/common"
	"github.com/supermetrolog/myvpn/internal/helpers/checkerr"
	"github.com/supermetrolog/myvpn/internal/helpers/command"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"golang.org/x/net/ipv4"
	"log"
)

type Server struct {
	cfg        *Config
	fromTunnel chan *protocol.TunnelPacket
	fromNet    chan *protocol.NetPacket
	tunnel     common.Tunnel
	net        Net

	tunnelFactory              common.TunnelFactory
	tunFactory                 TunFactory
	ipDistributorFactory       IpDistributorFactory
	ipDistributor              IpDistributor
	peersManager               PeersManager
	trafficRoutingConfigurator TrafficRoutingConfigurator
}

func NewServer(
	cfg *Config,
	tunnelF common.TunnelFactory,
	tunF TunFactory,
	ipDistributorFactory IpDistributorFactory,
	peersManager PeersManager,
	trafficRoutingConfigurator TrafficRoutingConfigurator,
) *Server {
	return &Server{
		cfg:                        cfg,
		fromTunnel:                 make(chan *protocol.TunnelPacket),
		fromNet:                    make(chan *protocol.NetPacket),
		tunnelFactory:              tunnelF,
		tunFactory:                 tunF,
		ipDistributorFactory:       ipDistributorFactory,
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

	checkerr.CheckErr("consume net error", s.fromNetConsumer())
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

	log.Printf("Created tunnel. Tunnel addr: %s", tunnel.LocalAddr())

	s.tunnel = tunnel

	err = s.trafficRoutingConfigurator.RouteToSubnet(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("traffic route to subnet error: %w", err)
	}

	ipDistributor, err := s.ipDistributorFactory.Create(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("create ip distributor error: %w", err)
	}

	s.ipDistributor = ipDistributor

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

		s.fromTunnel <- protocol.UnmarshalTunnelPacket(addr, buf[:n])
	}
}

func (s *Server) fromTunnelConsumer() error {
	log.Println("Consume tunnel")
	for packet := range s.fromTunnel {
		log.Printf("Readed from tunnel channel. Flag: %d", packet.Packet().Header().Flag())
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			err := s.ackHandler(packet)
			if err != nil {
				return fmt.Errorf("flag ACK error: %w", err) // TODO
			}
		case protocol.FlagFin:
			err := s.finHandler(packet)
			if err != nil {
				return fmt.Errorf("flag FIN error: %w", err) // TODO
			}

		case protocol.FlagData:
			err := s.dataHandler(packet)
			if err != nil {
				return fmt.Errorf("flag DATA error: %w", err) // TODO
			}

		default:
			return fmt.Errorf("unknown flag")
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

		n, err := s.WriteToTunnel(protocol.NewTunnelPacket(peer.Addr(), protocol.NewHeader(protocol.FlagData), *packet))

		if err != nil {
			return fmt.Errorf("write in tunnel error: %w", err)
		}

		log.Printf("Wrote to TUNNEL %d bytes", n)
	}

	return nil
}

func (s *Server) WriteToTunnel(packet *protocol.TunnelPacket) (int, error) {
	n, err := s.tunnel.WriteTo(packet.Packet().Marshal(), packet.Addr())

	if err != nil {
		return n, fmt.Errorf("write to tunnel error: %w", err)
	}

	return n, err
}
