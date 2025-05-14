package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"server/models"

	"github.com/gorilla/websocket"
)

type Game struct {
	ID                   string                                                  `json:"name"`
	Players              map[*models.Player][]models.Card                        `json:"-"`
	Private              bool                                                    `json:"private"`
	Password             string                                                  `json:"-"`
	Creator              *models.Player                                          `json:"creator"`
	Deck                 []models.Card                                           `json:"-"`
	Pile                 []models.Card                                           `json:"-"`
	Status               string                                                  `json:"status"`
	NumberOfPlayers      int                                                     `json:"numberOfPlayers"`
	SendMessageFunc      func(player *models.Player, message models.SendMessage) `json:"-"`
	BroadcastMessageFunc func(message models.SendMessage)                        `json:"-"`
}

func NewGame(id string, numberOfPlayers int, private bool, password string) (*Game, error) {
	if numberOfPlayers < 2 || numberOfPlayers > 5 {
		return nil, fmt.Errorf("number of players must be between 2 and 5")
	}

	if private && password == "" {
		return nil, fmt.Errorf("password is required for private games")
	}

	if id == "" {
		return nil, fmt.Errorf("game ID cannot be empty")
	}

	if len(id) < 1 || len(id) > 24 {
		return nil, fmt.Errorf("game ID must be between 1 and 24 characters")
	}

	game := &Game{
		ID:              id,
		Players:         make(map[*models.Player][]models.Card),
		NumberOfPlayers: numberOfPlayers,
		Private:         private,
		Deck:            GenerateDeck(),
		Status:          "Waiting For Players",
	}

	game.SendMessageFunc = func(player *models.Player, message models.SendMessage) {
		if player.Conn != nil {
			msgBytes, _ := json.Marshal(message)
			player.Conn.WriteMessage(websocket.TextMessage, msgBytes)
		}
	}

	game.BroadcastMessageFunc = func(message models.SendMessage) {
		for _, player := range game.GetPlayers() {
			if player.Conn != nil {
				msgBytes, _ := json.Marshal(message)
				player.Conn.WriteMessage(websocket.TextMessage, msgBytes)
			}
		}
	}

	return game, nil
}

func GenerateDeck() []models.Card {
	suits := []string{"clubs", "spades", "diamonds", "hearts"}
	ranks := []models.Rank{{Rank: "two", Value: 2}, {Rank: "three", Value: 3}, {Rank: "four", Value: 4}, {Rank: "five", Value: 5}, {Rank: "six", Value: 6}, {Rank: "seven", Value: 7}, {Rank: "eight", Value: 8}, {Rank: "nine", Value: 9}, {Rank: "ten", Value: 10}, {Rank: "jack", Value: 11}, {Rank: "queen", Value: 12}, {Rank: "king", Value: 13}, {Rank: "ace", Value: 14}}
	deck := []models.Card{}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, models.Card{Suit: suit, Rank: rank})
		}
	}

	deck = append(deck, models.Card{Suit: "joker", Rank: models.Rank{Rank: "small", Value: 0}}, models.Card{Suit: "joker", Rank: models.Rank{Rank: "big", Value: 1}})

	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func (g *Game) DrawCards(n int) []models.Card {
	if len(g.Deck) < n {
		return []models.Card{}
	}

	cards := g.Deck[:n]
	g.Deck = g.Deck[n:]
	return cards
}

func (g *Game) AddToPile(card models.Card) {
	g.Pile = append(g.Pile, card)
}

func (g *Game) SetCreator(creator *models.Player) error {
	_, exists := g.Players[creator]

	if !exists {
		if err := g.AddPlayer(creator); err != nil {
			return err
		}
	}

	g.Creator = creator

	return nil
}

func (g *Game) AddPlayer(player *models.Player) error {
	if len(g.Players) >= g.NumberOfPlayers {
		// g.SendMessage(player, models.SendMessage{
		// 	Type:    "error",
		// 	Content: "Game is full, cannot join.",
		// })
		return fmt.Errorf("game is full")
	}

	g.Players[player] = []models.Card{}
	// g.BroadcastMessage(models.SendMessage{Type: "new_player", Content: fmt.Sprintf("Player %s has joined the game!", player.Name)})
	// g.SendMessage(player, models.SendMessage{
	// 	Type:    "welcome",
	// 	Content: fmt.Sprintf("Welcome to the game %s, %s!", g.ID, player.Name),
	// })

	return nil
}

func (g *Game) HasPlayer(player *models.Player) bool {
	_, exists := g.Players[player]
	return exists
}

func (g *Game) GetPlayers() []*models.Player {
	keys := make([]*models.Player, 0, len(g.Players))
	for p := range g.Players {
		keys = append(keys, p)
	}
	return keys
}

func (g *Game) SendMessage(p *models.Player, msg models.SendMessage) error {
	_, exists := g.Players[p]
	if !exists {
		return fmt.Errorf("player %s doesn't exist in this game", p.Name) // Player not found
	}
	msgBytes, _ := json.Marshal(msg)
	p.Conn.WriteMessage(websocket.TextMessage, msgBytes)

	return nil
}

func (g *Game) BroadcastMessage(msg models.SendMessage) {
	msgBytes, _ := json.Marshal(msg)

	for _, player := range g.GetPlayers() {
		player.Conn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}

func (g *Game) UnmarshalJSON(data []byte) error {
	type Alias Game
	aux := &struct {
		Players map[string][]models.Card `json:"players"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	g.Password = aux.Password
	return nil
}
