package tunnel

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/common"
	"net"
	"time"
)

type ClientTunnelFactory struct {
}

func NewClientTunnelFactory() *ClientTunnelFactory {
	return &ClientTunnelFactory{}
}

type udpConnDecorator struct {
	conn *net.UDPConn
}

// WriteTo не поддерживается при UDP соединении
func (t *udpConnDecorator) WriteTo(p []byte, _ net.Addr) (n int, err error) {
	return t.conn.Write(p)
}

// ReadFrom не поддерживается при UDP соединении
func (t *udpConnDecorator) ReadFrom(p []byte) (int, net.Addr, error) {
	n, err := t.conn.Read(p)

	if err != nil {
		return n, t.conn.RemoteAddr(), fmt.Errorf("unable read from tunnel: %w", err)
	}

	return n, t.conn.RemoteAddr(), nil
}
func (t *udpConnDecorator) Close() error {
	return t.conn.Close()
}
func (t *udpConnDecorator) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}
func (t *udpConnDecorator) SetDeadline(tm time.Time) error {
	return t.conn.SetDeadline(tm)
}
func (t *udpConnDecorator) SetReadDeadline(tm time.Time) error {
	return t.conn.SetReadDeadline(tm)
}
func (t *udpConnDecorator) SetWriteDeadline(tm time.Time) error {
	return t.conn.SetWriteDeadline(tm)
}

func (t *ClientTunnelFactory) Create(addr net.Addr) (common.Tunnel, error) {
	udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
	if err != nil {
		return nil, fmt.Errorf("unable to resolve udp addr: %w", err)
	}

	conn, err := net.DialUDP(udpAddr.Network(), nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("unable to listen udp: %w", err)
	}

	return &udpConnDecorator{conn: conn}, nil
}
