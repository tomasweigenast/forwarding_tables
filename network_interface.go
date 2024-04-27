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
	ip     net.IP            // the interface ip
	device Device            // the device this network_interface is in
	rc     <-chan []byte     // the channel used to receive data
	sc     func(data []byte) // the function used to send data
}

func new_network_interfaces() network_interfaces {
	return network_interfaces{
		m:      make(map[string]network_interface),
		rindex: make(map[string]*network_interface),
	}
}

// add_network_interface adds a new network interface to a device and start listening
func (nis *network_interfaces) add_network_interface(name string, ip string, device Device) {
	pip, network, err := net.ParseCIDR(ip)
	if err != nil {
		panic(fmt.Errorf("unable to parse ip with cidr: %s", err))
	}

	ni := network_interface{name, pip, device, nil, nil}

	// subscribe to network
	ni.rc, ni.sc = default_pubsub.subscribe(device.id(), network.String())
	go ni.listen()

	// save to network_interfaces
	nis.m[name] = ni
	nis.rindex[pip.String()] = &ni
}

// get_interface returns a network_interface by its name
func (nis *network_interfaces) get_interface(name string) (*network_interface, error) {
	if ni, ok := nis.m[name]; ok {
		return &ni, nil
	}

	return nil, errors.New("device interface not found")
}

// output_data sends data thought its wire, if any, otherwise returns an error
func (ni *network_interface) output_data(data []byte) error {
	if ni.sc == nil {
		return fmt.Errorf("network interface [%s] %s is not attached to any wire\n", ni.name, ni.ip)
	}

	// send data
	// ni.wire.buffer <- wire_data{ni.device.id(), data}
	ni.sc(data)
	fmt.Printf("interface %s [%s] sent data\n", ni.name, ni.ip)
	return nil
}

func (ni *network_interface) listen() {
	for {
		// wait for data
		data, ok := <-ni.rc
		if !ok {
			fmt.Println("wire channel is closed")
			break
		}

		fmt.Printf("interface %s [%s] read packet\n", ni.name, ni.ip)

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
