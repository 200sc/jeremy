package game

import (
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/event"
)

var (
	envFriction   = 0.8
	levelInit     bool
	currentLevel  int
	levelComplete bool
)

func LevelStart(prevScene string, data interface{}) {
	levelName := levels[currentLevel]
	l := levelStore[levelName]
	l.Place()
	event.GlobalBind(func(int, interface{}) int {
		ok, d := oak.IsHeld("R")
		if ok && d > time.Millisecond*1500 {
			currentLevel--
			levelComplete = true
		}
		return 0
	}, "EnterFrame")
	levelInit = true
}

func LevelLoop() bool {
	return !levelComplete
}

func LevelEnd() (string, *oak.SceneResult) {
	levelInit = false
	levelComplete = false
	currentLevel++
	res := &oak.SceneResult{
		Transition: oak.TransitionFade(.001, 500),
	}
	return "level", res
}

type level [20][15]Tile

func (l level) Place() {
	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {
			l[x][y].Place(x, y)
		}
	}
}

var (
	levelStore = make(map[string]level)
	levels     = []string{"tutorial1", "tutorial2", "tutorial3", "tutorial4", "tutorial5", "tutorial6", "level1", "level2", "level3", "level4", "level5", "level6", "level7"}
)
