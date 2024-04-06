package server

import (
	"github.com/supermetrolog/myvpn/internal/protocol"
	"io"
	"net"
)

type TrafficRoutingConfigurator interface {
	RouteToSubnet(subnet net.IPNet) error
}

type Tunnel net.PacketConn

type TunnelFactory interface {
	Create(addr net.Addr) (Tunnel, error)
}

type Tun io.ReadWriteCloser

type TunFactory interface {
	Create(subnet net.IPNet, mtu int) (Tun, error)
}

type Net io.ReadWriter

type IpDistributor interface {
	AllocateIP() (net.IP, error)
	ReleaseIP(net.IP) error
}

type IpDistributorFactory interface {
	Create(ipNet net.IPNet) (IpDistributor, error)
}

type PeersManager interface {
	Add(peer *protocol.Peer) error
	Remove(peer *protocol.Peer) error
	FindByDedicatedIp(ip net.IP) (peer *protocol.Peer, exists bool, err error)
	FindByRealIp(ip net.IP) (peer *protocol.Peer, exists bool, err error)
}
