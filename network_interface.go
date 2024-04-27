package main

import (
	"errors"
	"fmt"
	"net"
)

type network_interfaces struct {
	m      map[string]network_interface
	rindex map[string]*network_interface
}

type network_interface struct {
	name   string
	ip     net.IP // the interface ip
	wire   *Wire  // the wire this interface is connected to
	device Device // the device this network_interface is in
}

func new_network_interfaces() network_interfaces {
	return network_interfaces{
		m:      make(map[string]network_interface),
		rindex: make(map[string]*network_interface),
	}
}

func (nis *network_interfaces) add_network_interface(name string, ip string, device Device) {
	pip, network, err := net.ParseCIDR(ip)
	if err != nil {
		panic(fmt.Errorf("unable to parse ip with cidr: %s", err))
	}

	ni := network_interface{name, pip, nil, device}

	// get wire
	wire := default_network_notifier.get_network_wire(*network)
	ni.attach_interface_to(wire)

	// save network_interfaces
	nis.m[name] = ni
	nis.rindex[pip.String()] = &ni

}

// attach_interface_to attaches this network interface to the wire w and start listening for incoming data
func (nis *network_interface) attach_interface_to(w *Wire) {
	nis.wire = w
	go nis.listen()
}

func (nis *network_interfaces) get_interface(name string) (*network_interface, error) {
	if ni, ok := nis.m[name]; ok {
		return &ni, nil
	}

	return nil, errors.New("device interface not found")
}

// output_data sends data thought its wire, if any, otherwise returns an error
func (ni *network_interface) output_data(data []byte) error {
	if ni.wire == nil {
		return fmt.Errorf("network interface %s is not attached to any wire", ni.ip)
	}

	// send data
	ni.wire.buffer <- data
	return nil
}

func (ni *network_interface) listen() {
	for {
		// wait for data
		data, ok := <-ni.wire.buffer
		if !ok {
			fmt.Printf("wire channel is closed\n")
			break
		}

		// parse packet
		packet := parse_packet(data)
		if packet == nil {
			fmt.Println("invalid packet format, frame dropped", data)
			continue
		}

		// let the device handle the packet
		ni.device.handle_packet(*packet, ni)
	}
}
