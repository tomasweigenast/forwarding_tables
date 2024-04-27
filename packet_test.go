package main

import (
	"net"
	"reflect"
	"testing"
)

func TestPacketEncode(t *testing.T) {
	data := []byte("hello world")
	packet := new_packet(net.ParseIP("10.4.1.10"), net.ParseIP("13.231.12.4"), data)
	buffer := packet.encode()

	parsed_packet := parse_packet(buffer)

	if !reflect.DeepEqual(packet, parsed_packet) {
		t.Fatalf("packets are not the same: expected [%v] got [%v]", packet, parsed_packet)
	}
}
