package client

import (
	"net"
)

type Config struct {
	BufferSize            int      // Кол-во байт, которые будут читаться из интерфейса и тунеля
	HeartBeatTimeInterval int      // Интервал времени, с которым нужно отпрввлять heartbeats для поддеркжи соединения
	ServerAddr            net.Addr // Адресс сервера
	Net                   string   // Протокол тоннеля: udp|tcp
	Mtu                   int
}

func NewConfig(
	bufferSize int,
	heartBeatTimeInterval int,
	serverAddr net.Addr,
	mtu int,
	net string,
) *Config {
	return &Config{
		BufferSize:            bufferSize,
		HeartBeatTimeInterval: heartBeatTimeInterval,
		ServerAddr:            serverAddr,
		Mtu:                   mtu,
		Net:                   net,
	}
}

func (c *Config) TunnelAddr() net.Addr {
	return c.ServerAddr
}
