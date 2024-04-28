package main

import "errors"

var (
	ErrForwardingError   = errors.New("forwarding-error")    // error that indicates a device does not know how to handle a packet and may drop it
	ErrInterfaceNotFound = errors.New("interface-not-found") // indicates that a given interface is not present in a device
	ErrNoNextHop         = errors.New("no-valid-next-hop")   // indicates that the device does not have a proper entry in its ftable for a redirection
)
