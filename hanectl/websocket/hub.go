package websocket

import (
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mutex sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}



func (h *Hub) Broadcast(message string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.broadcast <- []byte(message)
	h.broadcast <- []byte{0x1E}
}

func (h *Hub) run() {
	go func() {
		for {
			select {
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			case message := <-h.broadcast:
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}()

	go func() {
		for range time.Tick(time.Minute / 2) {
			h.Broadcast(utils.Json(model.Dictionary{"ping": 1}))
		}
	}()
}
