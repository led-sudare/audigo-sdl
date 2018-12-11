package player

import (
	"sync"

	"github.com/code560/audigo-sdl/util"
)

const (
	ChanSize = 50
)

type simpleProxy struct {
	act     chan *Action
	closing chan struct{}

	player  Player
	playMtx sync.Mutex
}

// NewProxy は、Playerを生成して返します。
func newSimpleProxy() Proxy {
	p := &simpleProxy{
		act:     make(chan *Action, ChanSize),
		closing: make(chan struct{}),
		player:  newMixPlayer(),
	}
	go p.work()
	return p
}

func (p *simpleProxy) GetChannel() chan<- *Action {
	return p.act
}

func (p *simpleProxy) work() {
	for {
		select {
		case v := <-p.act:
			if isDone(p.closing) {
				return
			}
			p.call(v)
		}
	}
}

func (p *simpleProxy) call(arg *Action) {
	switch arg.Act {
	case Play:
		a := arg.Args.(*PlayArgs)
		a.Src = dir + a.Src
		go func(p *simpleProxy, a *PlayArgs) {
			p.player.Play(a)
		}(p, a)
	case Stop:
		go p.player.Stop()
	case Volume:
		a := arg.Args.(*VolumeArgs)
		go p.player.Volume(a)
	default:
		log.Warn("nothing call player function")
	}
}

func isDone(c chan struct{}) bool {
	return util.IsDone(c)
}
