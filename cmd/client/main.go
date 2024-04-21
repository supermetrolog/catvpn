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
	clientIp := net.IPv4(192, 168, 16, 2)
	clientPort := 7070

	serverIp := net.IPv4(192, 168, 16, 3)
	serverPort := 9090

	serverAddr, err := net.ResolveUDPAddr("udp", serverIp.String()+":"+strconv.Itoa(serverPort))

	checkerr.CheckErr("Unable resolve server udp addr", err)

	clientAddr, err := net.ResolveUDPAddr("udp", clientIp.String()+":"+strconv.Itoa(clientPort))

	checkerr.CheckErr("Unable resolve client udp addr", err)

	cfg := client.NewConfig(
		2000,
		60,
		serverAddr,
		clientAddr,
		1300,
		"udp",
	)

	tunFactory := tuntap.New()
	tunnelFactory := tunnel.NewClientTunnelFactory()
	trafficRouteConfigurator := routeconfigurator.NewClientTrafficRouteConfigurator()

	s := client.NewClient(cfg, tunnelFactory, tunFactory, trafficRouteConfigurator)

	s.Serve()
}
