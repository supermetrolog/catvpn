package main

import (
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/client"
	_ "github.com/supermetrolog/myvpn/internal/logger"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/tunnel"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
	"strconv"
	"time"
)

func main() {
	// TODO: config
	clientIp := net.IPv4(192, 168, 16, 2)
	clientPort := 7070

	serverIp := net.IPv4(192, 168, 16, 3)
	serverPort := 9090

	serverAddr, err := net.ResolveUDPAddr("udp", serverIp.String()+":"+strconv.Itoa(serverPort))

	checkErr("Unable resolve server udp addr", err)

	clientAddr, err := net.ResolveUDPAddr("udp", clientIp.String()+":"+strconv.Itoa(clientPort))

	checkErr("Unable resolve client udp addr", err)

	cfg := client.NewConfig(
		2000,
		60,
		serverAddr,
		clientAddr,
		1500,
		"udp",
		time.Second*5,
	)

	tunFactory := tuntap.New()
	tunnelFactory := tunnel.NewClientTunnelFactory()
	trafficRouteConfigurator := routeconfigurator.NewClientTrafficRouteConfigurator()

	s := client.NewClient(cfg, tunnelFactory, tunFactory, trafficRouteConfigurator)

	s.Serve()
}

func checkErr(message string, e error) {
	if e != nil {
		logrus.Fatalf(message, e)
	}
}
