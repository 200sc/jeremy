package game

import (
	"time"

	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

type exit struct {
	r render.Renderable
	event.CID
	s *collision.Space
}

// In order to be an entity in oak you need an Init function.
func (e *exit) Init() event.CID {
	e.CID = event.NextID(e)
	return e.CID
}

func exitInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	e := new(exit)
	e.Init()
	e.r = r
	e.s = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(coralExit), e.CID)
	collision.Add(e.s)
	e.Bind(exitLevel, "OpenExit")
}

func exitLevel(id int, nothing interface{}) int {
	e := event.GetEntity(id).(*exit)
	event.Trigger("PausePlayer", nil)
	// Move the image of the exit
	e.Bind(func(id int, frame interface{}) int {
		f := frame.(int)
		if f%7 == 0 {
			e.r.ShiftX(1)
		}
		if f%14 == 0 {
			audio.Play(sounds, "PitEmpty.wav")
		}
		return 0
	}, "EnterFrame")
	// Eventually complete the level
	go timing.DoAfter(time.Second*2, func() {
		levelComplete = true
	})
	return event.UnbindEvent
}
