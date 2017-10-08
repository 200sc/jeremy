package game

import (
	"image/color"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

func initStartButton(x, y int, r render.Renderable) {
	box := render.NewColorBox(100, 50, color.RGBA{100, 150, 255, 255}).Modify(mod.CutRound(.30, .20))
	text := render.DefFont().NewStrText("Start", 40, 30)
	comp := render.NewCompositeR(box, text)
	s := entities.NewSolid(float64(x*16), float64(y*16), 100, 50, comp, 0)
	mouse.Add(s.Space)
	s.Bind(func(int, interface{}) int {
		levelComplete = true
		return 0
	}, "MouseRelease")

	render.Draw(comp, 3)

	title := render.DefFont().NewStrText("Jeremy The Clam", 100, 40)
	render.Draw(title, 3)
}
