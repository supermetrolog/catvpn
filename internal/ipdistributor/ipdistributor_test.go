package ipdistributor_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/supermetrolog/myvpn/internal/ipdistributor"
	"net"
	"testing"
)

func newIpDistributor(t *testing.T) *ipdistributor.IpDistributor {
	ipNet := net.IPNet{
		IP:   net.IPv4(10, 1, 1, 0),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}

	d, err := ipdistributor.New(ipNet)

	assert.NoError(t, err)

	return d
}

func TestCreateIPPool(t *testing.T) {
	d := newIpDistributor(t)

	assert.Len(t, d.GetIPPool(), 256)
	assert.Equal(t, expectedIpPool(), d.GetIPPool())
}

func expectedIpPool() map[string]bool {
	ipPool := make(map[string]bool, 256)

	for i := 0; i < 256; i++ {
		ipPool[net.IPv4(10, 1, 1, byte(i)).String()] = false
	}

	return ipPool
}

func TestAllocate(t *testing.T) {
	d := newIpDistributor(t)

	allocatedIp, err := d.AllocateIP()

	busyIpCount := 0

	assert.NoError(t, err)

	for ip, isBusy := range d.GetIPPool() {
		if isBusy {
			busyIpCount++

			assert.Equal(t, allocatedIp.String(), ip)
		}
	}

	assert.Equal(t, 1, busyIpCount)
}

func TestAllocateWithEndedError(t *testing.T) {
	d := newIpDistributor(t)

	for i := 0; i < 256; i++ {
		_, err := d.AllocateIP()
		assert.NoError(t, err)
	}

	_, err := d.AllocateIP()
	assert.Error(t, err)
}

func TestRelease(t *testing.T) {
	d := newIpDistributor(t)

	allocatedIp1, err := d.AllocateIP()
	assert.NoError(t, err)
	allocatedIp2, err := d.AllocateIP()
	assert.NoError(t, err)
	allocatedIp3, err := d.AllocateIP()
	assert.NoError(t, err)

	err = d.ReleaseIP(allocatedIp1)
	assert.NoError(t, err)
	err = d.ReleaseIP(allocatedIp1)
	assert.Error(t, err)

	err = d.ReleaseIP(allocatedIp2)
	assert.NoError(t, err)

	err = d.ReleaseIP(allocatedIp3)
	assert.NoError(t, err)
	err = d.ReleaseIP(allocatedIp3)
	assert.Error(t, err)
}
