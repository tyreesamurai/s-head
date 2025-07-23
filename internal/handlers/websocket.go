package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shithead/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// optionally restrict origin
		return true
	},
}

type HandlerFunc func(conn *websocket.Conn, msg models.Message)

var handlers = make(map[string]HandlerFunc)

func RegisterHandler(msgType string, handler HandlerFunc) {
	handlers[msgType] = handler
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection")

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg models.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		if handler, ok := handlers[msg.Type]; ok {
			handler(conn, msg)
		} else {
			log.Printf("No handler for message type: %s\n", msg.Type)
		}
	}
}
