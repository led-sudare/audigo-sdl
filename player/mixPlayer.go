package player

import (
	"math"
	"os"

	"github.com/code560/audigo-sdl/util"
)

const (
	SDL_INPUT_VOLUME_MIN  = 0.0
	SDL_INPUT_VOLUME_MAX  = 1.0
	SDL_OUTPUT_VOLUME_MIN = 0.0
	SDL_OUTPUT_VOLUME_MAX = 128.0
)

func newMixPlayer() Player {
	p := &mixPlayer{
		player: GetSdlPlayer(),
	}
	return p
}

type mixPlayer struct {
	player sdlPlayer
}

func (m *mixPlayer) Play(args *PlayArgs) {
	// open file
	if _, err := os.Stat(args.Src); err != nil {
		log.Warnf("not found music file: %s", args.Src)
		return
	}

	loop := 0
	if args.Loop {
		loop = -1
	}

	m.player.Play(args.Src, loop)
}

func (m *mixPlayer) Stop() {
	m.player.Stop()
}

func (m *mixPlayer) Volume(args *VolumeArgs) {
	volf := util.LinearTransF(float64(args.Vol),
		SDL_INPUT_VOLUME_MIN, SDL_INPUT_VOLUME_MAX,
		SDL_OUTPUT_VOLUME_MIN, SDL_OUTPUT_VOLUME_MAX)
	vol := int(math.Ceil(volf))
	log.Debugf("change vol unit: %f -> %d", args.Vol, vol)
	m.player.Volume(vol)
}
