package session

import (
	"shithead/internal/engine"
	"shithead/internal/models"
)

type Session struct {
	GameID  string
	Engine  *engine.Engine
	Players []models.Player
}
