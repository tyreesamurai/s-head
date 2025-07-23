package handlers

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/games", func(c *gin.Context) {
	})

	router.POST("/games", func(c *gin.Context) {
	})

	router.POST("/games/:id/join", func(c *gin.Context) {
	})

	router.POST("/players/register", func(c *gin.Context) {
	})

	router.Run()
}
