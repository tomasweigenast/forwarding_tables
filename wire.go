package main

// Wire represents a physical wire where devices can connect
type Wire struct {
	buffer chan []byte // channel where devices connected to this wire can write
}

func new_wire() *Wire {
	return &Wire{
		buffer: make(chan []byte),
	}
}
