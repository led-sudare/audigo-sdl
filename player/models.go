package player

import (
	"github.com/code560/audigo-sdl/util"
)

type Player interface {
	Play(args *PlayArgs)
	Stop()
	Volume(args *VolumeArgs)
}

type implPlayer interface {
	Player
}

type PlayArgs struct {
	Src  string `json:"src"`
	Loop bool   `json:"loop"`
}

type VolumeArgs struct {
	Vol float64 `json:"vol"`
}

const (
	dir = "asset/audio/"
)

type Proxy interface {
	GetChannel() chan<- *Action
}

type Action struct {
	Act  Actions
	Args interface{}
}

type Actions int

const (
	_ Actions = iota
	Play
	Stop
	Volume
)

type ctrler struct {
	Paused bool
}

var (
	log = util.GetLogger()
)
