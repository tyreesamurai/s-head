package models

import (
	"github.com/gorilla/websocket"
)

type Game struct {
	ID         string
	CreatorID  string
	Private    bool
	Password   string
	Status     string
	MaxPlayers int
	WinnerID   *string
}

type Player struct {
	ID   string
	Conn *websocket.Conn
}
