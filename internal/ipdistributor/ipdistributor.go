package ipdistributor

import (
	"github.com/supermetrolog/myvpn/internal/server"
	"net"
)

type IpDistributorFactory struct {
}

func NewIpDistributorFactory() *IpDistributorFactory {
	return &IpDistributorFactory{}
}

func (i IpDistributorFactory) Create(ipNet net.IPNet) (server.IpDistributor, error) {
	return New(ipNet), nil
}

type IpDistributor struct {
	subnet net.IPNet
	ipPool []net.IP
}

func New(ipNet net.IPNet) *IpDistributor {

	return &IpDistributor{
		subnet: ipNet,
	}
}

func (i *IpDistributor) AllocateIP() (net.IP, error) {
	return nil, nil
}

func (i *IpDistributor) ReleaseIP(net.IP) error {
	return nil
}

func (i *IpDistributor) LoadAvailableSubnet(ipNet net.IPNet) error {
	return nil
}
