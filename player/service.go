package player

import (
	"sync"
)

var (
	playerPool sync.Pool
	proxyPool  sync.Pool
)

func init() {
	initPool()
}

func NewProxy() Proxy {
	proxy := proxyPool.Get().(*simpleProxy)
	proxy.playerPool = &playerPool
	return proxy
}

func CloseProxy(p Proxy) {
	proxyPool.Put(p)
}

func ClosePlayer(p Player) {
	playerPool.Put(p)
}

func initPool() {
	playerPool = sync.Pool{
		New: func() interface{} {
			return newSimplePlayer() // todo trans internal
		},
	}

	proxyPool = sync.Pool{
		New: func() interface{} {
			return newSimpleProxy() // todo trans internal
		},
	}
}
