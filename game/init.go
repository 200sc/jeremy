package game

import (
	"encoding/csv"
	"log"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak/fileutil"
	"github.com/oakmound/oak/render"
)

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
				lev[x][y] = tile(t)
			}
		}
		levelStore[l] = lev
	}
	initFunctions = map[tile]func(int, int, render.Renderable){
		sandglob:           globInit,
		sandgeyser:         geyserInit,
		purpleCoralGate:    gateInit(purpleCoralGate),
		purpleCoralGateOff: offGateInit(purpleCoralGate),
		purpleCoralSwitch:  switchInit(purpleCoralSwitch),
		tealCoralGate:      gateInit(tealCoralGate),
		tealCoralGateOff:   offGateInit(tealCoralGate),
		tealCoralSwitch:    switchInit(tealCoralSwitch),
		blueCoralGate:      gateInit(blueCoralGate),
		blueCoralGateOff:   offGateInit(blueCoralGate),
		blueCoralSwitch:    switchInit(blueCoralSwitch),
		greenCoralGate:     gateInit(greenCoralGate),
		greenCoralGateOff:  offGateInit(greenCoralGate),
		greenCoralSwitch:   alternatingSwitchInit(greenCoralSwitch),
		coralExit:          exitInit,
		sandKey:            keyInit,
		sandtrap:           trapInit,
		treasure:           treasureInit,
		jeremyTile:         newJeremy,
		horizontalCrab:     newHorizontalCrab,
		verticalCrab:       newVerticalCrab,
		startButton:        initStartButton,
	}
}

// initTiles is not in the init() because sheets cannot be obtained from oak
// prior to oak's startup
func initTiles() {
	jshtt, _ := render.GetSheet(filepath.Join("16", "jeremy.png"))
	jsh := jshtt.ToSprites()
	tileRs[sand] = jsh[0][6].Copy()
	tileRs[coral] = jsh[1][6].Copy()
	tileRs[sandglob] = jsh[2][6].Copy()
	tileRs[sandgeyser] = jsh[3][6].Copy()
	tileRs[treasure] = jsh[1][5].Copy()

	tileRs[purpleCoralSwitch] = jsh[4][4].Copy()
	tileRs[blueCoralSwitch] = jsh[5][4].Copy()
	tileRs[tealCoralSwitch] = jsh[6][4].Copy()
	tileRs[greenCoralSwitch] = jsh[7][4].Copy()

	tileRs[purpleCoralGate] = render.NewSwitch("closed", map[string]render.Modifiable{
		"closed": jsh[4][6].Copy(),
		"open":   jsh[4][5].Copy(),
	})
	tileRs[purpleCoralGateOff] = render.NewSwitch("open", map[string]render.Modifiable{
		"closed": jsh[4][6].Copy(),
		"open":   jsh[4][5].Copy(),
	})
	tileRs[blueCoralGate] = render.NewSwitch("closed", map[string]render.Modifiable{
		"closed": jsh[5][6].Copy(),
		"open":   jsh[5][5].Copy(),
	})
	tileRs[blueCoralGateOff] = render.NewSwitch("open", map[string]render.Modifiable{
		"closed": jsh[5][6].Copy(),
		"open":   jsh[5][5].Copy(),
	})
	tileRs[tealCoralGate] = render.NewSwitch("closed", map[string]render.Modifiable{
		"closed": jsh[6][6].Copy(),
		"open":   jsh[6][5].Copy(),
	})
	tileRs[tealCoralGateOff] = render.NewSwitch("open", map[string]render.Modifiable{
		"closed": jsh[6][6].Copy(),
		"open":   jsh[6][5].Copy(),
	})
	tileRs[greenCoralGate] = render.NewSwitch("closed", map[string]render.Modifiable{
		"closed": jsh[7][6].Copy(),
		"open":   jsh[7][5].Copy(),
	})
	tileRs[greenCoralGateOff] = render.NewSwitch("open", map[string]render.Modifiable{
		"closed": jsh[7][6].Copy(),
		"open":   jsh[7][5].Copy(),
	})
	tileRs[sandKey] = jsh[2][5].Copy()
	tileRs[coralExit] = jsh[3][5].Copy()
	tileRs[sandtrap] = render.NewSwitch("hole", map[string]render.Modifiable{
		"hole":   jsh[3][7].Copy(),
		"filled": jsh[3][6].Copy(),
	})
}
