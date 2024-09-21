package common

import (
	"io"
	"net"
)

type Tunnel interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
	LocalAddr() net.Addr
	io.Closer
}

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type Tun interface {
	io.ReadWriteCloser
	Name() string
}
