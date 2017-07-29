package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type Key struct {
	Glob
}

func (k *Key) Init() event.CID {
	k.CID = event.NextID(k)
	return k.CID
}

func keyInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	k := new(Key)
	k.Vector = physics.NewVector(xf, yf)
	k.r = r
	k.Init()
	// Todo: this is duplicating a good amount of Glob
	k.s1 = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(SandKey), k.CID)
	k.s2 = collision.NewFullSpace(xf, yf, 16, 16, Blocking, k.CID)
	collision.Add(k.s1, k.s2)
	k.CID.Bind(keyDestroy, "Consume")
	// Only differentiating thing: when a key is dropped, if its on an exit,
	// complete the level
	hit := collision.HitLabel(k.s1, collision.Label(CoralExit))
	if hit != nil {
		hit.CID.Trigger("OpenExit", nil)
		r.UnDraw()
	}
}

func keyDestroy(id int, nothing interface{}) int {
	k := event.GetEntity(id).(*Key)
	k.r.UnDraw()
	collision.Remove(k.s1, k.s2)
	event.DestroyEntity(id)
	return 0
}
