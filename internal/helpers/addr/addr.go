package addr

import (
	"net"
	"strconv"
)

type Addr struct {
	net     string
	address string
}

func (a *Addr) Network() string {
	return a.net
}

func (a *Addr) String() string {
	return a.address
}

func NewForTransportLayer(network string, ip net.IP, port int) *Addr {
	return &Addr{
		net:     network,
		address: ip.String() + ":" + strconv.Itoa(port),
	}
}
