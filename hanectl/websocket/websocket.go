package websocket

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024}

var hub = newHub()

type StatusCallback = func(connected bool)

func StartWebSocket(rapi chi.Router, statusCallback StatusCallback) {
	hub.run()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	rapi.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Error().Msgf("Upgrade for websocket failed: %v", err)
			statusCallback(false)
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client
		statusCallback(true)

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()
	})
}

func Broadcast(message string) {
	hub.Broadcast(message)
}
