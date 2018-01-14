package game

import (
	"path/filepath"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type crab struct {
	// Could make crabs interactive, so they have speed, delta, etc, and could
	// be pushed with mass just like jeremy.
	entities.Solid
}

// Crabs are either vertical or horizontal.
// They follow the player along a single axis.

func newVerticalCrab(x, y int, r render.Renderable) {
	c := new(crab)
	xf, yf := float64(x*16), float64(y*16)
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	r = render.NewSwitch("still", map[string]render.Modifiable{
		"still":  jsh[2][3].Copy(),
		"moving": render.NewSequence(12, jsh[2][3].Copy(), jsh[2][4].Copy()),
	})
	c.Solid = entities.NewSolid(xf+2, yf+2, 12, 12, r, nil, c.Init())
	c.Space.UpdateLabel(blocking)
	render.Draw(c.R, 3)
	c.Bind(vCrabFollow, "EnterFrame")
}

func newHorizontalCrab(x, y int, r render.Renderable) {
	c := new(crab)
	xf, yf := float64(x*16), float64(y*16)
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	r = render.NewSwitch("still", map[string]render.Modifiable{
		"still":  jsh[3][3].Copy(),
		"moving": render.NewSequence(12, jsh[3][3].Copy(), jsh[3][4].Copy()),
	})
	c.Solid = entities.NewSolid(xf+2, yf+2, 12, 12, r, nil, c.Init())
	c.Space.UpdateLabel(blocking)
	render.Draw(c.R, 3)
	c.Bind(hCrabFollow, "EnterFrame")
}

func vCrabFollow(id int, nothing interface{}) int {
	c := event.GetEntity(id).(*entities.Solid)
	if !alg.F64eq(c.Y(), jeremyPos.Y()) {
		delta := jeremyPos.Y() - c.Y()
		if delta > 3.0 {
			delta = 3.0
		} else if delta < -3.0 {
			delta = -3.0
		}
		c.ShiftY(delta)
		if collision.HitLabel(c.Space, blocking, collision.Label(sandtrap)) != nil {
			c.ShiftY(-delta)
		} else {
			c.R.(*render.Switch).Set("moving")
			return 0
		}
	}
	c.R.(*render.Switch).Set("still")
	return 0
}

func hCrabFollow(id int, nothing interface{}) int {
	c := event.GetEntity(id).(*entities.Solid)
	if !alg.F64eq(c.X(), jeremyPos.X()) {
		delta := jeremyPos.X() - c.X()
		if delta > 3.0 {
			delta = 3.0
		} else if delta < -3.0 {
			delta = -3.0
		}
		c.ShiftX(delta)
		if collision.HitLabel(c.Space, blocking, collision.Label(sandtrap)) != nil {
			c.ShiftX(-delta)
		} else {
			c.R.(*render.Switch).Set("moving")
			return 0
		}
	}
	c.R.(*render.Switch).Set("still")
	return 0
}
