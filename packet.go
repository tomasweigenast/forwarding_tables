package main

import (
	"encoding/binary"
	"net"
	"unsafe"
)

// packet represents an IP packet data
type packet struct {
	sender      net.IP
	destination net.IP
	data        []byte
	len         int64
}

func new_packet(sender net.IP, destination net.IP, data []byte) *packet {
	sender = sender.To4()
	destination = destination.To4()

	if sender == nil {
		panic("IPv6 is not supported")
	}

	if destination == nil {
		panic("IPv6 is not supported")
	}

	return &packet{
		sender,
		destination,
		data,
		int64(len(data)),
	}
}

func parse_packet(buffer []byte) *packet {
	sender_ip := net.IP(buffer[0:4])
	dest_ip := net.IP(buffer[4:8])
	dataLen := int64(binary.LittleEndian.Uint64(buffer[8:16]))
	data := buffer[16:]
	return &packet{
		sender:      sender_ip,
		destination: dest_ip,
		len:         dataLen,
		data:        data,
	}
}

// encode encodes the packet to a sequence of bytes
func (p *packet) encode() []byte {

	lenBytes := *(*[8]byte)(unsafe.Pointer(&p.len))

	// 4 bytes per ip + 8 bytes of data len + data length
	buffer := make([]byte, 16+p.len)

	// copy data
	copy(buffer[0:], p.sender)
	copy(buffer[4:], p.destination)
	copy(buffer[8:], lenBytes[:])
	copy(buffer[16:], p.data)

	return buffer
}
