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
	clients map[string]*Client

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
		rooms: 		make(map[string]*Room),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

// handles register, unregister, and broadcast events
// broadcast messages to the right rooms
func (h *Hub) run() {
	for {
		select {
			case client := <-h.register:
				// on registration, add id -> client mapping to pool
				h.clients[client.id] = client

			case client := <-h.unregister:
				// remove client from memory and close the client send channel
				if _, ok := h.clients[client.id]; ok {
					delete(h.clients, client.id)
					close(client.send) // close the client send channel
				}
			case message := <-h.broadcast: // when a ws message comes in
				// transform it to a Message type
				var m Message
				json.Unmarshal(message, &m)

				// 3 types of messages - join/leave room, and actual send messages
				switch m.Action {
					case "joins":
						// find room, lookup client, add to room
						h.rooms[m.Room].clients[h.clients[m.Name]] = true
					case "leaves":
						// find room, remove client from room
						delete(h.rooms[m.Room].clients, h.clients[m.Name])
					case "sends":
						// send message to all clients in that room
						room := h.rooms[m.Room]
						for client := range room.clients {
							select {
							case client.send <- []byte(m.Message):

							default: // if send buffer is full, assume client is dead or stuck, unregister client, close websocket
								close(client.send)
								delete(h.clients, client.id)
							}
						}
					default:

				}


		}
	}
}

