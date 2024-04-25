package network

import (
	"fmt"
	"net"
)

func ResoleIpFromAddr(addr net.Addr) (net.IP, error) {
	host, _, err := net.SplitHostPort(addr.String())

	if err != nil {
		return nil, fmt.Errorf("make new peer error: split host port error: %w", err)
	}

	ip, err := net.ResolveIPAddr("ip", host)

	if err != nil {
		return nil, fmt.Errorf("resolve IP address error: %w", err)
	}

	return ip.IP, nil
}
