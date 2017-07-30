package game

import (
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type treasureBox struct {
	r render.Renderable
	s *collision.Space
	event.CID
}

func (t *treasureBox) Init() event.CID {
	t.CID = event.NextID(t)
	return t.CID
}

func treasureInit(x, y int, r render.Renderable) {
	//todo: shiny particles
	//todo: a reward for the treasure you pick up
	xf, yf := float64(x)*16, float64(y)*16
	t := new(treasureBox)
	t.r = r
	t.s = collision.NewFullSpace(xf+2, yf+2, 14, 14, collision.Label(treasure), t.Init())
	collision.Add(t.s)
	t.Bind(treasureDestroy, "Consume")
}

func treasureDestroy(id int, nothing interface{}) int {
	audio.Play(sounds, "Treasure.wav")
	t := event.GetEntity(id).(*treasureBox)
	t.r.UnDraw()
	collision.Remove(t.s)
	event.DestroyEntity(id)
	return 0
}
