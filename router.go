package main

import (
	"net"

	"tomasweigenast.com/forwarding_tables/packets"
)

type Router struct {
	ftable     ftable
	interfaces network_interfaces
	_name      string
	_id        string
}

func (r *Router) send(ip net.IP, data packets.IPv4Payload) error {
	packet := new_packet(nil, ip, data)
	return r.forward_packet(ip, packet, false)
}

func (r *Router) name() string {
	return r._name
}

func (r *Router) forward_packet(ip net.IP, packet *packet, is_fordward ...bool) error {
	must_track := true
	if len(is_fordward) == 1 {
		must_track = is_fordward[0]
	}

	// lookup in table
	next_hop, interface_name, err := r.ftable.lookup(ip)
	if err != nil {
		default_logger.infof("forwarding error on device %q trying to send to %q: %s", r._name, ip, err)
		return ErrForwardingError
	}

	interf, err := r.interfaces.get_interface(interface_name)
	if err != nil {
		return ErrInterfaceNotFound
	}

	if must_track {
		default_network_recorder.notify_jump(packet.id, interf.ip.String(), next_hop.String())
	}

	// set packet sender before sending
	// todo: packet sender must not be changed
	// packet.sender = interf.ip

	return interf.output_data(packet.encode())
}

func new_router(name string) *Router {
	return &Router{
		ftable:     newftable(),
		interfaces: new_network_interfaces(),
		_name:      name,
		_id:        name, //random_id(),
	}
}

func (r *Router) add_interface(name string, ip string) {
	r.interfaces.add_network_interface(name, ip, r)
}

func (r *Router) id() string {
	return r._id
}

func (r *Router) start() {
	default_logger.infof("Starting router: %s", r._name)
}

func (r *Router) handle_packet(p packet, i *network_interface) {
	// if destination is the router interface, accept it, otherwise try to send to another device
	if p.destination.Equal(i.ip) {
		dataString := string(p.data)
		default_logger.file_logf("router %s [%s: %s] received packet from %s: %s", r._name, i.name, i.ip, p.sender, dataString)
		handle_packet_payload(r, &p)
		return
	}

	// send to another device
	err := r.forward_packet(p.destination, &p)
	if err != nil {
		default_logger.infof("an error occurred trying to forward a packet to %s: %s", p.destination, err)
	} else {
		default_logger.file_logf("jumped from %s [%s] to %s", r._name, i.ip, p.destination)
	}
}
