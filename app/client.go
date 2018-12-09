package app

import (
	"github.com/code560/audigo/net"
	"github.com/code560/audigo/player"
)

type Client interface {
	Play(args *player.PlayArgs)
	Volume(args *player.VolumeArgs)
	Stop()
	Pause()
	Resume()
}

type simpleClient struct {
	r  *net.Requester
	id string
}

func NewClient(domain string, id string) Client {
	c := &simpleClient{
		id: id,
	}
	c.init(domain)
	return c
}

func (c *simpleClient) init(domain string) {
	c.r = net.NewRequester(domain, c.id)
}

func (c *simpleClient) Play(args *player.PlayArgs) {
	c.r.GetQueue() <- &net.Action{
		Act:  "play",
		Args: args,
	}
}

func (c *simpleClient) Volume(args *player.VolumeArgs) {
	c.r.GetQueue() <- &net.Action{
		Act:  "volume",
		Args: args,
	}
}

func (c *simpleClient) Stop() {
	c.r.GetQueue() <- &net.Action{
		Act:  "stop",
		Args: nil,
	}
}

func (c *simpleClient) Pause() {
	c.r.GetQueue() <- &net.Action{
		Act:  "pause",
		Args: nil,
	}
}

func (c *simpleClient) Resume() {
	c.r.GetQueue() <- &net.Action{
		Act:  "resume",
		Args: nil,
	}
}

func (c *simpleClient) Ping() {
	c.r.GetQueue() <- &net.Action{
		Act:  "ping",
		Args: nil,
	}
}
