package player

import (
	"os"

	"github.com/faiface/beep"
)

type simplePlayer struct {
	playerMaker
}

func newSimplePlayer() Player {
	p := &simplePlayer{}
	return p
}

func (p *simplePlayer) Play(args *PlayArgs) {
	// init
	p.reset()
	// and stop
	if args.Stop {
		p.Pause()
	}
	// open file
	if _, err := os.Stat(args.Src); err != nil {
		log.Warnf("not found music file: %s", args.Src)
		return
	}
	closer, format := p.openFile(args.Src)
	defer closer.Close()
	// set middlewares
	s := beep.Loop(loopCount(args.Loop), closer)
	s = p.setCtrlStream(s)
	s = p.setVolumeStream(s)
	p.mixer = p.makeMixer()
	// if err := p.makeOtoPlayer(format.SampleRate, format.SampleRate.N(time.Millisecond*CHUNK)); err != nil {
	if err := p.makeOtoPlayer(format.SampleRate, CHUNK_SIZE); err != nil {
		log.Warnf("dont create oto player: %s", err.Error())
		return
	}
	defer p.oto.Close()

	// play sound
	p.mixer.Play(s)
	p.sampling(closer) // blocking
	p.Stop()
}

func (p *simplePlayer) Volume(args *VolumeArgs) {
	if p.vol == nil {
		p.vol = p.makeVolume()
	}
	p.volume(args.Vol)
}

func (p *simplePlayer) Pause() {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	p.ctrl.Paused = true
}

func (p *simplePlayer) Resume() {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	p.ctrl.Paused = false
}

func (p *simplePlayer) reset() {
	p.close = false

	p.storeMutex.Lock()
	p.ctrl = nil
	p.vol = nil
	p.mixer = nil
	p.oto = nil
	p.samples = nil
	p.buf = nil
	p.storeMutex.Unlock()
}
