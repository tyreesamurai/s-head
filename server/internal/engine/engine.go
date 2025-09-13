package engine

import (
	"log"
	"math/rand"
	"shithead/internal/models"
)

type Engine struct {
	Deck            []Card
	Pile            []Card
	Players         []*models.Player
	Hands           map[string][]Card
	FaceDowns       map[string][]Card
	FaceUps         map[string][]Card
	TurnOrder       []string
	CurrentPlayerID string
}

type Card struct {
	Rank  string
	Suit  string
	Value int
}

func GenerateDeck() []Card {
	var deck []Card

	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}
	rankValues := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8,
		"9": 9, "10": 10, "Jack": 11, "Queen": 12, "King": 13, "Ace": 14,
	}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{
				Rank:  rank,
				Suit:  suit,
				Value: rankValues[rank],
			})
		}
	}

	deck = append(deck, Card{
		Rank:  "Small",
		Suit:  "Joker",
		Value: 0,
	})

	deck = append(deck, Card{
		Rank:  "Big",
		Suit:  "Joker",
		Value: 0,
	})

	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func NewEngine(players []*models.Player) *Engine {
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	turnOrder := make([]string, len(players))
	for i, p := range players {
		turnOrder[i] = p.ID
	}

	return &Engine{
		Deck:      GenerateDeck(),
		Players:   players,
		TurnOrder: turnOrder,
		Hands:     make(map[string][]Card),
		FaceUps:   make(map[string][]Card),
		FaceDowns: make(map[string][]Card),
		Pile:      []Card{},
	}
}

func (e *Engine) StartGame() {
	for _, p := range e.Players {
		e.FaceDowns[p.ID] = e.DrawCards(3)
		e.FaceUps[p.ID] = e.DrawCards(3)
		e.Hands[p.ID] = e.DrawCards(3)
	}
}

func (e *Engine) Run() {
	e.StartGame()

	// TODO: Send initial game state to players

	// TODO: Allow users to switch between face up cards & in their hand

	// TODO: Implement game loop

	e.SetCurrentPlayer(e.FindStartingPlayer())

	e.WaitForPlayerAction()

	e.ValidatePlayerAction()

	e.NextPlayer()
}

func (e *Engine) WaitForPlayerAction() {
}

func (e *Engine) ValidatePlayerAction() {
}

func (e *Engine) DrawCards(count int) []Card {
	if count > len(e.Deck) {
		count = len(e.Deck)
	}
	drawn := e.Deck[:count]
	e.Deck = e.Deck[count:]
	return drawn
}

func (e *Engine) FindStartingPlayer() string {
	lowest := 15
	var candidates []string

	for _, p := range e.Players {
		for _, c := range e.Hands[p.ID] {
			if c.Value > 2 && c.Value < lowest {
				lowest = c.Value
				candidates = []string{p.ID}
			} else if c.Value == lowest {
				candidates = append(candidates, p.ID)
			}
		}
	}

	if len(candidates) == 1 {
		return candidates[0]
	}

	return candidates[rand.Intn(len(candidates))]
}

func (e *Engine) SetCurrentPlayer(playerID string) {
	for _, p := range e.Players {
		if p.ID == playerID {
			e.CurrentPlayerID = playerID
			return
		}
	}
	log.Printf("player %s not found in game\n", playerID)
}

func (e *Engine) SendMessage(playerID string, payload interface{}) {
	for _, p := range e.Players {
		if p.ID == playerID {
			err := p.Conn.WriteJSON(payload)
			if err != nil {
				log.Printf("error sending message to %s: %v\n", playerID, err)
			}
			return
		}
	}
	log.Printf("player %s not found in game\n", playerID)
}

func (e *Engine) Broadcast(payload interface{}) {
	for _, p := range e.Players {
		err := p.Conn.WriteJSON(payload)
		if err != nil {
			log.Printf("error broadcasting to %s: %v\n", p.ID, err)
		}
	}
}

func (e *Engine) NextPlayer() {
	e.TurnOrder = append(e.TurnOrder[1:], e.TurnOrder[0])
	e.CurrentPlayerID = e.TurnOrder[0]
}
