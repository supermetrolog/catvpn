package server

import (
	"net"
)

type Config struct {
	BufferSize            int       // Кол-во байт, которые будут читаться из интерфейса и тунеля
	Subnet                net.IPNet // 10.1.1.0/24 Подсеть, в котором будет создан интерфейс, а так же отданы айпишники клиентам
	HeartBeatTimeInterval int       // Интервал времени, с которым нужно отпрввлять heartbeats для поддеркжи соединения
	ServerAddr            net.Addr  // Порт, который будет слушать UDP сервер
	Net                   string
	Mtu                   int
}

func NewConfig(
	bufferSize int,
	subnet net.IPNet,
	heartBeatTimeInterval int,
	serverAddr net.Addr,
	mtu int,
	net string,
) *Config {
	return &Config{
		BufferSize:            bufferSize,
		Subnet:                subnet,
		HeartBeatTimeInterval: heartBeatTimeInterval,
		ServerAddr:            serverAddr,
		Mtu:                   mtu,
		Net:                   net,
	}
}

func (c *Config) TunnelAddr() net.Addr {
	return c.ServerAddr
}
