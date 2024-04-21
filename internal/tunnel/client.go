package tunnel

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/common"
	"net"
)

type ClientTunnelFactory struct {
}

func NewClientTunnelFactory() *ClientTunnelFactory {
	return &ClientTunnelFactory{}
}

//func (t *ClientTunnelFactory) Create(addr net.Addr) (common.Tunnel, error) {
//	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
//	//udpAddr, err := net.ResolveUDPAddr(addr.Network(), "server:9090")
//	if err != nil {
//		return nil, fmt.Errorf("unable to resolve udp addr: %w", err)
//	}
//
//	conn, err := net.DialUDP(udpAddr.Network(), nil, udpAddr)
//	if err != nil {
//		return nil, fmt.Errorf("unable to listen udp: %w", err)
//	}
//
//	n, err := conn.WriteTo([]byte("ANAL"), udpAddr)
//
//	if err != nil {
//		log.Fatalf("ANAL: %v", err)
//	}
//
//	log.Println(n)
//
//	return conn, nil
//}

func (t *ClientTunnelFactory) Create(addr net.Addr) (common.Tunnel, error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), "192.168.16.2:9091") // TODO:
	if err != nil {
		return nil, fmt.Errorf("unable to resolve udp addr: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to listen udp: %w", err)
	}

	return conn, nil
}
