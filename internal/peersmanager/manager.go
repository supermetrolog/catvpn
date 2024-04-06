package peersmanager

import (
	"fmt"
	"github.com/supermetrolog/myvpn/internal/protocol"
	"net"
)

type PeersManager struct {
	peers []*protocol.Peer
}

func New() *PeersManager {
	return &PeersManager{}
}

func (pm *PeersManager) Add(peer *protocol.Peer) error {
	pm.peers = append(pm.peers, peer)

	return nil
}

func (pm *PeersManager) Remove(peer *protocol.Peer) error {
	for i, p := range pm.peers {
		if peer == p {
			pm.peers = append(pm.peers[i:], pm.peers[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("peer not found")
}

func (pm *PeersManager) FindByDedicatedIp(ip net.IP) (peer *protocol.Peer, exists bool, err error) {
	for _, peer := range pm.peers {
		if peer.DedicatedIP().Equal(ip) {
			return peer, true, nil
		}
	}

	return nil, false, nil
}

func (pm *PeersManager) FindByRealIp(ip net.IP) (peer *protocol.Peer, exists bool, err error) {
	for _, peer := range pm.peers {
		if peer.RealIP().Equal(ip) {
			return peer, true, nil
		}
	}

	return nil, false, nil
}
