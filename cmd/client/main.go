package main

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/supermetrolog/myvpn/internal/client"
	"github.com/supermetrolog/myvpn/internal/config"
	_ "github.com/supermetrolog/myvpn/internal/logger"
	"github.com/supermetrolog/myvpn/internal/routeconfigurator"
	"github.com/supermetrolog/myvpn/internal/tunnel/udp"
	"github.com/supermetrolog/myvpn/internal/tuntap"
	"net"
	"strconv"
	"time"
)

func main() {
	clientCfg, err := getClientConfig()
	checkErr("Unable get client config", err)

	serverAddr, err := net.ResolveUDPAddr("udp", clientCfg.ServerHost.Ip+":"+strconv.Itoa(int(clientCfg.ServerHost.Port)))

	checkErr("Unable resolve server udp addr", err)

	clientAddr, err := net.ResolveUDPAddr("udp", clientCfg.ClientHost.Ip+":"+strconv.Itoa(int(clientCfg.ClientHost.Port)))

	checkErr("Unable resolve client udp addr", err)

	cfg := client.NewConfig(
		int(clientCfg.BufferSize),
		int(clientCfg.HeartBeatTimeInterval),
		serverAddr,
		clientAddr,
		int(clientCfg.MTU),
		"udp",
		time.Second*time.Duration(clientCfg.ServerConnectionTimeout),
	)

	tunFactory := tuntap.New()
	tunnelFactory := udp.NewClientTunnelFactory()
	trafficRouteConfigurator := routeconfigurator.NewClientTrafficRouteConfigurator()

	s := client.NewClient(cfg, tunnelFactory, tunFactory, trafficRouteConfigurator)

	s.Serve()
}

func getClientConfig() (*config.ClientConfig, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("unable load env from .env file: %w", err)
	}

	var cfg config.ClientConfig

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
