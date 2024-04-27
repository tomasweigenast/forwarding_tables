package main

import (
	"fmt"
	"net"
	"sync"
)

// network_notifier is a manager that creates new wires or attaches to them network interfaces
type network_notifier struct {
	mx    sync.RWMutex
	wires map[string]*Wire
}

// default_network_notifier is the default instance of a network_notifier
var default_network_notifier *network_notifier = &network_notifier{wires: make(map[string]*Wire)}

// get_network_wire returns a new wire for a given ip
//
// the function gets the ip network and creates a new one for it if there are no other one
func (nn *network_notifier) get_network_wire(network net.IPNet) *Wire {
	nn.mx.Lock()
	defer nn.mx.Unlock()

	network_key := network.String()
	wire, ok := nn.wires[network_key]
	if !ok {
		wire = new_wire()
		nn.wires[network_key] = wire
		fmt.Printf("created a new wire for network %s\n", network.String())
	}

	return wire
}
