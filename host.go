package main

import (
	"fmt"
	"net"
)

type Host struct {
	ftable     ftable
	interfaces network_interfaces
	name       string
}

func new_host(name string) *Host {
	return &Host{
		ftable:     newftable(),
		name:       name,
		interfaces: new_network_interfaces(),
	}
}

func (r *Host) send(ip net.IP, data []byte) error {
	// lookup in table
	next_hop, name := r.ftable.lookup(ip)
	if next_hop == nil {
		fmt.Printf("forwarding error on device %q trying to send to %q\n", r.name, ip)
		return ErrForwardingError
	}

	interf, err := r.interfaces.get_interface(name)
	if err != nil {
		return ErrInterfaceNotFound
	}

	packet := new_packet(interf.ip, ip, data)

	return interf.output_data(packet.encode())
}

func (r *Host) add_interface(name string, ip string) {
	r.interfaces.add_network_interface(name, ip, r)
}

func (r *Host) start() {
	fmt.Printf("Starting host: %s\n", r.name)
}

func (r *Host) handle_packet(p packet, i *network_interface) {
	// if destination is the router interface, accept it, otherwise try to send to another device
	if !p.destination.Equal(i.ip) {
		return
	}

	dataString := string(p.data)
	fmt.Printf("network interface %s [%s] received packet from %s: %s\n", i.name, i.ip, p.sender, dataString)
}
