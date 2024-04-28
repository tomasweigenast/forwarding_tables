package main

import "fmt"

// network_recorder keeps tracks of jumps between devices
type network_recorder struct {
	tracks map[byte]*[]jump
}

var default_network_recorder *network_recorder = &network_recorder{tracks: make(map[byte]*[]jump)}

// jump represents a jump between two devices
type jump struct {
	source      string
	destination string
}

func (nr *network_recorder) new_recording(packet_id byte) {
	nr.tracks[packet_id] = new([]jump)

}

func (nr *network_recorder) notify_jump(packet_id byte, source, destination string) {
	track, ok := nr.tracks[packet_id]
	if !ok {
		panic(fmt.Errorf("network track %q not found", packet_id))
	}

	*track = append(*track, jump{source, destination})
}

func (nr *network_recorder) get_jumps(packet_id byte) []jump {
	track, ok := nr.tracks[packet_id]
	if !ok {
		return nil
	}

	return *track
}
