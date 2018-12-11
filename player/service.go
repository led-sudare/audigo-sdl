package player

import (
	"sync"
)

var (
	proxyPool sync.Pool
)

func init() {
	initPool()
}

func NewProxy() Proxy {
	proxy := proxyPool.Get().(*simpleProxy)
	return proxy
}

func CloseProxy(p Proxy) {
	proxyPool.Put(p)
}

func initPool() {
	proxyPool = sync.Pool{
		New: func() interface{} {
			return newSimpleProxy() // todo trans internal
		},
	}
}
