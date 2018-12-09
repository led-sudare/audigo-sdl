package app

import (
	"github.com/code560/audigo/net"
	"github.com/code560/audigo/player"
)

type voice int

const (
	_ voice = iota
	start
	finish
)

var sePlayer player.Player

func Serve(port string) {
	r := net.NewRouter()
	se(start)
	r.Run(port)
	se(finish)
}

func se(v voice) {
	sePlayer = player.NewInternalPlayer()
	var sound string
	switch v {
	case start:
		sound = "boyomi-chan_start.wav"
	case finish:
		sound = "rusuden_04-2.wav"
	default:
		return
	}
	sePlayer.Play(&player.PlayArgs{Src: sound}) // blocking
}
