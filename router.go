package main

import (
	"fmt"
	"net"
)

type Router struct {
	ftable     ftable
	interfaces network_interfaces
	name       string
	_id        string
}

func (r *Router) send(ip net.IP, data []byte) error {
	packet := new_packet(nil, ip, data)
	return r.forward_packet(ip, packet)
}

func (r *Router) forward_packet(ip net.IP, packet *packet) error {
	// lookup in table
	next_hop, interface_name := r.ftable.lookup(ip)
	if next_hop == nil {
		fmt.Printf("forwarding error on device %q trying to send to %q\n", r.name, ip)
		return ErrForwardingError
	}

	interf, err := r.interfaces.get_interface(interface_name)
	if err != nil {
		// return fmt.Errorf("device %s does not have an interface called %q", r.name, interface_name)
		return ErrInterfaceNotFound
	}

	// set packet sender before sending
	packet.sender = interf.ip

	return interf.output_data(packet.encode())
}

func new_router(name string) *Router {
	return &Router{
		ftable:     newftable(),
		interfaces: new_network_interfaces(),
		name:       name,
		_id:        random_id(),
	}
}

func (r *Router) add_interface(name string, ip string) {
	r.interfaces.add_network_interface(name, ip, r)
}

func (r *Router) id() string {
	return r._id
}

func (r *Router) start() {
	fmt.Printf("Starting router: %s\n", r.name)
}

func (r *Router) handle_packet(p packet, i *network_interface) {
	// if destination is the router interface, accept it, otherwise try to send to another device
	if p.destination.Equal(i.ip) {
		dataString := string(p.data)
		fmt.Printf("router %s [%s: %s] received packet from %s: %s\n", r.name, i.name, i.ip, p.sender, dataString)
		return
	}

	// send to another device
	err := r.forward_packet(p.destination, &p)
	if err != nil {
		fmt.Printf("an error occurred trying to forward a packet to %s: %s\n", p.destination, err)
	} else {
		fmt.Printf("jump from %s to %s\n", r.name, p.destination)
	}
}