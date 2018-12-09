package player

import (
	"sync"
	"testing"
	"time"

	"github.com/code560/audigo/util"
)

var (
	l = util.GetLogger()
)

func TestPlayer(t *testing.T) {
	dir := "../asset/audio/"
	args := []string{
		"bgm_wave.wav",
		"se_jump.wav",
		"info-girl1-bubu1.mp3",
	}

	wg := &sync.WaitGroup{}
	plist := make([]Player, len(args))
	for i, arg := range args {
		p := newSimplePlayer()
		plist[i] = p
		wg.Add(1)
		go func(p Player, name string) {
			p.Play(&PlayArgs{name, false, false})
			wg.Done()
		}(p, dir+arg)
	}

	wg.Wait()
	l.Debug("plaing sound")

	sec := time.Duration(2)
	time.Sleep(time.Second * sec)
	for _, p := range plist {
		log.Debug("call stop sound")
		p.Stop()
	}

	l.Debug("done routines")
}

func TestUnexpectedSound(t *testing.T) {
	loop := 10
	p := newSimplePlayer()
	for i := 0; i < loop; i++ {
		p.Stop()
	}
	for i := 0; i < loop; i++ {
		p.Pause()
	}
	for i := 0; i < loop; i++ {
		p.Resume()
	}
	arg := &VolumeArgs{1.8}
	for i := 0; i < loop; i++ {
		p.Volume(arg)
	}
	// for i := 0; i < loop; i++ {
	// 	p.Play(&PlayArgs{Src: "../asset/audio/bgm_wave.wav"})
	// 	fmt.Printf("%d ", i)
	// }
}
