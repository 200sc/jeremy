package game

import "github.com/oakmound/oak"

// The first scene does background intialization before the real game starts

func FirstSceneInit(string, interface{}) {
	initTiles()
}

func FirstSceneLoop() bool {
	return false
}

func FirstSceneEnd() (string, *oak.SceneResult) {
	return "level", nil
}
