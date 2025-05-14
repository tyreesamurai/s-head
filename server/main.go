package main

import (
	"log"
	"net/http"
	"server/connections"
)

func main() {
	http.HandleFunc("/", connections.HandleConnections)

	serverPort := ":8080"
	log.Println("WebSocket Server is running on port", serverPort)
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatal("Server Error:", err)
	}
}
