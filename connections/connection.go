package connections

import (
	"encoding/json"
	"log"
	"net/http"
	"server/globals"
	"server/handlers"
	"server/models"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	pm := globals.PlayerMgr
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer func() {
		pm.RemovePlayer(conn)
		conn.Close()
	}()

	log.Println("New Connection Established:", conn.RemoteAddr().String())

	for {
		var msg models.ReceiveMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				log.Println("Websocket closed by client")
			} else {
				log.Printf("Websocket Error: %v", err)
			}
			break
		}

		handler := handlers.GetHandler(msg.Type)

		if handler == nil {
			SendMessage(conn, models.SendMessage{Type: "error", Content: "Unknown message type"})
			continue
		}

		if err := handler(conn, msg); err != nil {
			continue
		}
	}
}

func SendMessage(conn *websocket.Conn, msg models.SendMessage) {
	msgBytes, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)
}
