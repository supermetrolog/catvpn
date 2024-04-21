package main

import (
	"github.com/supermetrolog/myvpn/internal/client"
	"github.com/supermetrolog/myvpn/internal/helpers/checkerr"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/tunnel"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
	"strconv"
)

func main() {
	serverIp := net.IPv4(192, 168, 16, 3)
	serverPort := 9090

	addr, err := net.ResolveUDPAddr("udp", serverIp.String()+":"+strconv.Itoa(serverPort))

	checkerr.CheckErr("Unable resolve udp addr", err)

	cfg := client.NewConfig(
		2000,
		60,
		addr,
		1300,
		"udp",
	)

	tunFactory := tuntap.New()
	tunnelFactory := tunnel.NewClientTunnelFactory()
	trafficRouteConfigurator := routeconfigurator.NewClientTrafficRouteConfigurator()

	s := client.NewClient(cfg, tunnelFactory, tunFactory, trafficRouteConfigurator)

	s.Serve()
}
