package main

import "net"

// Device is the basic unit of work that can send and receive data
type Device interface {
	// send sends data to ip
	send(ip net.IP, data []byte) error

	// handle_packet handles a packet p that came from i
	handle_packet(p packet, i *network_interface)

	id() string
}
