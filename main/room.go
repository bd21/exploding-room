package main

import (
	"math/rand"
	"time"
)

const (
	// room id length
	roomIdLength = 5

	// character set for room ids
	charSet = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// A room exists in memory and ties an ID to a list of Clients
type Room struct {
	// 5 letter id
	id string
	// Registered clients.
	clients map[string]*Client
	// TODO add message history here
}

//  generate new room with a random id
func newRoom() *Room {
	return &Room{
		id:      stringWithCharset(roomIdLength, charSet),
		clients: make(map[string]*Client),
	}
}

// generate a random string given a length and charset
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

