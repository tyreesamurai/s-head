package models

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Name string `json:"name"`
	Conn *websocket.Conn
}

type SendMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"` // Always a string
}

type ReceiveMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"` // Must be unmarshaled
}

type Card struct {
	Suit string `json:"suit"`
	Rank Rank   `json:"rank"`
}

type Rank struct {
	Rank  string `json:"name"`
	Value int    `json:"value,omitempty"`
}

type Handler func(conn *websocket.Conn, msg ReceiveMessage) error

type CreateGameRequest struct {
	Name            string `json:"name"`
	NumberOfPlayers int    `json:"numberOfPlayers"`
	Private         bool   `json:"private"`
	Password        string `json:"password,omitempty"`
}

type JoinGameRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

type GameStatus int

const (
	Created GameStatus = iota
	WaitingForPlayers
	Running
	Finished
)
