package server

import (
	"github.com/supermetrolog/myvpn/internal/helpers/addr"
	"net"
)

type Config struct {
	BufferSize            int       // Кол-во байт, которые будут читаться из интерфейса и тунеля
	Subnet                net.IPNet // 10.1.1.0/24 Подсеть, в котором будет создан интерфейс, а так же отданы айпишники клиентам
	HeartBeatTimeInterval int       // Интервал времени, с которым нужно отпрввлять heartbeats для поддеркжи соединения
	ServerIP              net.IP    // Адрес который будет слушать UDP сервер
	ServerPort            int       // Порт, который будет слушать UDP сервер
	Net                   string
	Mtu                   int
}

func NewConfig(
	bufferSize int,
	subnet net.IPNet,
	heartBeatTimeInterval int,
	ServerIP net.IP,
	serverPort int,
	mtu int,
	net string,
) *Config {
	return &Config{
		BufferSize:            bufferSize,
		Subnet:                subnet,
		HeartBeatTimeInterval: heartBeatTimeInterval,
		ServerIP:              ServerIP,
		ServerPort:            serverPort,
		Mtu:                   mtu,
		Net:                   net,
	}
}

func (c *Config) TunnelAddr() net.Addr {
	return addr.NewForTransportLayer(c.Net, c.ServerIP, c.ServerPort)
}
