package game

import (
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type floorSwitch struct {
	collision.Phase
	event.CID
	touching int
}

func (s *floorSwitch) Init() event.CID {
	s.CID = event.NextID(s)
	return s.CID
}

func switchInit(t tile) func(int, int, render.Renderable) {
	gc := t.gateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		s := new(floorSwitch)
		sp := collision.NewSpace(xf+6, yf+6, 4, 4, s.Init())
		collision.Add(sp)
		collision.PhaseCollision(sp)
		s.Bind(switchOn(gc), "CollisionStart")
		s.Bind(switchOff(gc), "CollisionStop")
	}
}

func switchOn(gc gateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		s := event.GetEntity(id).(*floorSwitch)
		switch label.(collision.Label) {
		case collision.Label(sandglob), collision.Label(jeremyTile), blocking:
			s.touching++
			if s.touching == 1 {
				audio.Play(sounds, "ButtonDown.wav")
				event.Trigger("Open"+gc.String(), nil)
			}
		}
		return 0
	}
}

func switchOff(gc gateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		s := event.GetEntity(id).(*floorSwitch)
		switch label.(collision.Label) {
		case collision.Label(sandglob), collision.Label(jeremyTile), blocking:
			s.touching--
			if s.touching == 0 {
				audio.Play(sounds, "ButtonDown.wav")
				event.Trigger("Close"+gc.String(), nil)
			}
		}
		return 0
	}
}

func alternatingSwitchInit(t tile) func(int, int, render.Renderable) {
	gc := t.gateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		s := new(floorSwitch)
		s.Init()
		sp := collision.NewSpace(xf+6, yf+6, 4, 4, s.CID)
		collision.Add(sp)
		collision.PhaseCollision(sp)
		s.Bind(alternatingSwitch(gc), "CollisionStart")
	}
}

var (
	greenSwitchTracker int
)

// theoretically only green gates use alternating switches
func alternatingSwitch(gc gateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		switch label.(collision.Label) {
		case collision.Label(sandglob), collision.Label(jeremyTile), blocking:
			greenSwitchTracker = (greenSwitchTracker + 1) % 2
			if greenSwitchTracker == 0 {
				audio.Play(sounds, "ButtonDown.wav")
				event.Trigger("Close"+gc.String(), nil)
			} else {
				audio.Play(sounds, "ButtonDown.wav")
				event.Trigger("Open"+gc.String(), nil)
			}
		}
		return 0
	}
}
