package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type glob struct {
	physics.Vector
	r      render.Renderable
	s1, s2 *collision.Space
	event.CID
}

func (g *glob) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func globInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	g := new(glob)
	g.Vector = physics.NewVector(xf, yf)
	g.r = r
	g.Init()
	g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, collision.Label(sandglob), g.CID)
	g.s2 = collision.NewFullSpace(xf+2, yf+2, 12, 12, blocking, g.CID)
	hit := collision.HitLabel(g.s1, collision.Label(sandtrap))
	collision.Add(g.s1, g.s2)
	// When the player picks this up, they'll trigger Consume on it
	g.CID.Bind(globDestroy, "Consume")
	// If this glob actually hit something when it was placed, then this should
	// immediately go away and if that thing has an interaction with globs,
	// trigger it.
	if hit != nil {
		hit.CID.Trigger("UseGlob", nil)
		g.CID.Trigger("Consume", nil)
	}
}

func globDestroy(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*glob)
	g.r.UnDraw()
	collision.Remove(g.s1, g.s2)
	event.DestroyEntity(id)
	return 0
}
