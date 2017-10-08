package game

import (
	"github.com/200sc/klangsynthese/font"
	"github.com/oakmound/oak/scene"
)

var (
	sounds *font.Font
)

// The first scene does background intialization before the real game starts

func FirstSceneInit(string, interface{}) {
	initTiles()
	sounds = font.New()
}

func FirstSceneLoop() bool {
	return false
}

func FirstSceneEnd() (string, *scene.Result) {
	return "level", nil
}
