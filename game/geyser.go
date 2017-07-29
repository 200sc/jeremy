package game

import (
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
)

func geyserInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
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
