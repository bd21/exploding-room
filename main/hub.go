package main

// Maintains state of clients and rooms

type Hub struct {
	// Registered clients.  Bool is always true
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}


func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}


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
				for client := range h.clients {
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

