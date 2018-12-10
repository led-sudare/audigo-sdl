package player

type internalPlayer struct {
	mixPlayer
}

func NewInternalPlayer() Player {
	p := &internalPlayer{}
	return p
}

func (p *internalPlayer) Play(args *PlayArgs) {
	args.Src = dir + args.Src
	log.Debugf("play file: %s", args.Src)
	p.mixPlayer.Play(args)
}
