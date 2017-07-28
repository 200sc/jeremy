package game

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render"
)

var (
	envFriction = 0.8
)

func LevelStart(prevScene string, data interface{}) {
	// Move this
	InitTiles()

	j := NewJeremy()
	render.Draw(j.R, 1)
	// Get name from data
	levelName := "level1"
	// Get player position from level
	j.SetPos(120, 120)
	l := levelStore[levelName]
	l.Place()
}

func LevelLoop() bool {
	return true
}

func LevelEnd() (string, *oak.SceneResult) {
	return "level", nil
}

type level [20][15]Tile

func (l level) Place() {
	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {
			l[x][y].Place(x, y)
		}
	}
}

var (
	levelStore = make(map[string]level)
	levels     = []string{"level1"}
)

func init() {
	for _, l := range levels {
		f, err := fileutil.Open(l + ".csv")
		if err != nil {
			log.Fatal(err)
		}
		records, err := csv.NewReader(f).ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		lev := level{}
		fmt.Println(len(records), len(records[0]))
		for x := 0; x < 20; x++ {
			for y := 0; y < 15; y++ {
				v := records[y][x]
				t, err := strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
				lev[x][y] = Tile(t)
			}
		}
		levelStore[l] = lev
	}
	fmt.Println(levelStore)
}
