package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/common"
	"github.com/supermetrolog/myvpn/internal/helpers/ippacket"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"golang.org/x/net/ipv4"
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
	err := s.setup()

	if err != nil {
		logrus.Fatalf("Setup error: %v", err)
	}

	go func() {
		err = s.listenTunnel()

		if err != nil {
			logrus.Fatalf("Listen tunnel error: %v", err)
		}
	}()
	go func() {
		err = s.listenNet()

		if err != nil {
			logrus.Fatalf("Listen net error: %v", err)
		}
	}()

	go func() {
		err = s.fromTunnelConsumer()

		if err != nil {
			logrus.Fatalf("Consume tunnel error: %v", err)
		}
	}()

	err = s.fromNetConsumer()

	if err != nil {
		logrus.Fatalf("Consume net error: %v", err)
	}
}

func (s *Server) setup() error {
	tun, err := s.tunFactory.Create(s.cfg.Subnet, s.cfg.Mtu)

	if err != nil {
		return fmt.Errorf("create tun iface error: %w", err)
	}

	s.net = tun
	logrus.Debugf("Created tun interface: %s", tun.Name())

	tunnel, err := s.tunnelFactory.Create(s.cfg.TunnelAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	s.tunnel = tunnel
	logrus.Debugf("Created tunnel. Tunnel addr: %s", tunnel.LocalAddr())

	err = s.trafficRoutingConfigurator.RouteToSubnet(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("traffic route to subnet error: %w", err)
	}
	logrus.Debug("Configure traffic routing")

	ipDistributor, err := s.ipDistributorFactory.Create(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("create ip distributor error: %w", err)
	}

	s.ipDistributor = ipDistributor
	logrus.Debug("Create ip distributor")

	return nil
}

func (s *Server) listenNet() error {
	for {
		buf := make([]byte, s.cfg.BufferSize)
		n, err := s.net.Read(buf)
		if err != nil {
			return fmt.Errorf("read from net error: %w", err)
		}

		logrus.Debugf("Readed bytes from NET %d", n)

		ippacket.LogHeader(buf[:n])

		p := protocol.NetPacket(buf[:n])
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

		logrus.Debugf("Readed bytes from TUNNEL: %d", n)

		s.fromTunnel <- protocol.UnmarshalTunnelPacket(addr, buf[:n])
	}
}

func (s *Server) fromTunnelConsumer() error {
	for packet := range s.fromTunnel {
		logrus.Debugf("Readed from tunnel channel. Flag: %d", packet.Packet().Header().Flag())
		switch packet.Packet().Header().Flag() {
		case protocol.FlagAcknowledge:
			err := s.ackHandler(packet)
			if err != nil {
				logrus.Errorf("Flag ACK error: %v", err)
			}
		case protocol.FlagFin:
			err := s.finHandler(packet)
			if err != nil {
				logrus.Errorf("Flag FIN error: %v", err)
			}
		case protocol.FlagData:
			err := s.dataHandler(packet)
			if err != nil {
				logrus.Errorf("Flag DATA error: %v", err)
			}

		default:
			logrus.Warnf("Unknown flag: %b", packet.Packet().Header().Flag())
		}
	}

	return nil
}

func (s *Server) fromNetConsumer() error {
	for packet := range s.fromNet {
		header, err := ipv4.ParseHeader(*packet)
		if err != nil {
			logrus.Warnf("Parse from net ip header error: %v", err)
			continue
		}

		peer, exists, err := s.peersManager.FindByDedicatedIp(header.Dst)

		if err != nil {
			logrus.Warnf("Find by dedicated port error: %v", err)
			continue
		}

		if !exists {
			logrus.Warnf("Peer not found")
			continue
		}

		n, err := s.WriteToTunnel(protocol.NewTunnelPacket(peer.Addr(), protocol.NewHeader(protocol.FlagData), *packet))

		if err != nil {
			logrus.Warnf("Write in tunnel error: %v", err)
			continue
		}

		logrus.Infof("Write to TUNNEL %d bytes", n)
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
