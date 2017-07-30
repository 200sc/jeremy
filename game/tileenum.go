package game

const (
	_ tile = iota
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
	// Subtract 7 from the line number for the csv number
	sand
	coral
	purpleCoralGate
	purpleCoralSwitch
	blueCoralGate
	blueCoralSwitch
	tealCoralGate
	tealCoralSwitch
	sandtrap
	sandglob
	sandgeyser
	jeremyTile
	verticalCrab
	treasure
	sandKey
	coralExit
	greenCoralGate
	greenCoralSwitch
	greenCoralGateOff
	purpleCoralGateOff
	blueCoralGateOff
	tealCoralGateOff
	horizontalCrab
)

type tile int
