package game

import (
	"image/color"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
)

type trap struct {
	r *render.Compound
	event.CID
	s         *collision.Space
	ps        *particle.Source
	emptyTime time.Time
}

func (t *trap) Init() event.CID {
	t.CID = event.NextID(t)
	return t.CID
}

func trapInit(x, y int, r render.Renderable) {
	xf, yf := float64(x)*16, float64(y)*16
	t := new(trap)
	t.r = r.(*render.Compound)
	t.s = collision.NewFullSpace(xf, yf, 16, 16, collision.Label(sandtrap), t.Init())
	collision.Add(t.s)
	t.Bind(trapFill, "UseGlob")
	t.Bind(trapEnter, "EnterFrame")
}

func trapEnter(id int, nothing interface{}) int {
	t := event.GetEntity(id).(*trap)
	if t.r.Get() == "filled" && t.emptyTime.Before(time.Now()) {
		audio.Play(sounds, "PitEmpty.wav")
		t.r.Set("hole")
		collision.Add(t.s)
		t.ps.Stop()
	}
	return 0
}

func trapFill(id int, nothing interface{}) int {
	t := event.GetEntity(id).(*trap)
	t.r.Set("filled")
	collision.Remove(t.s)
	// "Falling sand" particles
	t.ps = particle.NewColorGenerator(
		particle.NewPerFrame(floatrange.NewLinear(2, 4)),
		particle.Pos(t.r.GetX()+8, t.r.GetY()+8),
		particle.LifeSpan(floatrange.NewLinear(7, 12)),
		particle.Angle(floatrange.Constant(270)),
		particle.Speed(floatrange.NewLinear(.1, .4)),
		particle.Spread(5, 2),
		particle.Color(color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0}),
		particle.Size(intrange.NewLinear(1, 2)),
		particle.Shape(shape.Square),
	).Generate(2)
	// Once we pass this time, reset the hole to need to be filled again
	t.emptyTime = time.Now().Add(3 * time.Second)
	return 0
}
