package session

import (
	"sync"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	sync.RWMutex
	Players map[string]*Player
}

type Player struct {
	ID   string
	Conn *websocket.Conn
}
