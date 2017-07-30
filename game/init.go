package game

import (
	"encoding/csv"
	"log"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render"
)

func FirstSceneInit(string, interface{}) {
	InitTiles()
}

func FirstSceneLoop() bool {
	return false
}

func FirstSceneEnd() (string, *oak.SceneResult) {
	return "level", nil
}

func init() {
	for _, l := range levels {
		f, err := fileutil.Open(filepath.Join("levels", l+".csv"))
		if err != nil {
			log.Fatal(err)
		}
		records, err := csv.NewReader(f).ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		lev := level{}
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
	initFunctions = map[Tile]func(int, int, render.Renderable){
		Sandglob:           globInit,
		Sandgeyser:         geyserInit,
		PurpleCoralGate:    gateInit(PurpleCoralGate),
		PurpleCoralGateOff: offGateInit(PurpleCoralGate),
		PurpleCoralSwitch:  switchInit(PurpleCoralSwitch),
		TealCoralGate:      gateInit(TealCoralGate),
		TealCoralGateOff:   offGateInit(TealCoralGate),
		TealCoralSwitch:    switchInit(TealCoralSwitch),
		BlueCoralGate:      gateInit(BlueCoralGate),
		BlueCoralGateOff:   offGateInit(BlueCoralGate),
		BlueCoralSwitch:    switchInit(BlueCoralSwitch),
		GreenCoralGate:     gateInit(GreenCoralGate),
		GreenCoralGateOff:  offGateInit(GreenCoralGate),
		GreenCoralSwitch:   alternatingSwitchInit(GreenCoralSwitch),
		CoralExit:          exitInit,
		SandKey:            keyInit,
		Sandtrap:           trapInit,
		Treasure:           treasureInit,
		JeremyTile:         NewJeremy,
		HorizontalCrab:     NewHorizontalCrab,
		VerticalCrab:       NewVerticalCrab,
	}
}
