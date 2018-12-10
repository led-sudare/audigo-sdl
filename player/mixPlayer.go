package player

import "os"

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
	vol := (int)(args.Vol * 10)
	m.player.Volume(vol)
}

func (m *mixPlayer) Pause() {
	// nothing
	log.Warn("dont support.")
}

func (m *mixPlayer) Resume() {
	// nothing
	log.Warn("dont support.")
}
