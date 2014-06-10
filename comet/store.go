package comet

import (
	"sync"
)

var (
	uesers = &clientUs{
		m: make(map[string]*client),
		l: new(sync.RWMutex),
	}

	ctrls = &ctrlUs{
		m: make(map[string]*ControlServer),
		l: new(sync.RWMutex),
	}
)

type clientUs struct {
	m map[string]*client
	l *sync.RWMutex
}

func (c *clientUs) get(id string) *client {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.m[id]
}
func (c *clientUs) set(id string, v *client) {
	c.l.Lock()
	c.m[id] = v
	c.l.Unlock()
}
func (c *clientUs) del(id string) {
	c.l.Lock()
	delete(c.m, id)
	c.l.Unlock()
}

type ctrlUs struct {
	m map[string]*ControlServer
	l *sync.RWMutex
}

func (c *ctrlUs) get(id string) *ControlServer {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.m[id]
}
func (c *ctrlUs) set(id string, v *ControlServer) {
	c.l.Lock()
	c.m[id] = v
	c.l.Unlock()
}
func (c *ctrlUs) del(id string) {
	c.l.Lock()
	delete(c.m, id)
	c.l.Unlock()
}
