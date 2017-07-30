package game

import (
	"path/filepath"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/render"
)

type TileType int

const (
	Dood TileType = iota
	Solid
	Logical
	Entity
)

var (
	tileRs    = make(map[Tile]render.Modifiable)
	tileTypes = map[Tile]TileType{
		Sand:               Dood,
		Coral:              Solid,
		PurpleCoralGate:    Logical,
		PurpleCoralGateOff: Logical,
		PurpleCoralSwitch:  Logical,
		BlueCoralGate:      Logical,
		BlueCoralGateOff:   Logical,
		BlueCoralSwitch:    Logical,
		TealCoralGate:      Logical,
		TealCoralGateOff:   Logical,
		TealCoralSwitch:    Logical,
		GreenCoralGate:     Logical,
		GreenCoralGateOff:  Logical,
		GreenCoralSwitch:   Logical,
		Sandtrap:           Logical,
		Sandglob:           Logical,
		Sandgeyser:         Logical,
		VerticalCrab:       Entity,
		HorizontalCrab:     Entity,
		Treasure:           Logical,
		SandKey:            Logical,
		CoralExit:          Logical,
		JeremyTile:         Entity,
	}
	initFunctions map[Tile]func(int, int, render.Renderable)
)

func InitTiles() {
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	tileRs[Sand] = jsh[0][6].Copy()
	tileRs[Coral] = jsh[1][6].Copy()
	tileRs[Sandglob] = jsh[2][6].Copy()
	tileRs[Sandgeyser] = jsh[3][6].Copy()
	tileRs[Treasure] = jsh[1][5].Copy()

	tileRs[PurpleCoralSwitch] = jsh[4][4].Copy()
	tileRs[BlueCoralSwitch] = jsh[5][4].Copy()
	tileRs[TealCoralSwitch] = jsh[6][4].Copy()
	tileRs[GreenCoralSwitch] = jsh[7][4].Copy()

	tileRs[PurpleCoralGate] = render.NewCompound("closed", map[string]render.Modifiable{
		"closed": jsh[4][6].Copy(),
		"open":   jsh[4][5].Copy(),
	})
	tileRs[PurpleCoralGateOff] = render.NewCompound("open", map[string]render.Modifiable{
		"closed": jsh[4][6].Copy(),
		"open":   jsh[4][5].Copy(),
	})
	tileRs[BlueCoralGate] = render.NewCompound("closed", map[string]render.Modifiable{
		"closed": jsh[5][6].Copy(),
		"open":   jsh[5][5].Copy(),
	})
	tileRs[BlueCoralGateOff] = render.NewCompound("open", map[string]render.Modifiable{
		"closed": jsh[5][6].Copy(),
		"open":   jsh[5][5].Copy(),
	})
	tileRs[TealCoralGate] = render.NewCompound("closed", map[string]render.Modifiable{
		"closed": jsh[6][6].Copy(),
		"open":   jsh[6][5].Copy(),
	})
	tileRs[TealCoralGateOff] = render.NewCompound("open", map[string]render.Modifiable{
		"closed": jsh[6][6].Copy(),
		"open":   jsh[6][5].Copy(),
	})
	tileRs[GreenCoralGate] = render.NewCompound("closed", map[string]render.Modifiable{
		"closed": jsh[7][6].Copy(),
		"open":   jsh[7][5].Copy(),
	})
	tileRs[GreenCoralGateOff] = render.NewCompound("open", map[string]render.Modifiable{
		"closed": jsh[7][6].Copy(),
		"open":   jsh[7][5].Copy(),
	})
	tileRs[SandKey] = jsh[2][5].Copy()
	tileRs[CoralExit] = jsh[3][5].Copy()
	tileRs[Sandtrap] = render.NewCompound("hole", map[string]render.Modifiable{
		"hole":   jsh[3][7].Copy(),
		"filled": jsh[3][6].Copy(),
	})
}

func (t Tile) Place(x, y int) {
	xf := float64(x) * 16
	yf := float64(y) * 16
	switch tileTypes[t] {
	case Entity:
		if !levelInit {
			Sand.Place(x, y)
		}
		initFunctions[t](x, y, nil)
	case Solid:
		// Solids are doods that also block movement
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
		layer := 1
		if levelInit {
			// We want to draw placed objects above initially placed objects
			layer = 2
		} else if t != CoralExit {
			// Put sand underneath this if this is the initial placement
			Sand.Place(x, y)
		}
		render.Draw(r, layer)
		initFunctions[t](x, y, r)
	}
}
