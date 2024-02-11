package protocol

import "net"

type Peer struct {
	realIp      net.IP
	dedicatedIp net.IP
}
