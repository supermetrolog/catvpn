package tunnel

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/server"
	"net"
)

type TunnelFactory struct {
}

func NewTunnelFactory() *TunnelFactory {
	return &TunnelFactory{}
}

func (t *TunnelFactory) Create(addr net.Addr) (server.Tunnel, error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return nil, fmt.Errorf("unable to resolve udp addr: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to listen udp: %w", err)
	}

	return conn, nil
}
