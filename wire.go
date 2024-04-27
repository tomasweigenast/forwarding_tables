package main

// Wire represents a physical wire where devices can connect
type Wire struct {
	buffer chan wire_data // channel where devices connected to this wire can write
}

func new_wire() *Wire {
	return &Wire{
		buffer: make(chan wire_data),
	}
}

// wire_data is a special structure that wraps the []byte data
// and adds an initiator to prevent that a sender reads its own data.
//
// in a future this should be replaced with mac identification
type wire_data struct {
	initiator string
	data      []byte
}
