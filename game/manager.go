package game

import (
	"log"
	"server/models"
	"sync"
)

type GameManager struct {
	mutex sync.Mutex
	games map[string]Game
}

const (
	Created models.GameStatus = iota
	WaitingForPlayers
	Running
	Finished
)

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]Game),
	}
}

func (gm *GameManager) CreateGame(gameToCreate *Game) {
	if gm.games == nil {
		gm.games = make(map[string]Game)
	}

	gm.games[gameToCreate.ID] = *gameToCreate

	log.Println("Game created:", gameToCreate.ID)
}

func (gm *GameManager) GetGame(gameID string) (Game, bool) {
	game, exists := gm.games[gameID]
	return game, exists
}

func (gm *GameManager) GetAllGames() []Game {
	games := make([]Game, 0, len(gm.games))

	for _, game := range gm.games {
		games = append(games, game)
	}

	return games
}

func (gm *GameManager) RemoveGame(gameID string) {
	delete(gm.games, gameID)
}

func (gm *GameManager) AddPlayerToGame(game *Game, player *models.Player) {
	if len(game.Players) >= game.NumberOfPlayers {
		return
	}

	game.Players[player] = []models.Card{}
}
