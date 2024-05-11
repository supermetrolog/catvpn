package main

import (
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/ipdistributor"
	"github.com/supermetrolog/myvpn/internal/peersmanager"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/server"
	"github.com/supermetrolog/myvpn/internal/tunnel"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
	"strconv"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	subnet := net.IPNet{
		IP:   net.IPv4(10, 1, 1, 1),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}

	serverIp := net.IPv4(0, 0, 0, 0)
	serverPort := 9090

	addr, err := net.ResolveUDPAddr("udp", serverIp.String()+":"+strconv.Itoa(serverPort))

	checkErr("Unable resolve udp addr", err)

	cfg := server.NewConfig(
		2000,
		subnet,
		60,
		addr,
		1500,
		"udp",
	)

	tunFactory := tuntap.New()
	trafficRouteConfigurator := routeconfigurator.NewServerTrafficRouteConfigurator()
	tunnelFactory := tunnel.NewServerTunnelFactory()
	ipDistributorFactory := ipdistributor.NewIpDistributorFactory()
	peersManager := peersmanager.New()

	s := server.NewServer(cfg, tunnelFactory, tunFactory, ipDistributorFactory, peersManager, trafficRouteConfigurator)

	s.Serve()
}

func checkErr(message string, e error) {
	if e != nil {
		logrus.Fatalf(message, e)
	}
}
