package game

import "github.com/oakmound/oak"

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
	levelInit = true
}

func LevelLoop() bool {
	return !levelComplete
}

func LevelEnd() (string, *oak.SceneResult) {
	levelInit = false
	levelComplete = false
	//currentLevel++
	res := &oak.SceneResult{
		Transition: oak.TransitionFade(.001, 900),
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
	levels     = []string{"level1"}
)
