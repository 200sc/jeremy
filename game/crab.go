package game

import (
	"path/filepath"

	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type Crab struct {
	entities.Solid
}

func NewVerticalCrab(x, y int, r render.Renderable) {
	c := new(Crab)
	xf, yf := float64(x*16), float64(y*16)
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	r = render.NewCompound("still", map[string]render.Modifiable{
		"still":  jsh[2][3].Copy(),
		"moving": render.NewSequence([]render.Modifiable{jsh[2][3].Copy(), jsh[2][4].Copy()}, 12),
	})
	c.Solid = entities.NewSolid(xf+2, yf+2, 12, 12, r, c.Init())
	c.Space.UpdateLabel(Blocking)
	render.Draw(c.R, 3)
	c.Bind(vCrabFollow, "EnterFrame")
}

func NewHorizontalCrab(x, y int, r render.Renderable) {
	c := new(Crab)
	xf, yf := float64(x*16), float64(y*16)
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	r = render.NewCompound("still", map[string]render.Modifiable{
		"still":  jsh[3][3].Copy(),
		"moving": render.NewSequence([]render.Modifiable{jsh[3][3].Copy(), jsh[3][4].Copy()}, 12),
	})
	c.Solid = entities.NewSolid(xf+2, yf+2, 12, 12, r, c.Init())
	c.Space.UpdateLabel(Blocking)
	render.Draw(c.R, 3)
	c.Bind(hCrabFollow, "EnterFrame")
}

func vCrabFollow(id int, nothing interface{}) int {
	c := event.GetEntity(id).(*entities.Solid)
	if !alg.F64eq(c.Y(), JeremyPos.Y()) {
		delta := JeremyPos.Y() - c.Y()
		if delta > 2.0 {
			delta = 2.0
		} else if delta < -2.0 {
			delta = -2.0
		}
		c.ShiftY(delta)
		if collision.HitLabel(c.Space, Blocking, collision.Label(Sandtrap)) != nil {
			c.ShiftY(-delta)
		} else {
			c.R.(*render.Compound).Set("moving")
			return 0
		}
	}
	c.R.(*render.Compound).Set("still")
	return 0
}

func hCrabFollow(id int, nothing interface{}) int {
	c := event.GetEntity(id).(*entities.Solid)
	if !alg.F64eq(c.X(), JeremyPos.X()) {
		delta := JeremyPos.X() - c.X()
		if delta > 2.0 {
			delta = 2.0
		} else if delta < -2.0 {
			delta = -2.0
		}
		c.ShiftX(delta)
		if collision.HitLabel(c.Space, Blocking, collision.Label(Sandtrap)) != nil {
			c.ShiftX(-delta)
		} else {
			c.R.(*render.Compound).Set("moving")
			return 0
		}
	}
	c.R.(*render.Compound).Set("still")
	return 0
}
