package game

import (
	"log"
	"math"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var (
	jeremyPos physics.Vector
)

type jeremy struct {
	entities.Interactive
	physics.Mass
	eyes                     *render.Switch
	stopMovingX, stopMovingY bool
	overlap                  physics.Vector
	sand                     int
	dir                      physics.Vector
}

func (j *jeremy) Init() event.CID {
	j.CID = event.NextID(j)
	return j.CID
}

func newJeremy(x, y int, r render.Renderable) {
	j := new(jeremy)

	// Renderable setup
	// Jeremy's main sprite
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	cmp := render.NewSwitch("still_down", map[string]render.Modifiable{
		"still_down":        jsh[0][0].Copy(),
		"still_up":          jsh[0][2].Copy(),
		"still_left":        jsh[0][1].Copy().Modify(mod.FlipX),
		"still_right":       jsh[0][1].Copy(),
		"still_down_sand1":  jsh[1][0].Copy(),
		"still_up_sand1":    jsh[1][2].Copy(),
		"still_left_sand1":  jsh[1][1].Copy().Modify(mod.FlipX),
		"still_right_sand1": jsh[1][1].Copy(),
		"still_down_sand2":  jsh[2][0].Copy(),
		"still_up_sand2":    jsh[2][2].Copy(),
		"still_left_sand2":  jsh[2][1].Copy().Modify(mod.FlipX),
		"still_right_sand2": jsh[2][1].Copy(),
		"still_down_sand3":  jsh[3][0].Copy(),
		"still_up_sand3":    jsh[3][2].Copy(),
		"still_left_sand3":  jsh[3][1].Copy().Modify(mod.FlipX),
		"still_right_sand3": jsh[3][1].Copy(),
		"still_down_key":    jsh[4][0].Copy(),
		"still_up_key":      jsh[4][2].Copy(),
		"still_left_key":    jsh[4][1].Copy().Modify(mod.FlipX),
		"still_right_key":   jsh[4][1].Copy(),
	})
	// Jeremy's eyes
	eyes, err := render.LoadSheetSequence(filepath.Join("3", "eyes.png"), 3, 3, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 2, 0, 1, 0)
	if err != nil {
		log.Fatal(err)
	}
	//
	downeyes := render.NewComposite(eyes.Copy(), eyes.Copy().Modify(mod.FlipX))
	downeyes.Get(0).ShiftX(4)
	downeyes.Get(0).ShiftY(8)
	downeyes.Get(1).ShiftX(9)
	downeyes.Get(1).ShiftY(8)
	//
	lefteyes := eyes.Copy().Modify(mod.FlipX)
	lefteyes.ShiftX(1)
	lefteyes.ShiftY(8)
	//
	righteyes := eyes.Copy()
	righteyes.ShiftX(12)
	righteyes.ShiftY(8)
	//
	eyecmp := render.NewSwitch("still_down", map[string]render.Modifiable{
		"still_down":  downeyes,
		"still_up":    render.EmptyRenderable(),
		"still_left":  lefteyes,
		"still_right": righteyes,
	})
	//eyecmp.SetOffsets("still_left", physics.NewVector(1, 8))
	//eyecmp.SetOffsets("still_right", physics.NewVector(12, 8))
	//eyecmp.ShiftX(-14)
	//eyecmp.ShiftY(-24)
	// Draw both of those at the same time
	composite := render.NewComposite(cmp, eyecmp)

	// Non-renderable variables
	j.Interactive = entities.NewInteractive(0, 0, 14, 14, composite, nil, j.Init(), 0.4)
	j.Speed = physics.NewVector(.5, .5)
	j.overlap = physics.NewVector(0, 0)
	j.dir = physics.NewVector(0, 0)
	j.SetMass(10)
	j.RSpace.Space.UpdateLabel(collision.Label(jeremyTile))

	// Bindings
	// Do this every frame
	j.Bind(enterJeremy, "EnterFrame")
	// Do this when E is pressed
	j.Bind(consumeSand, "KeyDownE")
	// Do this when Spacebar is pressed
	j.Bind(placeGlob, "KeyDownSpacebar")
	// Do this when something triggers "pause" on Jeremy
	j.Bind(pauseJeremy, "PausePlayer")

	// React to hitting certain collision spaces with these bindings:
	j.RSpace.Add(blocking, jeremyStop)
	j.RSpace.Add(collision.Label(sandtrap), jeremyStop)
	j.RSpace.Add(collision.Label(treasure), pickUpTreasure)

	// Draw yourself and place yourself according to the input x,y
	render.Draw(j.R, 3)
	j.SetPos(float64(x*16), float64(y*16))

	// Track a global so crabs know where you are and can follow you
	jeremyPos = j.Vector
}

// When Jeremy is paused, they no longer respond to bindings.
func pauseJeremy(id int, nothing interface{}) int {
	j := event.GetEntity(id).(*jeremy)
	j.UnbindAll()
	return 0
}

// When space is pressed, Jeremy tries to spit out sand in front of themselves.
func placeGlob(id int, nothing interface{}) int {
	j := event.GetEntity(id).(*jeremy)
	x := int((j.X()+8)/16) + int(j.dir.X())
	y := int((j.Y()+8)/16) + int(j.dir.Y())
	s := collision.NewUnassignedSpace(float64(x*16)+2, float64(y*16)+2, 12, 12)
	// If there's something in the way, try to place a little farther
	if collision.HitLabel(s, blocking, collision.Label(jeremyTile)) != nil {
		collision.ShiftSpace(j.dir.X(), j.dir.Y(), s)
		if collision.HitLabel(s, blocking, collision.Label(jeremyTile)) != nil {
			return 0
		}
		x += int(j.dir.X())
		y += int(j.dir.Y())
	}

	// If Jeremy has a key right now, drop a key.
	if j.sand == 4 {
		sandKey.Place(x, y)
		j.sand = 0
		j.Speed.ShiftX(.3)
		j.Speed.ShiftY(.3)
		j.UpdateAnimation()
		audio.Play(sounds, "Sand.wav")
		// Otherwise drop sand
	} else if j.sand > 0 {
		sandglob.Place(x, y)
		j.sand--
		j.Speed.ShiftX(0.1)
		j.Speed.ShiftY(0.1)
		j.UpdateAnimation()
		audio.Play(sounds, "Sand.wav")
	}
	return 0
}

// When E is pressed, jeremy tries to suck up the sand in front of them
func consumeSand(id int, nothing interface{}) int {
	j := event.GetEntity(id).(*jeremy)
	s := collision.NewUnassignedSpace((j.X()+8)+16*j.dir.X(), (j.Y()+8)+16*j.dir.Y(), 1, 1)
	// You can pick up keys if you don't have any other sand
	if j.sand == 0 {
		hit := collision.HitLabel(
			s,
			collision.Label(sandKey),
		)
		if hit != nil {
			audio.Play(sounds, "Sand.wav")
			hit.CID.Trigger("Consume", nil)
			j.sand = 4
			j.Speed.ShiftX(-.3)
			j.Speed.ShiftY(-.3)
			j.UpdateAnimation()
			return 0
		}
	}
	// Otherwise you can pick up up to three sand globs
	if j.sand < 3 {
		hit := collision.HitLabel(
			s,
			collision.Label(sandglob),
		)
		if hit != nil {
			audio.Play(sounds, "Sand.wav")
			hit.CID.Trigger("Consume", nil)
			j.sand++
			j.Speed.ShiftX(-.1)
			j.Speed.ShiftY(-.1)
			j.UpdateAnimation()
		}
	}
	return 0
}

func enterJeremy(id int, frame interface{}) int {
	j := event.GetEntity(id).(*jeremy)
	// Slow down whatever Jeremy's old delta to move was
	j.ApplyFriction(envFriction)

	// Increase Jeremy's delta and set their direction according to what keys are
	// pressed.
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

	// Handle reactive collision
	<-j.RSpace.CallOnHits()
	// If something hit jeremy, determine how much Jeremy should move
	// back by to get out of it.
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

func (j *jeremy) UpdateAnimation() {
	cmp := j.R.(*render.Composite)
	if j.Delta.Magnitude() > 0.4 {
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
		// If we might have changed directions, update the eyes
		cmp.Get(1).(*render.Switch).Set(s)
	}
	s := cmp.Get(1).(*render.Switch).Get()
	s += j.SandString()
	// Always update the sand level that jeremy has consumed
	cmp.Get(0).(*render.Switch).Set(s)
}

// This converts how much sand jeremy has into what string that
// represents on their animation compound.
func (j *jeremy) SandString() string {
	switch j.sand {
	case 1:
		return "_sand1"
	case 2:
		return "_sand2"
	case 3:
		return "_sand3"
	case 4:
		return "_key"
	}
	return ""
}

func pickUpTreasure(s, s2 *collision.Space) {
	s2.CID.Trigger("Consume", nil)
}

func jeremyStop(s, s2 *collision.Space) {
	if s.CID != s2.CID {
		id1 := int(s.CID)
		j := event.GetEntity(id1).(*jeremy)
		j.SetStopXY(s2)
	}
}

func (j *jeremy) SetStopXY(s *collision.Space) {
	xOver, yOver := j.RSpace.Space.Overlap(s)

	if math.Abs(xOver) < math.Abs(yOver) {
		j.stopMovingX = true
		j.overlap.SetX(xOver)
	} else {
		j.stopMovingY = true
		j.overlap.SetY(yOver)
	}
}
