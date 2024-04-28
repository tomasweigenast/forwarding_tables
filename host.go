package main

import (
	"net"

	"tomasweigenast.com/forwarding_tables/packets"
)

type Host struct {
	ftable     ftable
	interfaces network_interfaces
	_name      string
	_id        string
}

func new_host(name string) *Host {
	return &Host{
		ftable:     newftable(),
		_name:      name,
		interfaces: new_network_interfaces(),
		_id:        random_id(),
	}
}

func (h *Host) name() string {
	return h._name
}

func (h *Host) send(ip net.IP, data packets.IPv4Payload) error {
	// lookup in table
	_, name, err := h.ftable.lookup(ip)
	if err != nil {
		default_logger.infof("forwarding error on device %q trying to send to %q: %s", h._name, ip, err)
		return ErrForwardingError
	}

	interf, err := h.interfaces.get_interface(name)
	if err != nil {
		return ErrInterfaceNotFound
	}

	packet := new_packet(interf.ip, ip, data)
	packet.id = random_packet_id()
	default_logger.file_logf("host %s [%s] sending data to %s [packet id: %d]", h._name, interf.ip, ip, packet.id)

	// start recording
	default_network_recorder.new_recording(packet.id)

	return interf.output_data(packet.encode())
}

func (r *Host) id() string {
	return r._id
}

func (r *Host) add_interface(name string, ip string) {
	r.interfaces.add_network_interface(name, ip, r)
}

func (r *Host) handle_packet(p packet, i *network_interface) {
	// if destination is the router interface, accept it, otherwise try to send to another device
	if !p.destination.Equal(i.ip) {
		return
	}

	dataString := string(p.data)
	default_logger.infof("network interface %s [%s] received packet from %s: %s", i.name, i.ip, p.sender, dataString)
	handle_packet_payload(r, &p)

	// if this is the end, stop recording
	default_logger.infof("jumps done ---------------------")
	for _, jump := range default_network_recorder.get_jumps(p.id) {
		default_logger.infof("from %q to %q", jump.source, jump.destination)
	}
	default_logger.infof("------------------------------------------")
}
