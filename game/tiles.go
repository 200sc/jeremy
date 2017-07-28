package game

import (
	"path/filepath"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type TileType int

const (
	Dood TileType = iota
	Solid
	Logical
)

var (
	tileRs    = make(map[Tile]render.Modifiable)
	tileTypes = map[Tile]TileType{
		Sand:              Dood,
		Coral:             Solid,
		PurpleCoralGate:   Logical,
		PurpleCoralSwitch: Logical,
		BlueCoralGate:     Logical,
		BlueCoralSwitch:   Logical,
		TealCoralGate:     Logical,
		TealCoralSwitch:   Logical,
		Sandtrap:          Logical,
		Sandglob:          Logical,
		Sandgeyser:        Logical,
		Crab:              Logical,
		Treasure:          Logical,
		SandKey:           Logical,
	}
	initFunctions map[Tile]func(int, int, render.Renderable)
)

func InitTiles() {
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	tileRs[Sand] = jsh[0][6].Copy()
	tileRs[Coral] = jsh[1][6].Copy()
	tileRs[Sandglob] = jsh[2][6].Copy()
	tileRs[Sandgeyser] = jsh[3][6].Copy()

	tileRs[PurpleCoralSwitch] = jsh[4][4].Copy()
	tileRs[BlueCoralSwitch] = jsh[5][4].Copy()
	tileRs[TealCoralSwitch] = jsh[6][4].Copy()

	tileRs[PurpleCoralGate] = jsh[4][6].Copy()
	tileRs[BlueCoralSwitch] = jsh[5][6].Copy()
	tileRs[TealCoralSwitch] = jsh[6][6].Copy()
}

func (t Tile) Place(x, y int) {
	xf := float64(x) * 16
	yf := float64(y) * 16
	switch tileTypes[t] {
	case Solid:
		collision.Add(collision.NewLabeledSpace(xf, yf, 16, 16, Blocking))
		fallthrough
	case Dood:
		// Draw this
		r := tileRs[t].Copy()
		r.SetPos(xf, yf)
		render.Draw(r, 0)
	case Logical:
		// Draw this and also do some initalization
		r := tileRs[t].Copy()
		r.SetPos(xf, yf)
		render.Draw(r, 1)
		initFunctions[t](x, y, r)
		// Put sand underneath this if this is the initial placement
		if !levelInit {
			Sand.Place(x, y)
		}
	}
}

func init() {
	initFunctions = map[Tile]func(int, int, render.Renderable){
		Sandglob:          globInit,
		Sandgeyser:        geyserInit,
		PurpleCoralGate:   gateInit(PurpleCoralGate),
		PurpleCoralSwitch: switchInit(PurpleCoralSwitch),
		TealCoralGate:     gateInit(TealCoralGate),
		TealCoralSwitch:   switchInit(TealCoralSwitch),
		BlueCoralGate:     gateInit(BlueCoralGate),
		BlueCoralSwitch:   switchInit(BlueCoralSwitch),
	}
}

type Gate struct {
	physics.Vector
	r  *render.Compound
	s1 *collision.Space
}

type Switch struct {
	collision.Phase
	trigger string
}

func switchInit(t Tile) func(int, int, render.Renderable) {
	return func(x, y int, r render.Renderable) {

	}
}

func gateInit(t Tile) func(int, int, render.Renderable) {
	return func(x, y int, r render.Renderable) {

	}
}
