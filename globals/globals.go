package globals

import (
	"server/game"
	"server/player"
)

var (
	PlayerMgr = player.NewPlayerManager()
	GameMgr   = game.NewGameManager()
)
