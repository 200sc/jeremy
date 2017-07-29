package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type Switch struct {
	collision.Phase
	event.CID
	touching int
}

func (s *Switch) Init() event.CID {
	s.CID = event.NextID(s)
	return s.CID
}

func switchInit(t Tile) func(int, int, render.Renderable) {
	gc := t.GateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		s := new(Switch)
		s.Init()
		sp := collision.NewSpace(xf, yf, 16, 16, s.CID)
		collision.Add(sp)
		collision.PhaseCollision(sp)
		s.Bind(switchOn(gc), "CollisionStart")
		s.Bind(switchOff(gc), "CollisionStop")
	}
}

func switchOn(gc GateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		s := event.GetEntity(id).(*Switch)
		switch label.(collision.Label) {
		case collision.Label(Sandglob), collision.Label(JeremyTile):
			s.touching++
			if s.touching == 1 {
				event.Trigger("Open"+gc.String(), nil)
			}
		}
		return 0
	}
}

func switchOff(gc GateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		s := event.GetEntity(id).(*Switch)
		switch label.(collision.Label) {
		case collision.Label(Sandglob), collision.Label(JeremyTile):
			s.touching--
			if s.touching == 0 {
				event.Trigger("Close"+gc.String(), nil)
			}
		}
		return 0
	}
}

func alternatingSwitchInit(t Tile) func(int, int, render.Renderable) {
	gc := t.GateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		s := new(Switch)
		s.Init()
		sp := collision.NewSpace(xf, yf, 16, 16, s.CID)
		collision.Add(sp)
		collision.PhaseCollision(sp)
		s.Bind(alternatingSwitch(gc), "CollisionStart")
	}
}

// theoretically only green gates use alternating switches
func alternatingSwitch(gc GateColor) func(id int, label interface{}) int {
	return func(id int, label interface{}) int {
		s := event.GetEntity(id).(*Switch)
		switch label.(collision.Label) {
		case collision.Label(Sandglob), collision.Label(JeremyTile):
			s.touching++
			if s.touching%2 == 0 {
				event.Trigger("Close"+gc.String(), nil)
			} else {
				event.Trigger("Open"+gc.String(), nil)
			}
		}
		return 0
	}
}
