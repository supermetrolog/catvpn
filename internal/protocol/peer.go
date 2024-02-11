package protocol

import (
	"net"
)

type Peer struct {
	realIP      net.IP
	dedicatedIP net.IP
	addr        net.Addr
}

func (p *Peer) RealIP() net.IP {
	return p.realIP
}

func (p *Peer) DedicatedIP() net.IP {
	return p.dedicatedIP
}

func (p *Peer) Addr() net.Addr {
	return p.addr
}

func NewPeer(realIP, dedicatedIP net.IP, addr net.Addr) *Peer {
	return &Peer{
		realIP:      realIP,
		dedicatedIP: dedicatedIP,
		addr:        addr,
	}
}
