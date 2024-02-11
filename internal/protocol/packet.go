package protocol

import "net"

const (
	FlagHandshake Flag = iota
	FlagAcknowledge
	FlagFin
	FlagData
)

type Flag byte

type Packet struct {
	flag    Flag
	payload []byte
}

func NewPacket(flag Flag, payload []byte) *Packet {
	return &Packet{
		flag:    flag,
		payload: payload,
	}
}

func (p *Packet) Payload() []byte {
	return p.payload
}

func (p *Packet) Flag() Flag {
	return p.flag
}

type TunnelPacket struct {
	packet *Packet
	addr   net.Addr
}

func (p *TunnelPacket) Packet() *Packet {
	return p.packet
}

func (p *TunnelPacket) Addr() net.Addr {
	return p.addr
}

func NewTunnelPacket(addr net.Addr, flag Flag, payload []byte) *TunnelPacket {
	return &TunnelPacket{
		packet: NewPacket(flag, payload),
		addr:   addr,
	}
}

type NetPacket []byte
