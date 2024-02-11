package peersmanager

import (
	"github.com/supermetrolog/myvpn/internal/protocol"
	"net"
)

type PeersManager struct {
}

func New() *PeersManager {
	return &PeersManager{}
}

func (p *PeersManager) Add(peer *protocol.Peer) error {
	return nil
}
func (p *PeersManager) Remove(peer *protocol.Peer) error {
	return nil
}
func (p *PeersManager) FindByDedicatedIp(ip net.IP) (peer *protocol.Peer, exists bool, err error) {
	return nil
}
func (p *PeersManager) FindByRealIp(ip net.IP) (peer *protocol.Peer, exists bool, err error) {
	return nil
}
