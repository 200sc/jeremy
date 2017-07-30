package game

import (
	"sync"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type gateColor int

const (
	purple gateColor = iota
	teal
	blue
	green
	nonGate
)

func (t tile) gateColor() gateColor {
	switch t {
	case purpleCoralGate, purpleCoralSwitch:
		return purple
	case tealCoralGate, tealCoralSwitch:
		return teal
	case blueCoralGate, blueCoralSwitch:
		return blue
	case greenCoralGate, greenCoralSwitch:
		return green
	}
	return nonGate
}

// Gate colors can be converted to strings so switches can trigger the
// opening / closing of the appropriate gates when pressed
func (gc gateColor) String() string {
	switch gc {
	case purple:
		return "Purple"
	case blue:
		return "Blue"
	case teal:
		return "Teal"
	case green:
		return "Green"
	}
	return ""
}

type gate struct {
	physics.Vector
	r  *render.Compound
	s1 *collision.Space
	event.CID
	lock   sync.Mutex
	active bool
}

func (g *gate) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func newGate(x, y float64, r render.Renderable) *gate {
	g := new(gate)
	g.lock = sync.Mutex{}
	g.r = r.(*render.Compound)
	g.s1 = collision.NewFullSpace(x+2, y+2, 12, 12, blocking, g.Init())
	return g
}

func gateInit(t tile) func(int, int, render.Renderable) {
	gc := t.gateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		g := newGate(xf, yf, r)
		collision.Add(g.s1)
		g.active = true
		g.Bind(gateOpen, "Open"+gc.String())
		g.Bind(gateClose, "Close"+gc.String())
	}
}

func offGateInit(t tile) func(int, int, render.Renderable) {
	gc := t.gateColor()
	return func(x, y int, r render.Renderable) {
		xf, yf := float64(x)*16, float64(y)*16
		g := newGate(xf, yf, r)
		g.active = false
		g.Bind(gateOpen, "Close"+gc.String())
		g.Bind(gateClose, "Open"+gc.String())
	}
}

func gateOpen(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*gate)
	g.lock.Lock()
	if g.active {
		g.r.Set("open")
		collision.Remove(g.s1)
		g.active = false
	}
	g.lock.Unlock()
	return 0
}

func gateClose(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*gate)
	// We lock and boolean check this to guarantee that gates don't add their
	// collision space multiple times
	g.lock.Lock()
	if !g.active {
		g.r.Set("closed")
		collision.Add(g.s1)
		g.active = true
	}
	g.lock.Unlock()
	return 0
}
