package main

import (
	"github.com/200sc/jeremy/game"
	"github.com/oakmound/oak"
)

func main() {
	oak.Add("level", game.LevelStart, game.LevelLoop, game.LevelEnd)
	oak.Add("first", game.FirstSceneInit, game.FirstSceneLoop, game.FirstSceneEnd)
	oak.LoadConf("oak.config")
	oak.Init("first")
}
