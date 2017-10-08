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

// The Level scene loads the next level in line (or resets, if we're out of
// levels), and places all of the tiles for that level.

func LevelStart(prevScene string, data interface{}) {
	levelName := levels[currentLevel]
	l := levelStore[levelName]
	l.Place()
	// This is the 'reset level' binding.
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

// Level scenes continue until something sets this boolean flag to true

func LevelLoop() bool {
	return !levelComplete
}

// When a level ends, it increments to the next level and sets up a fade
// transition.

func LevelEnd() (string, *oak.SceneResult) {
	levelInit = false
	levelComplete = false
	currentLevel = (currentLevel + 1) % len(levels)
	res := &oak.SceneResult{
		Transition: oak.TransitionFade(.1, 50),
	}
	return "level", res
}
