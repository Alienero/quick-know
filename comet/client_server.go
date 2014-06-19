// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"

	"github.com/golang/glog"
)

var notAlive = errors.New("Connection was dead")

type client struct {
	queue *PackQueue
	id    string // Owner+id

	offlines <-chan *store.Msg
	onlines  chan *store.Msg
	readChan <-chan *packAndErr

	onlineCache map[string]*store.Msg

	CloseChan   chan byte // Other gorountine Call notice exit
	isSendClose bool

	isStop bool
	lock   *sync.Mutex

	isLetClose bool
}

func newClient(rw *spp.Conn, id string) *client {
	return &client{
		queue:     NewPackQueue(rw),
		id:        id,
		CloseChan: make(chan byte),
		lock:      new(sync.Mutex),

		offlines: make(chan *store.Msg, Conf.MaxCacheMsg),
		onlines:  make(chan *store.Msg, Conf.MaxCacheMsg),
	}
}

// Push the msg and response the heart beat
func (c *client) listen_loop() (e error) {

	var (
		err     error
		msg     *store.Msg
		pAndErr *packAndErr
		pack    *spp.Pack

		noticeFin = make(chan byte)

		isOfflineClose bool // Is offline chan has been close
	)

	// Start the write queue
	go c.queue.writeLoop()

	c.offlines = store.GetOfflineMsg(c.id, noticeFin)
	c.readChan = c.queue.ReadPackInLoop(noticeFin)

	// Start push
loop:
	for {
		select {

		case msg = <-c.offlines:
			// Push the offline msg
			glog.Info("One offline msg")
			if msg == nil {
				isOfflineClose = true
				// Close the offline chan
				noticeFin <- 1
				glog.Errorln("offlines has been close")
				break
			}
			// msg.Typ = OFFLINE
			err = c.pushMsg(msg)
			if err != nil {
				break loop
			}
		case msg = <-c.onlines:
			// Push the online msg
			glog.Info("One online msg")
			// Add the msg into cache
			if msg == nil {
				glog.Errorln("onlines has been close")
				break
			}
			if len(c.onlineCache) > Conf.MaxCacheMsg && Conf.MaxCacheMsg != 0 {
				err = fmt.Errorf("Online msg is out of range:%v", len(c.onlineCache))
				break loop
			}
			msg.Typ = ONLINE
			c.onlineCache[msg.Msg_id] = msg
			err = c.pushMsg(msg)
			if err != nil {
				break
			}
		case pAndErr = <-c.readChan:
			// If connetion has a error, should break
			// if it return a timeout error, illustrate
			// hava not recive a heart beat pack at an
			// given time.
			// glog.Infof("One client msg(%v)\n",pAndErr.pack)
			if pAndErr.err != nil {
				glog.Info("Get error will break")
				err = pAndErr.err
				break loop
			}
			glog.Infof("One client msg(%v)\n", pAndErr.pack.Typ)

			// Choose the requst type
			switch pAndErr.pack.Typ {
			case OFFLINE:
				// Del the offline msg in the store
				glog.Info("Del a offline msg")
				store.DelOfflineMsg(string(pAndErr.pack.Body), c.id)
			case ONLINE:
				// Del the online msg in the msg cache
				glog.Info("Del a online msg")
				delete(c.onlineCache, string(pAndErr.pack.Body))
			case HEART_BEAT:
				// Reply the heart beat
				glog.Info("hb msg")
				pack, err = c.setPack(HEART_BEAT, []byte("OK"))
				if err != nil {
					break loop
				}
				err = c.queue.WritePack(pack)
				if err != nil {
					break loop
				}
			default:
				// Not define pack type
				glog.Errorf("The type not define:%v\n", pAndErr.pack.Typ)
			}
		case <-c.CloseChan:
			// Start close
			glog.Info("Will break new relogin")
			c.isSendClose = true
			break loop

		}
	}

	// Wrte the onlines msg to the db
	// Free resources
	// Close channels
	glog.Info("Has been break")
	if isOfflineClose {
		glog.Info("One fin")
		noticeFin <- 1
	} else {
		glog.Info("Two fin")
		for i := 0; i < 2; i++ {
			noticeFin <- 1
		}
	}
	glog.Info("Check Lock")

	// Wrte the onlines msg to the db
	for _, v := range c.onlineCache {
		v.Typ = OFFLINE
		store.InsertOfflineMsg(v)
	}

	// Close the online msg channel
	Users.Del(c.id)
	c.lock.Lock()
	c.isStop = true
	close(c.onlines)
	c.lock.Unlock()
	if c.isSendClose {
		c.CloseChan <- 0
	}
	close(c.CloseChan)

	return
}
func (c *client) pushMsg(msg *store.Msg) (err error) {
	var buf []byte
	buf, err = getMsg(msg)
	if err != nil {
		return
	}
	// Set a pack
	pack, err := c.setPack(PUSH_MSG, buf)
	if err != nil {
		return
	}
	// Write this pack
	err = c.queue.WritePack(pack)
	return
}
func (c *client) setPack(typ int, body []byte) (*spp.Pack, error) {
	return c.queue.rw.SetDefaultPack(typ, body)
}
func (c *client) LetClose() bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.isLetClose {
		c.isLetClose = true
		return true
	} else {
		return false
	}
}

func WriteOnlineMsg(msg *store.Msg) {
	c := Users.Get(msg.To_id)
	if c == nil {
		msg.Typ = OFFLINE
		store.InsertOfflineMsg(msg)
		return
	}
	// defer c.lock.Unlock()
	c.lock.Lock()
	if len(c.onlines) == Conf.MaxCacheMsg {
		msg.Typ = OFFLINE
		store.InsertOfflineMsg(msg)
		c.lock.Unlock()
		return
	}
	if c.isStop {
		c.lock.Unlock()
		msg.Typ = OFFLINE
		store.InsertOfflineMsg(msg)
	} else {
		msg.Typ = ONLINE
		c.onlines <- msg
		c.lock.Unlock()
	}
}
