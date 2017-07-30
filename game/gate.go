package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type GateColor int

const (
	Purple GateColor = iota
	Teal
	Blue
	Green
	NonGate
)

func (t Tile) GateColor() GateColor {
	switch t {
	case PurpleCoralGate, PurpleCoralSwitch:
		return Purple
	case TealCoralGate, TealCoralSwitch:
		return Teal
	case BlueCoralGate, BlueCoralSwitch:
		return Blue
		// Todo green
	}
	return NonGate
}

func (gc GateColor) String() string {
	switch gc {
	case Purple:
		return "Purple"
	case Blue:
		return "Blue"
	case Teal:
		return "Teal"
	case Green:
		return "Green"
	}
	return ""
}

type Gate struct {
	physics.Vector
	r  *render.Compound
	s1 *collision.Space
	event.CID
}

func (g *Gate) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func gateInit(t Tile) func(int, int, render.Renderable) {
	gc := t.GateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		g := new(Gate)
		g.Init()
		g.r = r.(*render.Compound)
		g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, Blocking, g.CID)
		collision.Add(g.s1)
		g.Bind(gateOpen, "Open"+gc.String())
		g.Bind(gateClose, "Close"+gc.String())
	}
}

func offGateInit(t Tile) func(int, int, render.Renderable) {
	gc := t.GateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		g := new(Gate)
		g.Init()
		g.r = r.(*render.Compound)
		g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, Blocking, g.CID)
		g.Bind(gateOpen, "Close"+gc.String())
		g.Bind(gateClose, "Open"+gc.String())
	}
}

func gateOpen(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*Gate)
	g.r.Set("open")
	collision.Remove(g.s1)
	return 0
}

func gateClose(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*Gate)
	g.r.Set("closed")
	collision.Add(g.s1)
	return 0
}
