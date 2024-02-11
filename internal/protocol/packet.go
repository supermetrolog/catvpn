package protocol

import "net"

const (
	FlagHandshake Flag = iota
	FlagAcknowledge
	FlagFin
	FlagData
)

const HeaderSize = 1

type Flag byte

type Header struct {
	flag Flag
}

func NewHeader(flag Flag) Header {
	return Header{
		flag: flag,
	}
}

func (h Header) Flag() Flag {
	return h.flag
}

type Packet struct {
	header  Header
	payload []byte
}

func NewPacket(header Header, payload []byte) *Packet {
	return &Packet{
		header:  header,
		payload: payload,
	}
}

func (p *Packet) Payload() []byte {
	return p.payload
}

func (p *Packet) Header() Header {
	return p.header
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

func NewTunnelPacket(addr net.Addr, header Header, payload []byte) *TunnelPacket {
	return &TunnelPacket{
		packet: NewPacket(header, payload),
		addr:   addr,
	}
}

type NetPacket []byte
