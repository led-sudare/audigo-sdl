package player

import (
	"testing"
)

func TestCreate(t *testing.T) {
	p := newSimpleProxy()
	if p == nil {
		t.Error("dont created proxy: player.NewProxy()")
	}
}

func TestChan(t *testing.T) {
	p := newSimpleProxy()
	c := p.GetChannel()
	act := &Action{}
	act.Act = Play
	act.Args = &PlayArgs{Src: "bgm_wave.wav"}
	c <- act

	act.Act = Volume
	act.Args = &VolumeArgs{1}
	c <- act

	act.Act = Stop
	c <- act
}
