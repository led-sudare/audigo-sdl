package player

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/oto"

	"github.com/code560/audigo/util"
)

// statics

var (
	mtxOto sync.Mutex
	log    = util.GetLogger()
)

// consts

const (
	CH         = 2
	BPS        = 2
	CHUNK      = 100
	CHUNK_SIZE = 2048

	VOLUME_BASE = 1.2
	VOLUME_INIT = 1
)

// functions

func loopCount(enable bool) int {
	if enable {
		return -1
	} else {
		return 1
	}
}

// create

type playerMaker struct {
	close         bool
	storeMutex    sync.Mutex
	ctrlFactory   func() *ctrler
	volumeFactory func() *effects.Volume

	ctrl    *ctrler
	vol     *effects.Volume
	mixer   *beep.Mixer
	oto     *oto.Player
	samples [][CH]float64
	buf     []byte
}

// interface methods

func (p *playerMaker) Stop() {
	p.close = true
}

func (p *playerMaker) Volume(args *VolumeArgs) {
	p.volume(args.Vol)
}

func (p *playerMaker) volume(vol float64) {
	if vol == 0 {
		p.vol.Silent = true
	} else {
		p.vol.Silent = false
	}
	p.vol.Volume = vol
}

func (p *playerMaker) Pause() {
	p.ctrl.Paused = true
}

func (p *playerMaker) Resume() {
	p.ctrl.Paused = false
}

// maker methods

func (p *playerMaker) makeMixer() *beep.Mixer {
	return &beep.Mixer{}
}

func (p *playerMaker) makeCtrl() *ctrler {
	if p.ctrlFactory != nil {
		return p.ctrlFactory()
	} else {
		return makeCtrl()
	}
}

func makeCtrl() *ctrler {
	return &ctrler{}
}

func (p *playerMaker) makeVolume() *effects.Volume {
	if p.volumeFactory != nil {
		return p.volumeFactory()
	} else {
		return makeVolume()
	}
}

func makeVolume() *effects.Volume {
	return &effects.Volume{Base: VOLUME_BASE, Volume: VOLUME_INIT}
}

func (p *playerMaker) makeClosing() *util.Closing {
	return util.NewClosing()
}

func (p *playerMaker) makeOtoPlayer(sampleRate beep.SampleRate, bufferSize int) error {
	var err error
	bufferNum := bufferSize * CH * BPS
	mtxOto.Lock()
	p.oto, err = oto.NewPlayer(int(sampleRate), CH, BPS, bufferNum)
	mtxOto.Unlock()
	if err != nil {
		return err
	}
	p.samples = make([][CH]float64, bufferSize)
	p.buf = make([]byte, bufferNum)
	return nil
}

func (p *playerMaker) sampling(s beep.StreamSeeker) {
	for pos, len := s.Position(), s.Len(); pos < len; {
		// read stream
		p.storeMutex.Lock()
		if p.samples == nil {
			p.storeMutex.Unlock()
			return
		}
		for p.ctrl != nil && p.ctrl.Paused {
			time.Sleep(time.Millisecond * 100)
			if p.close {
				return
			}
		}
		n, _ := p.mixer.Stream(p.samples)
		p.storeMutex.Unlock()
		pos += n

		if p.close {
			return
		}

		// write buffer
		var i16 int16
		var l, h byte
		var pos int
		var val float64
		for s := range p.samples {
			for rl := range p.samples[s] {
				val = p.samples[s][rl]
				if val < -1 {
					val = -1
				}
				if val > +1 {
					val = +1
				}
				i16 = int16(val * math.MaxInt16)
				l = byte(i16)
				h = byte(i16 >> 8)
				pos = s*4 + rl*2
				p.buf[pos] = l
				p.buf[pos+1] = h
			}
		}
		p.oto.Write(p.buf)
	}
}

func (p *playerMaker) setCtrlStream(s beep.Streamer) beep.Streamer {
	if p.ctrl == nil {
		p.ctrl = p.makeCtrl()
	}
	return s
}

func (p *playerMaker) setVolumeStream(s beep.Streamer) beep.Streamer {
	if p.vol == nil {
		p.vol = p.makeVolume()
	}
	p.vol.Streamer = s
	return p.vol
}

func (p *playerMaker) setCtrlFactory(f func() *ctrler) {
	p.ctrlFactory = f
}

func (p *playerMaker) setVolumeFactory(f func() *effects.Volume) {
	p.volumeFactory = f
}

func (p *playerMaker) openFile(src string) (beep.StreamSeekCloser, *beep.Format) {

	// open file
	f, err := os.Open(src)
	if err != nil {
		log.Errorf("dont open file: %s\n%s", src, err.Error())
		return nil, nil
	}
	// decode file
	ext := strings.ToLower(filepath.Ext(src))
	var closer beep.StreamSeekCloser
	var format beep.Format
	switch ext {
	case ".wav":
		closer, format, err = wav.Decode(f)
	case ".mp3":
		closer, format, err = mp3.Decode(f)
	default:
		log.Errorf("dont support file: %s", src)
		return nil, nil
	}
	if err != nil {
		log.Error("dont decode file: %s\n%s", src, err.Error())
		return nil, nil
	}

	return closer, &format
}
