package client

import (
	"github.com/supermetrolog/myvpn/internal/common"
	"io"
	"net"
)

type Net io.ReadWriter

type TunFactory interface {
	Create(subnet net.IPNet, mtu int) (common.Tun, error)
}

type TrafficRoutingConfigurator interface {
	RouteToIface(ifaceName string) error
}
