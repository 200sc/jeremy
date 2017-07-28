package game

import (
	"fmt"
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
	sand                     int
	dir                      physics.Vector
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
		"still_down":        jsh[0][0].Copy(),
		"still_up":          jsh[0][2].Copy(),
		"still_left":        jsh[0][1].Copy().Modify(render.FlipX),
		"still_right":       jsh[0][1].Copy(),
		"still_down_sand1":  jsh[1][0].Copy(),
		"still_up_sand1":    jsh[1][2].Copy(),
		"still_left_sand1":  jsh[1][1].Copy().Modify(render.FlipX),
		"still_right_sand1": jsh[1][1].Copy(),
		"still_down_sand2":  jsh[2][0].Copy(),
		"still_up_sand2":    jsh[2][2].Copy(),
		"still_left_sand2":  jsh[2][1].Copy().Modify(render.FlipX),
		"still_right_sand2": jsh[2][1].Copy(),
		"still_down_sand3":  jsh[3][0].Copy(),
		"still_up_sand3":    jsh[3][2].Copy(),
		"still_left_sand3":  jsh[3][1].Copy().Modify(render.FlipX),
		"still_right_sand3": jsh[3][1].Copy(),
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
	j.dir = physics.NewVector(0, 0)
	j.SetMass(10)
	collision.Add(j.RSpace.Space)

	// Bindings
	j.Bind(enterJeremy, "EnterFrame")
	j.Bind(consumeSand, "KeyDownE")
	j.Bind(placeGlob, "KeyDownSpacebar")
	j.RSpace.Add(Blocking, jeremyStop)
	return j
}

func placeGlob(id int, nothing interface{}) int {
	j := event.GetEntity(id).(*Jeremy)
	if j.sand > 0 {
		// Todo: right now you can place infinite sand on the same spot. It is tracked, but there's no
		// visual indication how much sand is there.
		x := int((j.X()+8)/16) + int(j.dir.X())
		y := int((j.Y()+8)/16) + int(j.dir.Y())
		fmt.Println(x, y)
		Sandglob.Place(x, y)
		j.sand--
		j.Speed.ShiftX(.1)
		j.Speed.ShiftY(.1)
		j.UpdateAnimation()
	}
	return 0
}

func consumeSand(id int, nothing interface{}) int {
	j := event.GetEntity(id).(*Jeremy)
	s := collision.NewUnassignedSpace((j.X()+8)+16*j.dir.X(), (j.Y()+8)+16*j.dir.Y(), 1, 1)
	hit := collision.HitLabel(
		s,
		collision.Label(Sandglob))
	if hit != nil {
		hit.CID.Trigger("Consume", nil)
		j.sand++
		j.Speed.ShiftX(-.1)
		j.Speed.ShiftY(-.1)
		j.UpdateAnimation()
	}
	fmt.Println("Consumed?", s.GetX(), s.GetY(), j.sand)
	return 0
}

func enterJeremy(id int, frame interface{}) int {
	j := event.GetEntity(id).(*Jeremy)
	j.ApplyFriction(envFriction)
	if oak.IsDown("W") {
		j.Delta.SetY(j.Delta.Y() - j.Speed.Y())
		j.dir.SetPos(0, -1)
	}
	if oak.IsDown("S") {
		j.Delta.SetY(j.Delta.Y() + j.Speed.Y())
		j.dir.SetPos(0, 1)
	}
	if oak.IsDown("A") {
		j.Delta.SetX(j.Delta.X() - j.Speed.X())
		j.dir.SetPos(-1, 0)
	}
	if oak.IsDown("D") {
		j.Delta.SetX(j.Delta.X() + j.Speed.X())
		j.dir.SetPos(1, 0)
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
	if j.Delta.Magnitude() > 0.3 {
		var s string
		if math.Abs(j.Delta.X()) > math.Abs(j.Delta.Y()) {
			if j.Delta.X() < 0 {
				s = "still_left"
			} else {
				s = "still_right"
			}
		} else {
			if j.Delta.Y() < 0 {
				s = "still_up"
			} else {
				s = "still_down"
			}
		}
		cmp.Get(1).(*render.Compound).Set(s)
	}
	s := cmp.Get(1).(*render.Compound).Get()
	s += j.SandString()
	cmp.Get(0).(*render.Compound).Set(s)
}

func (j *Jeremy) SandString() string {
	switch j.sand {
	case 1:
		return "_sand1"
	case 2:
		return "_sand2"
	case 3:
		return "_sand3"
	}
	return ""
}

func jeremyStop(s, s2 *collision.Space) {
	if s.CID != s2.CID {
		id1 := int(s.CID)
		j := event.GetEntity(id1).(*Jeremy)
		j.SetStopXY(s2)
		// Push what we ran into (if we end up having pushable things)
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
