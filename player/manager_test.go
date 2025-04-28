package player

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestManagerAddPlayer(t *testing.T) {
	pm := NewPlayerManager()

	conn := new(websocket.Conn)

	player := pm.AddPlayer(conn, "testPlayer")

	if player.Name != "testPlayer" {
		t.Errorf("Expected player name to be 'testPlayer', got '%s'", player.Name)
	}
}

func TestManagerAddNilConnection(t *testing.T) {
	pm := NewPlayerManager()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Expected no panic when adding nil connection, got panic: %v", r)
		}
	}()

	player := pm.AddPlayer(nil, "nilPlayer")
	if player.Name != "nilPlayer" {
		t.Errorf("Expected player name to be 'nilPlayer', got '%s'", player.Name)
	}
}

func TestManagerRemovePlayer(t *testing.T) {
	pm := NewPlayerManager()
	conn := new(websocket.Conn)
	pm.AddPlayer(conn, "testPlayer")
	pm.RemovePlayer(conn)
	if _, exists := pm.players[conn]; exists {
		t.Errorf("Expected player to be removed, but it still exists")
	}
}

func TestManagerRemoveNilConnection(t *testing.T) {
	pm := NewPlayerManager()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Expected no panic when removing nil connection, got panic: %v", r)
		}
	}()

	pm.RemovePlayer(nil)
}

func TestManagerGetPlayer(t *testing.T) {
	pm := NewPlayerManager()
	conn := new(websocket.Conn)
	pm.AddPlayer(conn, "testPlayer")
	player, exists := pm.GetPlayer(conn)
	if !exists {
		t.Errorf("Expected player to exist, but it does not")
	}
	if player.Name != "testPlayer" {
		t.Errorf("Expected player name to be 'testPlayer', got '%s'", player.Name)
	}
}

func TestManagerAddPlayerWithExistingConnection(t *testing.T) {
	pm := NewPlayerManager()
	conn := new(websocket.Conn)

	pm.AddPlayer(conn, "testPlayer")
	player := pm.AddPlayer(conn, "updatedPlayer")

	if player.Name != "updatedPlayer" {
		t.Errorf("Expected player name to be updated to 'updatedPlayer', got '%s'", player.Name)
	}
}

func TestManagerRemoveNonExistentPlayer(t *testing.T) {
	pm := NewPlayerManager()
	conn := new(websocket.Conn)

	pm.RemovePlayer(conn) // Should not panic or cause errors
	if _, exists := pm.players[conn]; exists {
		t.Errorf("Expected no player to exist for the connection, but found one")
	}
}

func TestManagerGetNonExistentPlayer(t *testing.T) {
	pm := NewPlayerManager()
	conn := new(websocket.Conn)

	player, exists := pm.GetPlayer(conn)
	if exists {
		t.Errorf("Expected player to not exist, but found one: %+v", player)
	}
}

func TestManagerAddMultiplePlayers(t *testing.T) {
	pm := NewPlayerManager()
	conn1 := new(websocket.Conn)
	conn2 := new(websocket.Conn)

	pm.AddPlayer(conn1, "player1")
	pm.AddPlayer(conn2, "player2")

	if len(pm.players) != 2 {
		t.Errorf("Expected 2 players, but found %d", len(pm.players))
	}

	if pm.players[conn1].Name != "player1" || pm.players[conn2].Name != "player2" {
		t.Errorf("Player names do not match expected values")
	}
}

func TestManagerGetAllPlayers(t *testing.T) {
	pm := NewPlayerManager()
	conn1 := new(websocket.Conn)
	conn2 := new(websocket.Conn)

	pm.AddPlayer(conn1, "player1")
	pm.AddPlayer(conn2, "player2")

	players := pm.GetAllPlayers()

	if len(players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(players))
	}
}

func TestManagerGetAllPlayersEmpty(t *testing.T) {
	pm := NewPlayerManager()

	players := pm.GetAllPlayers()

	if len(players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(players))
	}
}

func FuzzAddPlayer(f *testing.F) {
	pm := NewPlayerManager()

	f.Add("testPlayer")
	f.Fuzz(func(t *testing.T, name string) {
		conn := new(websocket.Conn)
		pm.AddPlayer(conn, name)
	})
}
