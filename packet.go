package main

import (
	"encoding/binary"
	"net"
	"unsafe"

	"tomasweigenast.com/forwarding_tables/packets"
)

// packet represents an IP packet data
type packet struct {
	sender      net.IP
	destination net.IP
	protocol    packets.IPv4Protocol
	data        []byte
	len         int64
}

const ipv4_packet_base_size = 4 + 4 + 1 + 8

func new_packet(sender net.IP, destination net.IP, payload packets.IPv4Payload) *packet {
	sender = sender.To4()
	destination = destination.To4()

	if sender == nil {
		panic("IPv6 is not supported")
	}

	if destination == nil {
		panic("IPv6 is not supported")
	}

	return &packet{
		sender:      sender,
		destination: destination,
		protocol:    payload.Protocol(),
		data:        payload.Encode(),
		len:         payload.Len(),
	}
}

func parse_packet(buffer []byte) *packet {
	sender_ip := net.IP(buffer[0:4])
	dest_ip := net.IP(buffer[4:8])
	protocol := buffer[8]
	dataLen := int64(binary.LittleEndian.Uint64(buffer[9:17]))
	data := buffer[17:]
	return &packet{
		sender:      sender_ip,
		destination: dest_ip,
		protocol:    packets.IPv4Protocol(protocol),
		len:         dataLen,
		data:        data,
	}
}

// encode encodes the packet to a sequence of bytes
func (p *packet) encode() []byte {

	lenBytes := *(*[8]byte)(unsafe.Pointer(&p.len))

	// 4 bytes per ip + 8 bytes of data len + data length
	buffer := make([]byte, ipv4_packet_base_size+p.len)

	// copy data
	copy(buffer[0:], p.sender)
	copy(buffer[4:], p.destination)
	buffer[8] = byte(p.protocol)
	copy(buffer[9:], lenBytes[:])
	copy(buffer[17:], p.data)

	return buffer
}
