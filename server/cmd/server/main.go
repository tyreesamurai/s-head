package main

import (
	"log"
	"net/http"
	"shithead/internal/httpserver"
	"shithead/internal/ws"
)

func main() {
	http.HandleFunc("/ws", ws.Handler)
	http.HandleFunc("/", httpserver.Handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
