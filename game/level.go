package game

type level [20][15]tile

func (l level) Place() {
	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {
			l[x][y].Place(x, y)
		}
	}
}

var (
	levelStore = make(map[string]level)
	levels     = []string{"menu",
		"tutorial1", "tutorial2", "tutorial3", "tutorial4", "tutorial5", "tutorial6",
		"level1", "level2", "level3", "level4", "level5", "level6", "level7", "level8", "level9"}
	//"victory"}
)
