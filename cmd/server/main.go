package main

import (
	"log"
	"net/http"
	"shithead/internal/handlers"
	"shithead/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gorilla/websocket"
)

func main() {
	db, err := gorm.Open(sqlite.Open("games.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto migrate your models
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("failed to migrate:", err)
	}

	// register handlers
	handlers.RegisterHandler("chat", handleChat)
	handlers.RegisterHandler("ping", handlePing)

	http.HandleFunc("/ws", handlers.WebSocketHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleChat(conn *websocket.Conn, msg models.Message) {
	log.Println("Chat message received:", msg.Data)
	conn.WriteJSON(models.Message{
		Type: "chat_ack",
		Data: "Received chat message",
	})
}

func handlePing(conn *websocket.Conn, msg models.Message) {
	log.Println("Ping received")
	conn.WriteJSON(models.Message{
		Type: "pong",
	})
}
