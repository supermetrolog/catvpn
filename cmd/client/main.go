package main

import (
	"github.com/supermetrolog/myvpn/internal/client"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/tunnel"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
)

func main() {

	serverIp := net.IPv4(10, 1, 1, 1)

	cfg := client.NewConfig(
		2000,
		60,
		serverIp,
		9090,
		1300,
		"udp",
	)

	tunFactory := tuntap.New()
	tunnelFactory := tunnel.NewTunnelFactory()
	trafficRouteConfigurator := routeconfigurator.NewClientTrafficRouteConfigurator()

	s := client.NewClient(cfg, tunnelFactory, tunFactory, trafficRouteConfigurator)

	s.Serve()
}
