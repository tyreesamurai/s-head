package handlers

import (
	"encoding/json"
	"log"
	"server/game"
	"server/globals"
	"server/models"
	"strings"

	"github.com/gorilla/websocket"
)

var handlers = map[string]models.Handler{
	"register":    AddUser,
	"create_game": CreateGame,
	"unregister":  RemoveUser,
	"join_game":   JoinGame,
	"all_games":   GetAllGames,
}

func GetHandler(msgType string) models.Handler {
	if handler, exists := handlers[msgType]; exists {
		return handler
	}
	return nil
}

func AddUser(conn *websocket.Conn, msg models.ReceiveMessage) error {
	var name string

	if err := json.Unmarshal([]byte(msg.Content), &name); err != nil {
		log.Println("User provided an invalid name format")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Invalid name format"})
	}

	if name == "" {
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Name is required"})
		log.Println("User attempted to register without a name")
		return websocket.ErrBadHandshake
	}

	if name == "" {
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Name is required"})
		log.Println("User attempted to register without a name")
		return websocket.ErrBadHandshake
	}

	if len(name) < 1 || len(name) > 24 {
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Name must be less than 24 characters"})
		log.Println("User provided an invalid name length")
		return websocket.ErrBadHandshake
	}

	if strings.ContainsAny(name, "\x00\x1f\x7f") {
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Name contains invalid control characters"})
		log.Println("User provided a name with invalid control characters")
		return websocket.ErrBadHandshake
	}

	playerManager := globals.PlayerMgr

	players := playerManager.GetAllPlayers()
	for _, player := range players {
		if player.Name == name {
			SendMessage(conn, models.SendMessage{Type: "error", Content: "Name already taken"})
			log.Println("User attempted to register with a name in use")
			return websocket.ErrBadHandshake
		}
	}

	playerManager.AddPlayer(conn, name)

	SendMessage(conn, models.SendMessage{Type: "success", Content: "Welcome " + name + "!"})
	gm := globals.GameMgr

	all_games := gm.GetAllGames()

	all_games_string, err := json.Marshal(all_games)
	if err != nil {
		log.Println("Error marshalling game data:", err)
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Error fetching games"})
	}

	SendMessage(conn, models.SendMessage{Type: "all_games", Content: string(all_games_string)})

	log.Printf("Active players: %d", len(players)+1)

	return nil
}

func CreateGame(conn *websocket.Conn, msg models.ReceiveMessage) error {
	playerManager := globals.PlayerMgr

	player, exists := playerManager.GetPlayer(conn)

	if !exists || player.Name == "" {
		log.Println("User attempted to create a game without being registered")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "User not registered"})
		return websocket.ErrBadHandshake
	}

	var gameReq models.CreateGameRequest

	if err := json.Unmarshal([]byte(msg.Content), &gameReq); err != nil {
		log.Println("User provided an invalid game request format")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Invalid game request format"})
		return websocket.ErrBadHandshake
	}

	if gameReq.Name == "" {
		log.Println("User provided an empty game name")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Game name is required"})
		return websocket.ErrBadHandshake
	}

	if len(gameReq.Name) < 1 || len(gameReq.Name) > 24 {
		log.Println("User provided an invalid game name length")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Game name must be less than 24 characters"})
		return websocket.ErrBadHandshake
	}

	if strings.ContainsAny(gameReq.Name, "\x00\x1f\x7f") {
		log.Println("User provided a name with invalid control characters")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Name contains invalid control characters"})
		return websocket.ErrBadHandshake
	}

	if gameReq.NumberOfPlayers < 2 || gameReq.NumberOfPlayers > 5 {
		log.Println("User provided an invalid number of players")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Number of players must be between 2 and 5"})
		return websocket.ErrBadHandshake
	}

	if gameReq.Private && (strings.TrimSpace(gameReq.Password) == "") {
		log.Println("User attempted to create a private game without a password")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Password is required for private games"})
		return websocket.ErrBadHandshake
	}

	game, err := game.NewGame(gameReq.Name, gameReq.NumberOfPlayers, gameReq.Private, gameReq.Password)
	if err != nil {
		log.Println("Error creating game:", err)
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Error creating game"})
	}

	if err := game.SetCreator(player); err != nil {
		log.Println("Error setting game creator:", err)
		return err
	}

	gameManager := globals.GameMgr

	gameManager.CreateGame(game)

	gameBytes, err := json.Marshal(game)
	if err != nil {
		log.Println("Error marshalling game data:", err)
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Error creating game"})
	}

	log.Println("Game created successfully:", game.ID)
	SendMessage(conn, models.SendMessage{Type: "success", Content: string(gameBytes)})

	playerManager.BroadcastMessage(models.SendMessage{Type: "new_game", Content: string(gameBytes)})

	return nil
}

func RemoveUser(conn *websocket.Conn, msg models.ReceiveMessage) error {
	playerManager := globals.PlayerMgr
	player, exists := playerManager.GetPlayer(conn)
	if !exists {
		log.Println("User attempted to unregister without being registered")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "User not registered"})
		return websocket.ErrBadHandshake
	}
	playerManager.RemovePlayer(conn)
	SendMessage(conn, models.SendMessage{Type: "success", Content: "Goodbye " + player.Name + "!"})
	return nil
}

func JoinGame(conn *websocket.Conn, msg models.ReceiveMessage) error {
	playerManager := globals.PlayerMgr
	player, exists := playerManager.GetPlayer(conn)

	if !exists {
		log.Println("User attempted to join a game without being registered")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "User not registered"})
		return websocket.ErrBadHandshake
	}

	var gameReq models.JoinGameRequest

	if err := json.Unmarshal([]byte(msg.Content), &gameReq); err != nil {
		log.Println("User provided an invalid game request format")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Invalid game request format"})
		return websocket.ErrBadHandshake
	}
	gameManager := globals.GameMgr

	game, exists := gameManager.GetGame(gameReq.Name)

	if !exists {
		log.Println("Game not found")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Game not found"})
		return websocket.ErrBadHandshake
	}
	if game.Private && (strings.TrimSpace(gameReq.Password) == "") {
		log.Println("User attempted to join a private game without a password")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Password is required for private games"})
		return websocket.ErrBadHandshake
	}

	if game.Private && (game.Password != gameReq.Password) {
		log.Println("User provided an incorrect password for the private game")
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Incorrect password"})
		return websocket.ErrBadHandshake
	}

	if err := game.AddPlayer(player); err != nil {
		log.Println("Error adding player to the game:", err)
		return err
	}

	log.Printf("Player %s joined the game %s", player.Name, game.ID)
	SendMessage(conn, models.SendMessage{Type: "success", Content: "Joined the game successfully!"})
	return nil
}

func GetAllGames(conn *websocket.Conn, msg models.ReceiveMessage) error {
	gameManager := globals.GameMgr
	games := gameManager.GetAllGames()
	gamesString, err := json.Marshal(games)
	if err != nil {
		log.Println("Error marshalling game data:", err)
		SendMessage(conn, models.SendMessage{Type: "error", Content: "Error fetching games"})
		return err
	}
	SendMessage(conn, models.SendMessage{Type: "success", Content: string(gamesString)})
	return nil
}

func SendMessage(conn *websocket.Conn, msg models.SendMessage) {
	msgBytes, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, msgBytes)
}
