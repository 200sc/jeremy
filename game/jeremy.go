package game

import (
	"log"
	"math"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type Jeremy struct {
	entities.Interactive
	physics.Mass
	eyes                     *render.Compound
	stopMovingX, stopMovingY bool
	overlap                  physics.Vector
}

func (j *Jeremy) Init() event.CID {
	j.CID = event.NextID(j)
	return j.CID
}

func NewJeremy() *Jeremy {
	j := new(Jeremy)

	// Renderable setup
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	cmp := render.NewCompound("still_down", map[string]render.Modifiable{
		"still_down":  jsh[0][0].Copy(),
		"still_up":    jsh[0][2].Copy(),
		"still_left":  jsh[0][1].Copy().Modify(render.FlipX),
		"still_right": jsh[0][1].Copy(),
	})
	eyes, err := render.LoadSheetAnimation(filepath.Join("3", "eyes.png"), 3, 3, 0, 6, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, 0, 1, 0})
	if err != nil {
		log.Fatal(err)
	}
	//
	downeyes := render.NewComposite([]render.Modifiable{eyes.Copy(), eyes.Copy().Modify(render.FlipX)})
	downeyes.Get(0).ShiftX(4)
	downeyes.Get(0).ShiftY(8)
	downeyes.Get(1).ShiftX(9)
	downeyes.Get(1).ShiftY(8)
	//
	lefteyes := eyes.Copy().Modify(render.FlipX)
	lefteyes.ShiftX(1)
	lefteyes.ShiftY(8)
	//
	righteyes := eyes.Copy()
	righteyes.ShiftX(12)
	righteyes.ShiftY(8)
	//
	eyecmp := render.NewCompound("still_down", map[string]render.Modifiable{
		"still_down":  downeyes,
		"still_up":    render.EmptyRenderable(),
		"still_left":  lefteyes,
		"still_right": righteyes,
	})
	composite := render.NewComposite([]render.Modifiable{cmp, eyecmp})

	// Non-renderable variables
	j.Interactive = entities.NewInteractive(0, 0, 16, 16, composite, j.Init(), 0.4)
	j.Speed = physics.NewVector(.5, .5)
	j.overlap = physics.NewVector(0, 0)
	j.SetMass(10)
	collision.Add(j.RSpace.Space)

	// Bindings
	j.Bind(enterJeremy, "EnterFrame")
	j.RSpace.Add(Blocking, jeremyStop)
	return j
}

func enterJeremy(id int, frame interface{}) int {
	j := event.GetEntity(id).(*Jeremy)
	j.ApplyFriction(envFriction)
	if oak.IsDown("W") {
		j.Delta.SetY(j.Delta.Y() - j.Speed.Y())
	}
	if oak.IsDown("S") {
		j.Delta.SetY(j.Delta.Y() + j.Speed.Y())
	}
	if oak.IsDown("A") {
		j.Delta.SetX(j.Delta.X() - j.Speed.X())
	}
	if oak.IsDown("D") {
		j.Delta.SetX(j.Delta.X() + j.Speed.X())
	}
	j.ShiftPos(j.Delta.X(), j.Delta.Y())
	j.UpdateAnimation()

	// Handle collision (with blocking things, jeremy doesn't collide with anything else)
	<-j.RSpace.CallOnHits()
	if j.stopMovingX || j.stopMovingY {
		v := j.overlap.Copy().Scale(-1)
		if math.Abs(v.X()) > math.Abs(j.Delta.X()) {
			j.Delta.SetX(0)
		} else {
			j.Delta.SetX(j.Delta.X() + v.X())
		}
		if math.Abs(v.Y()) > math.Abs(j.Delta.Y()) {
			j.Delta.SetY(0)
		} else {
			j.Delta.SetY(j.Delta.Y() + v.Y())
		}
		j.Delta.Add(v)
		j.ShiftPos(j.Delta.X(), j.Delta.Y())
	}
	j.stopMovingX = false
	j.stopMovingY = false
	j.overlap.SetPos(0, 0)
	return 0
}

func (j *Jeremy) UpdateAnimation() {
	// Todo: make this composite setting easier
	cmp := j.R.(*render.Composite)
	if j.Delta.Magnitude() < 0.3 {
		return
	}
	if math.Abs(j.Delta.X()) > math.Abs(j.Delta.Y()) {
		if j.Delta.X() < 0 {
			for i := 0; i < 2; i++ {
				cmp.Get(i).(*render.Compound).Set("still_left")
			}
		} else {
			for i := 0; i < 2; i++ {
				cmp.Get(i).(*render.Compound).Set("still_right")
			}
		}
	} else {
		if j.Delta.Y() < 0 {
			for i := 0; i < 2; i++ {
				cmp.Get(i).(*render.Compound).Set("still_up")
			}
		} else {
			for i := 0; i < 2; i++ {
				cmp.Get(i).(*render.Compound).Set("still_down")
			}
		}
	}
}

func jeremyStop(s, s2 *collision.Space) {
	if s.CID != s2.CID {
		id1 := int(s.CID)
		j := event.GetEntity(id1).(*Jeremy)
		j.SetStopXY(s2)
		// Push what we ran into
		v := j.Delta.Copy()
		if j.stopMovingX {
			v.SetY(0)
		} else if j.stopMovingY {
			v.SetX(0)
		}
		s2.CID.Trigger("push", physics.DefaultForceVector(v, j.GetMass()))
	}
}

func (j *Jeremy) SetStopXY(s *collision.Space) {
	xOver, yOver := j.RSpace.Space.Overlap(s)

	if math.Abs(xOver) < math.Abs(yOver) {
		j.stopMovingX = true
		j.overlap.SetX(xOver)
	} else {
		j.stopMovingY = true
		j.overlap.SetY(yOver)
	}
}
