package player

import (
	"log"
	"server/models"
	"sync"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	mutex   sync.Mutex
	players map[*websocket.Conn]*models.Player
}

func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[*websocket.Conn]*models.Player),
	}
}

func (pm *PlayerManager) AddPlayer(conn *websocket.Conn, name string) models.Player {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.players[conn] = &models.Player{Name: name, Conn: conn}

	log.Println("User registered with name:", name)

	return *pm.players[conn]
}

func (pm *PlayerManager) RemovePlayer(conn *websocket.Conn) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// log.Println("Attempting to remove player for connection:", conn.RemoteAddr().String())

	player, exists := pm.players[conn]
	if !exists {
		log.Println("User not found for unregistration - connection may not exist in map")
		return
	}

	delete(pm.players, conn)
	log.Println("User unregistered with name:", player.Name)
	log.Println("Active connections:", len(pm.players))
}

func (pm *PlayerManager) GetPlayer(conn *websocket.Conn) (*models.Player, bool) {
	if pm.players == nil {
		return nil, false
	}

	pm.mutex.Lock()
	player, exists := pm.players[conn]
	pm.mutex.Unlock()

	return player, exists
}

func (pm *PlayerManager) GetAllPlayers() []*models.Player {
	players := make([]*models.Player, 0, len(pm.players))
	for _, player := range pm.players {
		players = append(players, player)
	}
	return players
}

func (pm *PlayerManager) BroadcastMessage(msg models.SendMessage) error {
	for _, player := range pm.players {
		err := player.Conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error broadcasting message to player:", player.Name, "Error:", err)
			return err
		}
	}

	return nil
}
