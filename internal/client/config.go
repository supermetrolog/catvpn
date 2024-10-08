package client

import (
	"net"
	"time"
)

type Config struct {
	BufferSize              int           // Кол-во байт, которые будут читаться из интерфейса и тунеля
	HeartBeatTimeInterval   int           // Интервал времени, с которым нужно отпрввлять heartbeats для поддеркжи соединения
	ServerAddr              net.Addr      // Адресс сервера
	ClientAddr              net.Addr      // Адресс клиента
	Net                     string        // Протокол тоннеля: udp|tcp
	Mtu                     int           // MTU
	ServerConnectionTimeout time.Duration // Время ожидания соединения с сервером
}

func NewConfig(
	bufferSize int,
	heartBeatTimeInterval int,
	serverAddr net.Addr,
	clientAddr net.Addr,
	mtu int,
	net string,
	serverConnectionTimeout time.Duration,
) *Config {
	return &Config{
		BufferSize:              bufferSize,
		HeartBeatTimeInterval:   heartBeatTimeInterval,
		ServerAddr:              serverAddr,
		ClientAddr:              clientAddr,
		Mtu:                     mtu,
		Net:                     net,
		ServerConnectionTimeout: serverConnectionTimeout,
	}
}

func (c *Config) TunnelServerAddr() net.Addr {
	return c.ServerAddr
}

func (c *Config) TunnelClientAddr() net.Addr {
	return c.ClientAddr
}
