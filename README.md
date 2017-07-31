Jeremy the Clam is a game made for the [Gopher Game jam](https://itch.io/jam/gopher-jam), in the last three days of the game jam instead of in the preceding three months.
Jeremy is a simple tile-based puzzle game, but I think it uses some nice patterns that made coding it really enjoyable.
I recommend playing with the window maximized.

## Installing

Run `go get -u github.com/200sc/jeremy/...`

Then `go run core.go` in that folder to run the program.

### Linux

You'll also need ALSA: `sudo apt-get install alsa-base libasound2-dev`

## Controls

E to pick up sand from geysers or hills, or keys.

Space to drop.

Hold R to reset a level if you're stuck.

## Modding?

If you want to try to make your own levels, just look at the csvs in the level folder-- the numbering refers to indices in `tileenum.go`.

## Support

Tested primarily on Windows, a little on Linux, should work on OSX (but no SFX)
