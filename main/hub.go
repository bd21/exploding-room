package main

import (
	"encoding/json"
)

// Maintains state of clients and rooms
//
type Hub struct {

	// map the room id to the room
	// this is used to route messages
	rooms map[string]*Room

	// Registered clients.  Bool is always true
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	Name string
	Action string
	Room string
	Message string
}


func newHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// handles register, unregister, and broadcast events
// broadcast messages to the right rooms
func (h *Hub) run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send) // close the client send channel
				}
			case message := <-h.broadcast:
				// parse message
				var msg Message
				json.Unmarshal([]byte(message), &msg)

				// route message to correct client
				for client := range h.rooms[msg.Room].clients {
					select {
					case client.send <- message:
					default: // if send buffer is full, assume client is dead or stuck, unregister client, close websocket
						close(client.send)
						delete(h.clients, client)
					}
				}
		}
	}
}

