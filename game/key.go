package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type key struct {
	glob
}

func (k *key) Init() event.CID {
	k.CID = event.NextID(k)
	return k.CID
}

func keyInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	k := new(key)
	k.Vector = physics.NewVector(xf, yf)
	k.r = r
	k.Init()
	// Todo: this is duplicating a good amount of Glob
	k.s1 = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(sandKey), k.CID)
	k.s2 = collision.NewFullSpace(xf, yf, 16, 16, blocking, k.CID)
	collision.Add(k.s1, k.s2)
	k.CID.Bind(keyDestroy, "Consume")
	// Only differentiating thing: when a key is dropped, if its on an exit,
	// complete the level
	hit := collision.HitLabel(k.s1, collision.Label(coralExit))
	if hit != nil {
		hit.CID.Trigger("OpenExit", nil)
		r.Undraw()
	}
}

func keyDestroy(id int, nothing interface{}) int {
	k := event.GetEntity(id).(*key)
	k.r.Undraw()
	collision.Remove(k.s1, k.s2)
	event.DestroyEntity(id)
	return 0
}
