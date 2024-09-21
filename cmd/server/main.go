package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/config"
	"github.com/supermetrolog/myvpn/internal/ipdistributor"
	_ "github.com/supermetrolog/myvpn/internal/logger"
	"github.com/supermetrolog/myvpn/internal/peersmanager"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/server"
	"github.com/supermetrolog/myvpn/internal/tunnel/udp"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
	"strconv"
)

func main() {
	serverCfg, err := getClientConfig()
	checkErr("Unable get server config", err)

	_, subnet, err := net.ParseCIDR(serverCfg.Subnet)
	checkErr("Unable parse CIDR", err)

	addr, err := net.ResolveUDPAddr("udp", serverCfg.ServerHost.Ip+":"+strconv.Itoa(int(serverCfg.ServerHost.Port)))

	checkErr("Unable resolve udp addr", err)

	cfg := server.NewConfig(
		int(serverCfg.BufferSize),
		*subnet,
		int(serverCfg.HeartBeatTimeInterval),
		addr,
		int(serverCfg.MTU),
		"udp",
	)

	tunFactory := tuntap.New()
	trafficRouteConfigurator := routeconfigurator.NewServerTrafficRouteConfigurator()
	tunnelFactory := udp.NewServerTunnelFactory()
	ipDistributorFactory := ipdistributor.NewIpDistributorFactory()
	peersManager := peersmanager.New()

	s := server.NewServer(cfg, tunnelFactory, tunFactory, ipDistributorFactory, peersManager, trafficRouteConfigurator)

	s.Serve()
}

func getClientConfig() (*config.ServerConfig, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("unable load env from .env file: %w", err)
	}

	var cfg config.ServerConfig

	err = cleanenv.ReadEnv(&cfg)

	if err != nil {
		return nil, fmt.Errorf("unable map env to config struct: %w", err)
	}

	return &cfg, nil
}

func checkErr(message string, e error) {
	if e != nil {
		logrus.Fatalf(message, e)
	}
}
