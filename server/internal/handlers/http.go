package handlers

import (
	"net/http"
	"shithead/internal/db"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	Store *db.Store
}

func NewHTTPHandler(store *db.Store) *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) CreateGame(ctx *gin.Context) {
	var game db.Game
	if err := ctx.ShouldBindJSON(&game); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := h.Store.CreateGame(&game); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create game"})
		return
	}

	ctx.JSON(http.StatusCreated, game)
}

func (h *HTTPHandler) GetActiveGames(ctx *gin.Context) {
	games, err := h.Store.GetActiveGames()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve active games"})
		return
	}

	if len(games) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "no active games found"})
		return
	}

	ctx.JSON(http.StatusOK, games)
}

func (h *HTTPHandler) GetGameByID(ctx *gin.Context) {
	gameID := ctx.Param("id")

	game, err := h.Store.GetGameByID(gameID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve game"})
		return
	}

	if game == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "game not found"})
		return
	}

	ctx.JSON(http.StatusOK, game)
}

func (h HTTPHandler) DeleteGameByID(ctx *gin.Context) {
	var game db.Game

	if err := ctx.ShouldBindJSON(&game); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	h.Store.DeleteGameByID(game.ID)
}
