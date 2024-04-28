package main

import (
	"crypto/rand"
	"encoding/hex"
)

func random_id() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func random_packet_id() byte {
	b := make([]byte, 1)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return b[0]
}
