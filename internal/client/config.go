package client

import (
	"github.com/supermetrolog/myvpn/internal/helpers/addr"
	"net"
)

type Config struct {
	BufferSize            int    // Кол-во байт, которые будут читаться из интерфейса и тунеля
	HeartBeatTimeInterval int    // Интервал времени, с которым нужно отпрввлять heartbeats для поддеркжи соединения
	ServerIp              net.IP // IP сервера
	ServerPort            int    // Порт, который будет слушать UDP сервер
	Net                   string // Протокол тоннеля: udp|tcp
	Mtu                   int
}

func NewConfig(
	bufferSize int,
	heartBeatTimeInterval int,
	serverIP net.IP,
	serverPort int,
	mtu int,
	net string,
) *Config {
	return &Config{
		BufferSize:            bufferSize,
		HeartBeatTimeInterval: heartBeatTimeInterval,
		ServerIp:              serverIP,
		ServerPort:            serverPort,
		Mtu:                   mtu,
		Net:                   net,
	}
}

func (c *Config) TunnelAddr() net.Addr {
	return addr.NewForTransportLayer(c.Net, c.ServerIp, c.ServerPort)
}
