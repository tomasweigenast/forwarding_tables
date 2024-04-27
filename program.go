package main

import (
	"fmt"

	"tomasweigenast.com/forwarding_tables/packets"
)

// ProgramExec represents the actual function that parses an IPv4 payload
type ProgramExec func(device Device, payload []byte)

// Program handles a specific IPv4 packet
type Program struct {
	protocol  packets.IPv4Protocol // the protocol this program handles
	execution ProgramExec          // the program's code
}

var programs map[packets.IPv4Protocol]*Program = map[packets.IPv4Protocol]*Program{
	packets.ICMP_Protocol: {
		protocol:  packets.ICMP_Protocol,
		execution: icmp_handle,
	},
}

func handle_packet_payload(d Device, p *packet) {
	program := programs[p.protocol]
	if program == nil {
		panic(fmt.Errorf("there is no program to handle protocol %d", p.protocol))
	}

	program.execution(d, p.data)
}

func icmp_handle(device Device, payload []byte) {
	icmp_packet := packets.DecodeICMP(payload)
	default_logger.infof("icmp_packet: type [%d] code [%d]", icmp_packet.Type, icmp_packet.Code)
}
