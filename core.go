package main

import (
	"github.com/200sc/jeremy/game"
	"github.com/oakmound/oak"
)

func main() {
	oak.AddScene("level", game.LevelStart, game.LevelLoop, game.LevelEnd)
	oak.LoadConf("oak.config")
	oak.Init("level")
}
