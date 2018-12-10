package player

import (
	"log"
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

type Player interface {
	Close()
	Play(file string)
	Stop()
	Volume(vol int)
}

func GetPlayer(contentId_ string) Player {
	return &implPlayer{
		contentId: contentId_,
		channels:  make(map[int]bool, SDL_CHANNELS),
		sdl:       getInstance(),
		statVol:   128,
	}
}

// implPlayer

type implPlayer struct {
	contentId string
	sdl       *sdlPlayer

	channels  map[int]bool
	chanMutex sync.Mutex

	// 0 - 128
	statVol int
}

func (i *implPlayer) Close() {
	i.sdl.close()
}

func (i *implPlayer) Play(file string) {
	close := make(closing)
	ch := i.sdl.play(file, i.statVol, close)

	i.chanMutex.Lock()
	i.channels[ch] = true
	log.Printf("[%s] play ch: %d", i.contentId, ch)
	i.chanMutex.Unlock()
	select {
	case <-close: // blocking
		i.chanMutex.Lock()
		delete(i.channels, ch)
		log.Printf("[%s] close ch: %d", i.contentId, ch)
		i.chanMutex.Unlock()
		break
	}
}

func (i *implPlayer) Stop() {
	i.chanMutex.Lock()
	for ch, enable := range i.channels {
		if enable {
			mix.HaltChannel(ch)
			log.Printf("[%s] stop ch: %d", i.contentId, ch)
		}
	}
	i.chanMutex.Unlock()
}

func (i *implPlayer) Volume(vol int) {
	if !validVolume(vol) {
		return
	}

	i.statVol = vol
	i.chanMutex.Lock()
	for ch, enable := range i.channels {
		if enable {
			i.sdl.volume(ch, vol)
			log.Printf("[%s] set volume: %d -> %d", i.contentId, ch, vol)
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

// sdlPlayer

func getInstance() *sdlPlayer {
	return instance_player
}

var instance_player *sdlPlayer = newSdlPlayer()

type sdlPlayer struct {
	finishes    map[int]closing
	finishMutex sync.Mutex
}

func newSdlPlayer() *sdlPlayer {
	p := &sdlPlayer{}
	p.finishes = make(map[int]closing, SDL_CHANNELS)
	p.init()
	return p
}

func (p *sdlPlayer) init() {
	err := mix.Init(0)
	if err != nil {
		log.Fatal(err)
	}
	p.open()
}

func (p *sdlPlayer) open() {
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

func (p *sdlPlayer) close() {
	mix.CloseAudio()
	mix.Quit()
}

func (p *sdlPlayer) play(file string, vol int, close closing) int {
	chunk := p.load(file)
	if chunk == nil {
		return SDL_INVALID_CH
	}
	chunk.Volume(vol)

	ch, err := chunk.Play(-1, 0)
	if err != nil {
		log.Print(err)
		return SDL_INVALID_CH
	}
	p.finishMutex.Lock()
	p.finishes[ch] = close
	p.finishMutex.Unlock()
	return ch
}

func (p *sdlPlayer) volume(ch int, vol int) {
	p.finishMutex.Lock()
	if _, ok := p.finishes[ch]; ok {
		mix.Volume(ch, vol)
	}
	p.finishMutex.Unlock()
}

func (p *sdlPlayer) load(file string) *mix.Chunk {
	chunk, err := mix.LoadWAV(file)
	if err != nil {
		log.Print(err)
		return nil
	}
	return chunk
}
