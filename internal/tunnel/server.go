package tunnel

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/common"
	"net"
)

type ServerTunnelFactory struct {
}

func NewServerTunnelFactory() *ServerTunnelFactory {
	return &ServerTunnelFactory{}
}

func (t *ServerTunnelFactory) Create(addr net.Addr) (common.Tunnel, error) {
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
