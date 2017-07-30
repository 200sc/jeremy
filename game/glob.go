package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type Glob struct {
	physics.Vector
	r      render.Renderable
	s1, s2 *collision.Space
	event.CID
}

func (g *Glob) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func globInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	g := new(Glob)
	g.Vector = physics.NewVector(xf, yf)
	g.r = r
	g.Init()
	g.s1 = collision.NewFullSpace(xf+2, yf+2, 12, 12, collision.Label(Sandglob), g.CID)
	g.s2 = collision.NewFullSpace(xf+2, yf+2, 12, 12, Blocking, g.CID)
	hit := collision.HitLabel(g.s1, collision.Label(Sandtrap))
	collision.Add(g.s1, g.s2)
	g.CID.Bind(globDestroy, "Consume")
	if hit != nil {
		hit.CID.Trigger("UseGlob", nil)
		g.CID.Trigger("Consume", nil)
	}
}

func globDestroy(id int, nothing interface{}) int {
	g := event.GetEntity(id).(*Glob)
	g.r.UnDraw()
	collision.Remove(g.s1, g.s2)
	event.DestroyEntity(id)
	return 0
}
