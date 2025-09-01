package main

import (
	"log"
	"shithead/internal/db"
	"shithead/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	store, err := db.NewStore()
	if err != nil {
		log.Fatal("failed to connect to database: %v", err)
	}

	h := handlers.NewHTTPHandler(store)

	r := gin.Default()

	r.POST("/games", h.CreateGame)
	r.GET("/games", h.GetActiveGames)
	r.DELETE("/games/:id", h.DeleteGameByID)
}
