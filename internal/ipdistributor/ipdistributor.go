package ipdistributor

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/server"
	"net"
	"sync"
)

type Factory struct {
}

func NewIpDistributorFactory() *Factory {
	return &Factory{}
}

func (i Factory) Create(ipNet net.IPNet) (server.IpDistributor, error) {
	return New(ipNet)
}

type IpDistributor struct {
	subnet net.IPNet
	m      sync.Mutex
	ipPool map[string]bool
}

func New(ipNet net.IPNet) (*IpDistributor, error) {
	d := &IpDistributor{
		subnet: ipNet,
		ipPool: make(map[string]bool, 256),
	}

	d.generateIpPool()

	return d, nil
}

func (ipd *IpDistributor) AllocateIP() (net.IP, error) {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	for ip, isBusy := range ipd.ipPool {
		if !isBusy {
			ipd.ipPool[ip] = true
			return net.ParseIP(ip), nil
		}
	}

	return nil, fmt.Errorf("ip pool has ended")
}

func (ipd *IpDistributor) ReleaseIP(ip net.IP) error {
	ipd.m.Lock()
	defer ipd.m.Unlock()

	for ipp, isBusy := range ipd.ipPool {
		if ipp == ip.String() {
			if !isBusy {
				return fmt.Errorf("ip is not busy")
			}

			ipd.ipPool[ipp] = false
			return nil
		}
	}

	return fmt.Errorf("ip not found")
}

func (ipd *IpDistributor) GetIPPool() map[string]bool {
	return ipd.ipPool
}

func (ipd *IpDistributor) generateIpPool() {
	for ip := ipd.subnet.IP.Mask(ipd.subnet.Mask); ipd.subnet.Contains(ip); ipd.incIP(ip) {
		ipd.ipPool[ip.String()] = false
	}
}

func (ipd *IpDistributor) incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
