package client

import (
	"io"
	"net"
)

type Tunnel net.PacketConn
type Net io.ReadWriter

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type Tun interface {
	io.ReadWriteCloser
	Name() string
}

type TunFactory interface {
	Create(subnet net.IPNet, mtu int) (Tun, error)
}

type TrafficRoutingConfigurator interface {
	RouteToIface(ifaceName string) error
}
