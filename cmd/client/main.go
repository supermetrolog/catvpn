package main

import (
	"github.com/supermetrolog/myvpn/internal/client"
	"github.com/supermetrolog/myvpn/internal/ipdistributor"
	"github.com/supermetrolog/myvpn/internal/peersmanager"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/server"
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
	trafficRouteConfigurator := routeconfigurator.New()
	tunnelFactory := tunnel.NewTunnelFactory()
	ipDistributorFactory := ipdistributor.NewIpDistributorFactory()
	peersManager := peersmanager.New()

	s := server.NewServer(cfg, tunnelFactory, tunFactory, ipDistributorFactory, peersManager, trafficRouteConfigurator)

	s.Serve()
}
