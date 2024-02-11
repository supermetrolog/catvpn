package server

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"io"
	"log"
	"net"
)

type Tunnel net.PacketConn

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type TunFactory interface {
	Create(subnet net.IPNet) (io.ReadWriteCloser, error)
}

type Net io.ReadWriter

type Server struct {
	cfg   *Config
	peers []*protocol.Peer

	toTunnel chan *protocol.Packet
	toNet    chan *protocol.Packet

	fromTunnel chan *protocol.TunnelPacket
	fromNet    chan *protocol.NetPacket

	tunnelFactory TunnelFactory
	tunFactory    TunFactory

	tunnel Tunnel
	net    Net
}

func NewServer(cfg *Config, tunnelF TunnelFactory, tunF TunFactory) *Server {
	return &Server{
		cfg:           cfg,
		tunnelFactory: tunnelF,
		tunFactory:    tunF,
		toTunnel:      make(chan *protocol.Packet),
		toNet:         make(chan *protocol.Packet),
		peers:         make([]*protocol.Peer, 0),
	}
}

func (s *Server) Serve() error {
	if err := s.setup(); err != nil {
		return err
	}

	go s.listenTunnel()
	go s.listenNet()
}

func (s *Server) setup() error {
	tun, err := s.tunFactory.Create(s.cfg.Subnet)
	if err != nil {
		return fmt.Errorf("create tun iface error: %w", err)
	}
	s.net = tun

	tunnel, err := s.tunnelFactory.Create(s.cfg.TunnelAddr())
	if err != nil {
		return fmt.Errorf("create tunnel error: %w", err)
	}

	s.tunnel = tunnel

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

		s.fromTunnel <- protocol.NewTunnelPacket(addr, protocol.Flag(buf[0]), buf[1:])
	}
}

func (s *Server) toNetConsumer() error {
	for packet := range s.toNet {
		n, err := s.net.Write(packet.Payload())
		if err != nil {
			return fmt.Errorf("write to net error: %w", err)
		}

		log.Printf("Write bytes to net: %d", n)
	}

	return nil
}

func (s *Server) toTunnelConsumer() error {
	// TODO

	return nil
}

func (s *Server) fromTunnelConsumer() error {
	for packet := range s.fromTunnel {
		switch packet.Packet().Flag() {
		case protocol.FlagAcknowledge:
			// TODO: create tunnel (save real ip map to assigned dedicated ip)
		case protocol.FlagData:
			s.toNet <- packet.Packet()

		}
	}

	return nil
}

func (s *Server) fromNetConsumer() error {
	for packet := range s.fromNet {
		// TODO: get dst ip, get tunnel from tunnel map, produce to tunnel chan TunnelPacket
	}

	return nil
}
