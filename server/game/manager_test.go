package game

import (
	"server/models"
	"testing"
)

func TestNewGameManager(t *testing.T) {
	gm := NewGameManager()
	if gm == nil {
		t.Fatal("Expected GameManager to be initialized, got nil")
	}
	if len(gm.games) != 0 {
		t.Fatalf("Expected no games in GameManager, got %d", len(gm.games))
	}
}

func TestCreateGame(t *testing.T) {
	gm := NewGameManager()
	game := &Game{ID: "game1"}

	gm.CreateGame(game)

	if _, exists := gm.games["game1"]; !exists {
		t.Fatal("Expected game to be created, but it does not exist")
	}
}

func TestCreateGameWhenGamesNil(t *testing.T) {
	gm := &GameManager{} // Manually make gm.games nil

	game := &Game{ID: "game_nil_test"}
	gm.CreateGame(game)

	if _, exists := gm.games["game_nil_test"]; !exists {
		t.Fatal("Expected game to be created even when games map was initially nil")
	}
}

func TestCreateGameOverwrite(t *testing.T) {
	gm := NewGameManager()

	game1 := &Game{ID: "game1", NumberOfPlayers: 2}
	game2 := &Game{ID: "game1", NumberOfPlayers: 4} // Same ID, different data

	gm.CreateGame(game1)
	gm.CreateGame(game2)

	retrieved, _ := gm.GetGame("game1")
	if retrieved.NumberOfPlayers != 4 {
		t.Fatalf("Expected game to be overwritten with new data, got %d players", retrieved.NumberOfPlayers)
	}
}

func TestGetGame(t *testing.T) {
	gm := NewGameManager()
	game := &Game{ID: "game1"}
	gm.CreateGame(game)

	retrievedGame, exists := gm.GetGame("game1")
	if !exists {
		t.Fatal("Expected game to exist, but it does not")
	}
	if retrievedGame.ID != "game1" {
		t.Fatalf("Expected game ID to be 'game1', got '%s'", retrievedGame.ID)
	}
}

func TestGetNonExistingGame(t *testing.T) {
	gm := NewGameManager()

	_, exists := gm.GetGame("nonexistent")

	if exists {
		t.Fatal("Expected game to not exist, but it does")
	}
}

func TestGetAllGames(t *testing.T) {
	gm := NewGameManager()
	gm.CreateGame(&Game{ID: "game1"})
	gm.CreateGame(&Game{ID: "game2"})

	allGames := gm.GetAllGames()
	if len(allGames) != 2 {
		t.Fatalf("Expected 2 games, got %d", len(allGames))
	}
}

func TestGetAllGamesEmpty(t *testing.T) {
	gm := NewGameManager()

	allGames := gm.GetAllGames()

	if len(allGames) != 0 {
		t.Fatalf("Expected 0 games, got %d", len(allGames))
	}
}

func FuzzCreateAndGetGame(f *testing.F) {
	gm := NewGameManager()

	f.Add("randomID1")
	f.Add("anotherGameID")
	f.Fuzz(func(t *testing.T, gameID string) {
		game := &Game{ID: gameID}
		gm.CreateGame(game)

		retrieved, exists := gm.GetGame(gameID)
		if !exists {
			t.Errorf("Game with ID '%s' should exist but does not", gameID)
		}
		if retrieved.ID != gameID {
			t.Errorf("Retrieved game ID '%s' does not match expected '%s'", retrieved.ID, gameID)
		}
	})
}

func TestRemoveGame(t *testing.T) {
	gm := NewGameManager()
	gm.CreateGame(&Game{ID: "game1"})

	gm.RemoveGame("game1")

	if _, exists := gm.GetGame("game1"); exists {
		t.Fatal("Expected game to be removed, but it still exists")
	}
}

func TestRemoveNonExistingGame(t *testing.T) {
	gm := NewGameManager()

	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Expected no panic when removing non-existing game")
		}
	}()

	gm.RemoveGame("ghost_game")
}

func TestAddPlayerToGame(t *testing.T) {
	gm := NewGameManager()

	game := &Game{
		ID:              "game1",
		Players:         make(map[*models.Player][]models.Card),
		NumberOfPlayers: 2,
	}
	player := &models.Player{Name: "player1"}

	gm.AddPlayerToGame(game, player)

	if len(game.Players) != 1 {
		t.Fatalf("Expected 1 player in the game, got %d", len(game.Players))
	}
	if _, exists := game.Players[player]; !exists {
		t.Fatal("Expected player to be added to the game, but it does not exist")
	}
}

func TestAddPlayerToFullGame(t *testing.T) {
	gm := NewGameManager()

	game := &Game{
		ID:              "game1",
		Players:         make(map[*models.Player][]models.Card),
		NumberOfPlayers: 1,
	}
	player1 := &models.Player{Name: "player1"}
	player2 := &models.Player{Name: "player2"}

	gm.AddPlayerToGame(game, player1)
	gm.AddPlayerToGame(game, player2)

	if len(game.Players) != 1 {
		t.Fatalf("Expected 1 player in the game, got %d", len(game.Players))
	}
}

func TestAddSamePlayerTwice(t *testing.T) {
	gm := NewGameManager()

	game := &Game{
		ID:              "sameplayertest",
		Players:         make(map[*models.Player][]models.Card),
		NumberOfPlayers: 2,
	}
	player := &models.Player{Name: "player1"}

	gm.AddPlayerToGame(game, player)
	gm.AddPlayerToGame(game, player)

	if len(game.Players) != 1 {
		t.Fatalf("Expected only 1 entry for player, got %d", len(game.Players))
	}
}

func TestAddNilPlayer(t *testing.T) {
	gm := NewGameManager()

	game := &Game{
		ID:              "nilplayertest",
		Players:         make(map[*models.Player][]models.Card),
		NumberOfPlayers: 2,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Expected no panic when adding nil player")
		}
	}()

	gm.AddPlayerToGame(game, nil)

	if len(game.Players) != 1 {
		t.Fatalf("Expected 1 entry (nil key), got %d", len(game.Players))
	}
}
