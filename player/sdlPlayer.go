package player

import (
	"sync"

	"github.com/veandco/go-sdl2/mix"
)

type closing chan int

const (
	SDL_CH          = 2
	SDL_SAMPLE_RATE = 44100
	SDL_BUFFER      = 1024

	SDL_CHANNELS = 50

	SDL_INVALID_CH = -2
)

type sdlPlayer interface {
	Close()
	Play(file string, loop int)
	Stop()
	Volume(vol int)
}

func GetSdlPlayer() sdlPlayer {
	return &sdlController{
		channels: make(map[int]bool, SDL_CHANNELS),
		sdl:      getInstance(),
		statVol:  128,
	}
}

// implPlayer

type sdlController struct {
	sdl *sdlWrapper

	channels  map[int]bool
	chanMutex sync.Mutex

	// 0 - 128
	statVol int
}

func (i *sdlController) Close() {
	i.sdl.close()
}

func (i *sdlController) Play(file string, loop int) {
	close := make(closing)
	ch := i.sdl.play(file, loop, i.statVol, close)

	i.chanMutex.Lock()
	i.channels[ch] = true
	log.Debugf("play ch: %d", ch)
	i.chanMutex.Unlock()
	select {
	case <-close: // blocking
		i.chanMutex.Lock()
		delete(i.channels, ch)
		log.Debugf("close ch: %d", ch)
		i.chanMutex.Unlock()
		break
	}
}

func (i *sdlController) Stop() {
	i.chanMutex.Lock()
	for ch, enable := range i.channels {
		if enable {
			mix.HaltChannel(ch)
			log.Debugf("stop ch: %d", ch)
		}
	}
	i.chanMutex.Unlock()
}

func (i *sdlController) Volume(vol int) {
	if !validVolume(vol) {
		return
	}

	i.statVol = vol
	i.chanMutex.Lock()
	for ch, enable := range i.channels {
		if enable {
			i.sdl.volume(ch, vol)
			log.Debugf("set volume: %d -> %d", ch, vol)
		}
	}
	i.chanMutex.Unlock()
}

func validVolume(vol int) bool {
	if vol < 0 {
		return false
	} else if vol > 128 {
		return false
	}
	return true
}

// sdlWrapper

func getInstance() *sdlWrapper {
	return instanceSdl
}

var instanceSdl *sdlWrapper = newSdlWrapper()

type sdlWrapper struct {
	finishes    map[int]closing
	finishMutex sync.Mutex
}

func newSdlWrapper() *sdlWrapper {
	p := &sdlWrapper{}
	p.finishes = make(map[int]closing, SDL_CHANNELS)
	p.init()
	return p
}

func (p *sdlWrapper) init() {
	err := mix.Init(0)
	if err != nil {
		log.Fatal(err)
	}
	p.open()
}

func (p *sdlWrapper) open() {
	err := mix.OpenAudio(SDL_SAMPLE_RATE,
		mix.DEFAULT_FORMAT, SDL_CH, SDL_BUFFER)
	if err != nil {
		log.Fatal(err)
	}
	mix.AllocateChannels(SDL_CHANNELS)
	mix.ChannelFinished(func(ch int) {
		if ch < 0 {
			return
		}

		chunk := mix.GetChunk(ch)
		if chunk != nil {
			chunk.Free()
			p.finishMutex.Lock()
			close(p.finishes[ch])
			delete(p.finishes, ch)
			p.finishMutex.Unlock()
		}
	})
}

func (p *sdlWrapper) close() {
	mix.CloseAudio()
	mix.Quit()
}

func (p *sdlWrapper) play(file string, loop int, vol int, close closing) int {
	chunk := p.load(file)
	if chunk == nil {
		return SDL_INVALID_CH
	}
	chunk.Volume(vol)

	ch, err := chunk.Play(-1, loop)
	if err != nil {
		log.Error(err)
		return SDL_INVALID_CH
	}
	p.finishMutex.Lock()
	p.finishes[ch] = close
	p.finishMutex.Unlock()
	return ch
}

func (p *sdlWrapper) volume(ch int, vol int) {
	p.finishMutex.Lock()
	if _, ok := p.finishes[ch]; ok {
		mix.Volume(ch, vol)
	}
	p.finishMutex.Unlock()
}

func (p *sdlWrapper) load(file string) *mix.Chunk {
	chunk, err := mix.LoadWAV(file)
	if err != nil {
		log.Error(err)
		return nil
	}
	return chunk
}
