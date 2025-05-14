package game

import (
	"math/rand"
	"server/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGameValidation(t *testing.T) {
	_, err := NewGame("", 3, false, "")
	assert.Error(t, err, "Game ID should not be empty")

	_, err = NewGame("toolongidtoolongidtoolongidtoolongid", 3, false, "")
	assert.Error(t, err, "Game ID should not be longer than 24 characters")

	_, err = NewGame("validid", 1, false, "")
	assert.Error(t, err, "Number of players too small should fail")

	_, err = NewGame("validid", 6, false, "")
	assert.Error(t, err, "Number of players too big should fail")

	_, err = NewGame("validid", 3, true, "")
	assert.Error(t, err, "Private game without password should fail")
}

func FuzzNewGame(f *testing.F) {
	f.Add("gameid", 2, false, "")
	f.Add("", 3, false, "")
	f.Add("id", 5, true, "pass123")

	f.Fuzz(func(t *testing.T, id string, numberOfPlayers int, private bool, password string) {
		_, _ = NewGame(id, numberOfPlayers, private, password)
	})
}

func TestGenerateDeck(t *testing.T) {
	deck := GenerateDeck()
	assert.Equal(t, 54, len(deck), "Deck should contain 54 cards (52 + 2 jokers)")
}

func FuzzDrawCards(f *testing.F) {
	f.Add(5) // Example starting seed

	f.Fuzz(func(t *testing.T, n int) {
		deckSize := rand.Intn(55) // Up to 54
		game := &Game{Deck: GenerateDeck()[:deckSize]}
		cards := game.DrawCards(n)

		if n > deckSize {
			assert.Equal(t, 0, len(cards), "Should not allow overdraw")
		} else if n >= 0 {
			assert.Equal(t, n, len(cards), "Should draw exactly n cards")
		}
	})
}

func TestDrawCards(t *testing.T) {
	game := &Game{Deck: GenerateDeck()}
	cards := game.DrawCards(5)
	assert.Equal(t, 5, len(cards), "Should draw 5 cards")
	assert.Equal(t, 49, len(game.Deck), "Deck should have 49 cards left")
}

func TestDrawCards_EmptyDeck(t *testing.T) {
	game := &Game{Deck: []models.Card{}}
	cards := game.DrawCards(5)
	assert.Equal(t, 0, len(cards), "Should draw 0 cards when deck is empty")
}

func TestDrawCards_Overdraw(t *testing.T) {
	game := &Game{Deck: GenerateDeck()[:3]} // Only 3 cards
	cards := game.DrawCards(5)
	assert.Equal(t, 0, len(cards), "Should draw 0 cards when trying to draw more than deck size")
}

func TestDrawCardsMoreThanDeck(t *testing.T) {
	game := &Game{Deck: GenerateDeck()}
	cards := game.DrawCards(100)

	assert.Equal(t, 0, len(cards), "Should return empty slice when drawing more than available")
}

func TestAddToPile(t *testing.T) {
	game := &Game{}
	card := models.Card{Suit: "hearts", Rank: models.Rank{Rank: "ace", Value: 14}}
	game.AddToPile(card)
	assert.Equal(t, 1, len(game.Pile), "Pile should contain 1 card")
	assert.Equal(t, card, game.Pile[0], "Card in pile should match the added card")
}

func TestAddPlayer(t *testing.T) {
	game, err := NewGame("testgame", 2, false, "")

	assert.NoError(t, err, "Creating a new game should not return an error")

	game.SendMessageFunc = func(player *models.Player, message models.SendMessage) {}
	game.BroadcastMessageFunc = func(message models.SendMessage) {}

	player := &models.Player{Name: "Player1"}
	err = game.AddPlayer(player)
	assert.NoError(t, err, "Adding a player should not return an error")
	assert.True(t, game.HasPlayer(player), "Player should be added to the game")
}

func TestAddPlayerGameFull(t *testing.T) {
	game, _ := NewGame("fullgame", 2, false, "")

	game.SendMessageFunc = func(player *models.Player, message models.SendMessage) {}
	game.BroadcastMessageFunc = func(message models.SendMessage) {}

	player1 := &models.Player{Name: "Player1"}
	player2 := &models.Player{Name: "Player2"}
	player3 := &models.Player{Name: "Player3"}

	_ = game.AddPlayer(player1)
	_ = game.AddPlayer(player2)
	err := game.AddPlayer(player3)

	assert.Error(t, err, "Adding player to full game should error")
}

func TestSetCreator(t *testing.T) {
	game, err := NewGame("testgame", 2, false, "")

	assert.NoError(t, err, "Creating a new game should not return an error")

	game.SendMessageFunc = func(player *models.Player, message models.SendMessage) {}
	game.BroadcastMessageFunc = func(message models.SendMessage) {}
	player := &models.Player{Name: "Creator"}

	err = game.SetCreator(player)
	assert.NoError(t, err, "Setting the creator should not return an error")
	assert.Equal(t, player, game.Creator, "Creator should be set correctly")
}
