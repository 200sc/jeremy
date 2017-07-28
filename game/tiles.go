package game

import (
	"image/color"
	"path/filepath"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
)

// Tile enum
const (
	_ Tile = iota
	_
	_
	_
	_
	_
	_
	_
	_
	_
	// We start at 10, so the csvs can be more easily read (evenly spaced columns)
	// Subtract 21 from line number for key
	Sand
	Coral
	PurpleCoralGate
	PurpleCoralSwitch
	BlueCoralGate
	BlueCoralSwitch
	TealCoralGate
	TealCoralSwitch
	Sandtrap
	Sandglob
	Sandgeyser
	Crab
	Treasure
	SandKey
)

type Tile int

var (
	tileRs = make(map[Tile]render.Modifiable)
)

func InitTiles() {
	jsh := render.GetSheet(filepath.Join("16", "jeremy.png"))
	tileRs[Sand] = jsh[0][6].Copy()
	tileRs[Coral] = jsh[1][6].Copy()
	tileRs[Sandglob] = jsh[2][6].Copy()
	tileRs[Sandgeyser] = jsh[3][6].Copy()
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
		render.Draw(r, 0)
		initFunctions[t](x, y)
	}
}

type TileType int

const (
	Dood TileType = iota
	Solid
	Logical
)

var (
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
	initFunctions = map[Tile]func(int, int){
		Sandglob:   globInit,
		Sandgeyser: geyserInit,
	}
)

type Glob struct {
	s1, s2 *collision.Space
	event.CID
}

func (g *Glob) Init() event.CID {
	g.CID = event.NextID(g)
	return g.CID
}

func globInit(x, y int) {
	xf := float64(x) * 16
	yf := float64(y) * 16
	g := new(Glob)
	g.Init()
	g.s1 = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(Sandglob), g.CID)
	g.s2 = collision.NewFullSpace(xf, yf, 16, 16, Blocking, g.CID)
	collision.Add(g.s1, g.s2)
}

func geyserInit(x, y int) {
	xf := float64(x) * 16
	yf := float64(y) * 16
	// Set up a sandglob collision space, but don't bind anything to happen on consume
	collision.Add(collision.NewLabeledSpace(xf, yf, 16, 16, collision.Label(Sandglob)))
	collision.Add(collision.NewLabeledSpace(xf, yf, 16, 16, Blocking))
	// Make particles
	particle.NewColorGenerator(
		particle.NewPerFrame(floatrange.NewLinear(2, 4)),
		particle.Pos(xf+8, yf+8),
		particle.LifeSpan(floatrange.NewLinear(7, 12)),
		particle.Angle(floatrange.Constant(90)),
		particle.Speed(floatrange.NewLinear(.1, .4)),
		particle.Spread(5, 0),
		particle.Color(color.RGBA{127, 201, 255, 255}, color.RGBA{0, 0, 0, 0}, color.RGBA{127, 201, 255, 255}, color.RGBA{0, 0, 0, 0}),
		particle.Size(intrange.NewLinear(1, 2)),
		particle.Shape(shape.Square),
	).Generate(2)
}
