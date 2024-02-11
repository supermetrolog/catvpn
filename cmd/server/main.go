package main

import (
	"github.com/supermetrolog/myvpn/internal/server"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
)

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

	tunFactory := tuntap.NewFactory()

	server := server.NewServer(cfg)
}
