package ipdistributor

import (
	"net"
)

type IpDistributor struct {
}

func New() *IpDistributor {
	return &IpDistributor{}
}

func (i *IpDistributor) AllocateIP() (net.IP, error) {

}

func (i *IpDistributor) ReleaseIP(net.IP) error {

}

func (i *IpDistributor) LoadAvailableSubnet(ipNet net.IPNet) error {

}
