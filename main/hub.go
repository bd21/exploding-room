package main

import (
	"encoding/json"
	"fmt"
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
	Name    string
	Action  string
	Room    string
	Message string
}

func newHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
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

			// send the generated client id back to the client
			// so that they can use this for future requests
			client.send <- []byte("{\"client_id\": \"" + client.id + "\"}")

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

			fmt.Println(m)

			// 3 types of messages - join/leave room, and actual send messages
			switch m.Action {
			case "creates":

				// TODO this is O(n).  Fix this by storing a mapping if client id-> room
				// lookup the old room the client was in
				// and remove them from that room
				for _, oldRoom := range h.rooms {
					for _, tempClient := range oldRoom.clients {
						if tempClient.id == m.Name {
							delete(oldRoom.clients, m.Name)
							break
						}
					}
				}


				// create a room
				room := newRoom()

				// add room mapping
				h.rooms[room.id] = room

				// add client to room
				h.rooms[room.id].clients[m.Name] = h.clients[m.Name]

				// send client the room id
				h.clients[m.Name].send <- []byte("{\"room_id\": \"" + room.id + "\"}")

			case "joins":
				// TODO this is O(n).  Fix this by storing a mapping if client id-> room
				// lookup the old room the client was in
				// and remove them from that room
				for _, oldRoom := range h.rooms {
					for _, tempClient := range oldRoom.clients {
						if tempClient.id == m.Name {
							delete(oldRoom.clients, m.Name)
							break
						}
					}
				}


				// find room, lookup client
				room, roomExists := h.rooms[m.Room]
				if roomExists {
					// add client to new room
					h.rooms[room.id].clients[m.Name] = h.clients[m.Name]
					h.clients[m.Name].send <- []byte("{\"room_id\": \"" + room.id + "\"}")
				}

			case "leaves": // leaving in case a "leave" button is added back in
				// find room, remove client from room
				delete(h.rooms[m.Room].clients, m.Name)
			case "sends":
				// send message to all clients in that room
				room := h.rooms[m.Room]

				if room == nil {
					return
				}
				for clientKey := range room.clients {
					select {
					case room.clients[clientKey].send <- []byte(m.Name + ": " + m.Message):

					default: // if send buffer is full, assume client is dead or stuck, unregister client, close websocket
						close(room.clients[clientKey].send)
						delete(h.clients, room.clients[clientKey].id)
					}
				}
			default:

			}

		}
	}
}
