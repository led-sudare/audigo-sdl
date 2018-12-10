package net

import (
	"testing"
	"time"

	"github.com/code560/audigo-sdl/player"
)

func TestRequest(t *testing.T) {
	req := NewRequester("http://localhost:8080", "abc")
	req.GetQueue() <- &Action{
		Act: "play",
		Args: &player.PlayArgs{
			Src: "bgm_wave.wav",
		},
	}

	time.Sleep(time.Second * 10)
}
