package util

import (
	"fmt"
	"sync"
)

func RecoveGo(f func() error) error {
	e := make(chan error, 1)
	go func() {
		if err := Recove(f); err != nil {
			e <- err
		}
		e <- nil
	}()
	return <-e
}

func Recove(f func() error) error {
	log := GetLogger()
	defer func() error {
		if r := recover(); r != nil {
			log.Error("recorvered from ", r)
			return fmt.Errorf("error %s", r)
		}
		return nil
	}()
	return f()
}

var muClosing sync.Mutex

type Closing struct {
	done chan struct{}
}

func NewClosing() *Closing {
	muClosing.Lock()
	defer muClosing.Unlock()
	c := &Closing{
		done: make(chan struct{}),
	}
	return c
}

func (c *Closing) Reset() {
	c.Close()

	muClosing.Lock()
	defer muClosing.Unlock()
	c.done = make(chan struct{})
}

func (c *Closing) Close() {
	muClosing.Lock()
	defer muClosing.Unlock()
	if !IsDone(c.done) {
		close(c.done)
	}
}

func (c *Closing) GetDone() <-chan struct{} {
	return c.done
}

func (c *Closing) GetDo() chan<- struct{} {
	return c.done
}

func (c *Closing) IsDone() bool {
	return IsDone(c.done)
}

func IsDone(c <-chan struct{}) bool {
	select {
	case <-c:
		return true
	default:
		return false
	}
}
