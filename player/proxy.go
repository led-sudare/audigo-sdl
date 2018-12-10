package player

import (
	"sync"

	"github.com/code560/audigo-sdl/util"
)

const (
	ChanSize = 50
)

type simpleProxy struct {
	// playerPool *sync.Pool
	act     chan *Action
	closing chan struct{}

	// plays   map[int]Player
	player  Player
	playMtx sync.Mutex

	// playerCtrl *ctrler
	// playerVol  *effects.Volume
}

// NewProxy は、Playerを生成して返します。
func newSimpleProxy() Proxy {
	p := &simpleProxy{
		act:     make(chan *Action, ChanSize),
		closing: make(chan struct{}),

		// plays: make(map[int]Player, 32),
		player: newMixPlayer(),

		// playerCtrl: makeCtrl(),
		// playerVol:  makeVolume(),
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
			// player := p.playerPool.Get().(Player)
			// p.setFactory(player)
			// var i int
			// p.playLock(func() {
			// 	i = p.pushPlayer(player)
			// })
			// player.Play(a)
			// p.playLock(func() {
			// 	p.popPlayer(i)
			// })
			// p.playerPool.Put(player)
		}(p, a)
	case Stop:
		go p.player.Stop()
		// p.playLock(func() {
		// 	for _, player := range p.plays {
		// 		go player.Stop()
		// 	}
		// })
	case Pause:
		// p.playerCtrl.Paused = true
		go p.player.Pause()
		// p.playLock(func() {
		// 	for _, player := range p.plays {
		// 		go player.Pause()
		// 	}
		// })
	case Resume:
		go p.player.Resume()
		// p.playerCtrl.Paused = false
		// p.playLock(func() {
		// 	for _, player := range p.plays {
		// 		go player.Resume()
		// 	}
		// })
	case Volume:
		a := arg.Args.(*VolumeArgs)
		go p.player.Volume(a)
		// p.volume(a)
		// p.playLock(func() {
		// 	for _, player := range p.plays {
		// 		go player.Volume(a)
		// 	}
		// })
	default:
		log.Warn("nothing call player function")
	}
}

// func (p *simpleProxy) playLock(f func()) {
// 	p.playMtx.Lock()
// 	f()
// 	p.playMtx.Unlock()
// }

// func (p *simpleProxy) atPlayer(i int) Player {
// 	v, ok := p.plays[i]
// 	if ok {
// 		return v
// 	} else {
// 		return nil
// 	}
// }

// func (p *simpleProxy) pushPlayer(player Player) int {
// 	i := len(p.plays)
// 	p.plays[i] = player
// 	return i
// }

// func (p *simpleProxy) popPlayer(i int) {
// 	_, ok := p.plays[i]
// 	if ok {
// 		delete(p.plays, i)
// 	}
// }

func isDone(c chan struct{}) bool {
	return util.IsDone(c)
}

// func (p *simpleProxy) volume(args *VolumeArgs) {
// 	if args.Vol == 0 {
// 		p.playerVol.Silent = true
// 	} else {
// 		p.playerVol.Silent = false
// 	}
// 	p.playerVol.Volume = args.Vol
// }

// func (p *simpleProxy) setFactory(player interface{}) {
// 	impl, ok := player.(implPlayer)
// 	if !ok {
// 		return
// 	}

// 	ctrl := p.playerCtrl
// 	impl.setCtrlFactory(func() *ctrler {
// 		c := makeCtrl()
// 		c.Paused = ctrl.Paused
// 		return c
// 	})

// 	vol := p.playerVol
// 	impl.setVolumeFactory(func() *effects.Volume {
// 		v := makeVolume()
// 		v.Base = vol.Base
// 		v.Silent = vol.Silent
// 		v.Volume = vol.Volume
// 		return v
// 	})
// }
