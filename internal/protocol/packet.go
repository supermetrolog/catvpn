package protocol

import (
	"bytes"
	"log"
	"net"
)

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

func (p *Packet) Marshal() []byte {
	var buf []byte

	buffer := bytes.NewBuffer(buf)

	buffer.WriteByte(byte(p.Header().Flag()))
	buffer.Write(p.Payload())

	log.Println(buffer.Bytes())

	return buffer.Bytes()
}

func UnmarshalTunnelPacket(addr net.Addr, bytes []byte) *TunnelPacket {
	return NewTunnelPacket(addr, NewHeader(Flag(bytes[0])), bytes[HeaderSize:])
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
	log.Printf("Create new tunnel packet: %s", addr.String())
	return &TunnelPacket{
		packet: NewPacket(header, payload),
		addr:   addr,
	}
}

type NetPacket []byte
