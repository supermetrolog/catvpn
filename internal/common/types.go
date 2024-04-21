package common

import (
	"io"
	"net"
)

type Tunnel net.PacketConn

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type Tun interface {
	io.ReadWriteCloser
	Name() string
}
