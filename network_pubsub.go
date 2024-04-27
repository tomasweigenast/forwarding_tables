package main

import (
	"sync"
)

// the default buffer size for new channels
const channel_buffer_size = 0

type pubsub struct {
	mx       sync.RWMutex
	channels map[string][]client
	closed   bool
}

type client struct {
	channel chan []byte
	id      string
}

var default_pubsub *pubsub = &pubsub{
	mx:       sync.RWMutex{},
	channels: make(map[string][]client),
}

// subscribe subscribes a new client with id to a network
//
// todo: id should be changed with mac address in a future
func (ps *pubsub) subscribe(id, network string) (receive <-chan []byte, send func(data []byte)) {
	if ps.closed {
		panic("pubsub closed")
	}

	ps.mx.Lock()
	defer ps.mx.Unlock()

	channel := make(chan []byte, channel_buffer_size)
	ps.channels[network] = append(ps.channels[network], client{
		channel: channel,
		id:      id,
	})

	return channel, func(data []byte) {
		ps.publish(id, network, data)
	}
}

// publish publishes data to a network from the device with id
func (ps *pubsub) publish(id, network string, data []byte) {
	ps.mx.RLock()
	defer ps.mx.RUnlock()

	clients := ps.channels[network]
	default_logger.log(0, "--------------")
	default_logger.infof("there are %d clients for network %s", len(clients), network)
	default_logger.infof("clients:")
	for _, client := range clients {
		default_logger.infof("	device id: %s", client.id)
	}
	default_logger.log(0, "--------------")
	for _, client := range clients {
		// ignore sender
		if client.id == id {
			continue
		}

		client.channel <- data
	}
}

func (ps *pubsub) close() {
	if ps.closed {
		return
	}

	ps.mx.Lock()
	defer ps.mx.Unlock()

	for _, channel := range ps.channels {
		for _, client := range channel {
			close(client.channel)
		}
	}

	ps.closed = true
}
