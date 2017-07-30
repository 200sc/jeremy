package game

import (
	"fmt"
	"sync"

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
	lock   sync.Mutex
	active bool
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
		g.lock = sync.Mutex{}
		g.Init()
		g.r = r.(*render.Compound)
		g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, Blocking, g.CID)
		collision.Add(g.s1)
		g.active = true
		g.Bind(gateOpen, "Open"+gc.String())
		g.Bind(gateClose, "Close"+gc.String())
	}
}

func offGateInit(t Tile) func(int, int, render.Renderable) {
	gc := t.GateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		g := new(Gate)
		g.lock = sync.Mutex{}
		g.Init()
		g.r = r.(*render.Compound)
		g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, Blocking, g.CID)
		g.active = false
		g.Bind(gateOpen, "Close"+gc.String())
		g.Bind(gateClose, "Open"+gc.String())
	}
}

func gateOpen(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*Gate)
	g.lock.Lock()
	if g.active {
		fmt.Println("Removing ", g.s1)
		g.r.Set("open")
		collision.Remove(g.s1)
		g.active = false
	}
	g.lock.Unlock()
	return 0
}

func gateClose(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*Gate)
	g.lock.Lock()
	if !g.active {
		fmt.Println("Adding ", g.s1)
		g.r.Set("closed")
		collision.Add(g.s1)
		g.active = true
	}
	g.lock.Unlock()
	return 0
}
