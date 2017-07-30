package game

import "github.com/oakmound/oak/collision"

// The only collision label used is "blocking", for blocking terrain.
// Tiles are also cast to collision labels, but they don't need to
// be re-iterated here. Tiles start at 10, so they don't overlap with 
// blocking.
const (
	blocking collision.Label = iota
)
