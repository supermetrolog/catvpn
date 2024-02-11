package main

import (
	"github.com/supermetrolog/myvpn/internal/helpers/checkerr"
	"github.com/supermetrolog/myvpn/internal/server"
	"net"
)

type TrafficRoutingConfigurator interface {
	RouteToSubnet(subnet net.IPNet) error
}

func main() {
	subnet := net.IPNet{
		IP:   net.IPv4(10, 1, 1, 1),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}

	cfg := &server.Config{
		BufferSize:            2000,
		Subnet:                subnet,
		HeartBeatTimeInterval: 60,
		ServerPort:            9090,
	}

	iface := createTun(subnet)
	defer iface.Close() // TODO

	configureOSTrafficRouting(subnet)

	server := server.NewServer(cfg, iface, iface)
}

func configureOSTrafficRouting(subnet net.IPNet) {
	var trafficRoutingConfigurator TrafficRoutingConfigurator

	err := trafficRoutingConfigurator.RouteToSubnet(subnet)
	checkerr.CheckErr("route traffic to subnet error", err)
}

func createTun(subnet net.IPNet) TunIface {
	var tunCreator TunCreator

	iface, err := tunCreator.Create(subnet)
	checkerr.CheckErr("create tun iface error", err)

	return iface
}
