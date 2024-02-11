package tunnel

import (
	"github.com/supermetrolog/myvpn/internal/server"
	"net"
)

type TunnelFactory struct {
}

func NewTunnelFactory() *TunnelFactory {
	return &TunnelFactory{}
}

func (t *TunnelFactory) Create(addr net.Addr) (server.Tunnel, error) {

}
