package game

import (
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
)

type Exit struct {
	r render.Renderable
	event.CID
	s *collision.Space
}

func (e *Exit) Init() event.CID {
	e.CID = event.NextID(e)
	return e.CID
}

func exitInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	e := new(Exit)
	e.Init()
	e.r = r
	e.s = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(CoralExit), e.CID)
	collision.Add(e.s)
	e.Bind(exitLevel, "OpenExit")
}

func exitLevel(id int, nothing interface{}) int {
	e := event.GetEntity(id).(*Exit)
	// todo: Pause player / level
	event.Trigger("PausePlayer", nil)
	e.Bind(func(id int, frame interface{}) int {
		f := frame.(int)
		if f%11 == 0 {
			e.r.ShiftX(1)
		}
		return 0
	}, "EnterFrame")
	go timing.DoAfter(time.Second*3, func() {
		levelComplete = true
	})
	return event.UnbindEvent
}
