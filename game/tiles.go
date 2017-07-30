package game

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/render"
)

type tileType int

const (
	dood tileType = iota
	solid
	logical
	entity
)

var (
	tileRs    = make(map[tile]render.Modifiable)
	tileTypes = map[tile]tileType{
		sand:               dood,
		coral:              solid,
		purpleCoralGate:    logical,
		purpleCoralGateOff: logical,
		purpleCoralSwitch:  logical,
		blueCoralGate:      logical,
		blueCoralGateOff:   logical,
		blueCoralSwitch:    logical,
		tealCoralGate:      logical,
		tealCoralGateOff:   logical,
		tealCoralSwitch:    logical,
		greenCoralGate:     logical,
		greenCoralGateOff:  logical,
		greenCoralSwitch:   logical,
		sandtrap:           logical,
		sandglob:           logical,
		sandgeyser:         logical,
		verticalCrab:       entity,
		horizontalCrab:     entity,
		treasure:           logical,
		sandKey:            logical,
		coralExit:          logical,
		jeremyTile:         entity,
	}
	initFunctions map[tile]func(int, int, render.Renderable)
)

func (t tile) Place(x, y int) {
	xf := float64(x) * 16
	yf := float64(y) * 16
	layer := 0
	switch tileTypes[t] {
	case entity:
		if !levelInit {
			sand.Place(x, y)
		}
		// Entities are expected to do everything themselves
		initFunctions[t](x, y, nil)
	case solid:
		// Solids are doods that also block movement
		if !levelInit {
			sand.Place(x, y)
		}
		layer++
		collision.Add(collision.NewLabeledSpace(xf, yf, 16, 16, blocking))
		fallthrough
	case dood:
		// Draw this
		r := tileRs[t].Copy()
		r.SetPos(xf, yf)
		render.Draw(r, layer)
	case logical:
		// Draw this and also do some initalization
		r := tileRs[t].Copy()
		r.SetPos(xf, yf)
		layer++
		if levelInit {
			// We want to draw placed objects above initially placed objects
			layer = 2
		} else if t != coralExit {
			// Put sand underneath this if this is the initial placement
			sand.Place(x, y)
		}
		render.Draw(r, layer)
		initFunctions[t](x, y, r)
	}
}
