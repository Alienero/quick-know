// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"sync"
)

var (
	Users = &clientUs{
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

func (c *clientUs) Get(id string) *client {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.m[id]
}
func (c *clientUs) Set(id string, v *client) {
	c.l.Lock()
	c.m[id] = v
	c.l.Unlock()
}
func (c *clientUs) Del(id string) {
	c.l.Lock()
	delete(c.m, id)
	c.l.Unlock()
}

type ctrlUs struct {
	m map[string]*ControlServer
	l *sync.RWMutex
}

func (c *ctrlUs) Get(id string) *ControlServer {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.m[id]
}
func (c *ctrlUs) Set(id string, v *ControlServer) {
	c.l.Lock()
	c.m[id] = v
	c.l.Unlock()
}
func (c *ctrlUs) Del(id string) {
	c.l.Lock()
	delete(c.m, id)
	c.l.Unlock()
}
