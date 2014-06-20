// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"fmt"
	"sync"
	"time"

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

	onlineCache  map[string]*store.Msg
	offlineCache []*user_msg

	CloseChan   chan byte // Other gorountine Call notice exit
	isSendClose bool

	isStop bool
	lock   *sync.Mutex

	isLetClose bool

	isOfflineClose bool // Is offline chan has been close
}

func newClient(rw *spp.Conn, id string) *client {
	return &client{
		queue:     NewPackQueue(rw),
		id:        id,
		CloseChan: make(chan byte),
		lock:      new(sync.Mutex),

		offlines: make(chan *store.Msg, Conf.MaxCacheMsg),
		onlines:  make(chan *store.Msg, Conf.MaxCacheMsg),

		onlineCache: make(map[string]*store.Msg),
		// offlineCache: make([]*user_msg, 1),
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

		findEndChan <-chan byte
	)

	// Start the write queue
	go c.queue.writeLoop()

	c.offlines, findEndChan = store.GetOfflineMsg(c.id, noticeFin)
	c.readChan = c.queue.ReadPackInLoop(noticeFin)

	// Start push
loop:
	for {
		select {
		case <-findEndChan:
			glog.Info("Offline has been finish send.")
			c.isOfflineClose = true
			if err = c.pushOfflineMsg(nil); err != nil {
				break loop
			}
			break
		case msg = <-c.offlines:
			// Push the offline msg
			glog.Info("Get a offline msg")
			err = c.pushOfflineMsg(msg)
			if err != nil {
				break loop
			}

		case msg = <-c.onlines:
			// Push the online msg
			glog.Info("Get a online msg")
			// Check the msg time
			if time.Now().UTC().Unix() > msg.Expired {
				// cancel send the msg
				break
			}
			// Add the msg into cache
			if msg == nil {
				glog.Warning("onlines has been close")
				break loop
			}
			if len(c.onlineCache) > Conf.MaxCacheMsg && Conf.MaxCacheMsg != 0 {
				err = fmt.Errorf("Online msg is out of range:%v", len(c.onlineCache))
				break loop
			}
			msg.Typ = ONLINE
			c.onlineCache[msg.Msg_id] = msg
			err = c.pushMsg(msg)
			if err != nil {
				break loop
			}
		case pAndErr = <-c.readChan:
			// If connetion has a error, should break
			// if it return a timeout error, illustrate
			// hava not recive a heart beat pack at an
			// given time.
			if pAndErr.err != nil {
				glog.Info("Get a connection error , will break")
				err = pAndErr.err
				break loop
			}
			glog.Infof("Client msg(%v)\n", pAndErr.pack.Typ)

			// Choose the requst type
			switch pAndErr.pack.Typ {
			case SINGEL:
				sg := new(singleMsg)
				if err = unMarshalJson(pAndErr.pack.Body, sg); err != nil {
					glog.Errorf("Marshal a json error:%v", err)
					break
				}
				c.delMsg(sg.Typ, sg.Id)
			case MUTIL:
				mutil := new(mutilMsg)
				if err = unMarshalJson(pAndErr.pack.Body, mutil); err != nil {
					glog.Errorf("Marshal a json error:%v", err)
					break
				}
				for _, sg := range mutil.Msgs {
					c.delMsg(sg.Typ, sg.Id)
				}

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
	for i := 0; i < 2; i++ {
		noticeFin <- 1
	}
	// Wrte the onlines msg to the db
	for _, v := range c.onlineCache {
		v.Typ = OFFLINE
		store.InsertOfflineMsg(v)
	}

	// Close the online msg channel
	c.lock.Lock()
	c.isStop = true
	close(c.onlines)
	c.lock.Unlock()
	if c.isSendClose {
		c.CloseChan <- 0
	}
	close(c.CloseChan)

	Users.Del(c.id)

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

func (c *client) delMsg(typ int, msg_id string) {
	switch typ {
	case OFFLINE:
		// Del the offline msg in the store
		glog.Info("Del a offline msg")
		store.DelOfflineMsg(msg_id, c.id)
	case ONLINE:
		// Del the online msg in the msg cache
		glog.Info("Del a online msg")
		delete(c.onlineCache, c.id)
	default:
		glog.Errorf("No define the type:%v", typ)
	}
}

// func pushMsgs(msg []*store.Msg) (err error) {
// }
func (c *client) pushOfflineMsg(msg *store.Msg) (err error) {
	// The max cache size is 20
	if c.isOfflineClose && msg == nil {
		// Send
		err = c.sendOfflineMsg(&offineMsg{
			Ms: c.offlineCache,
		})
	} else {
		// Wait
		if len(c.offlineCache) > 20 {
			// Send
			err = c.sendOfflineMsg(&offineMsg{
				Ms: c.offlineCache,
			})
			// Clean the cache
			c.offlineCache = []*user_msg{getUserMsg(msg)}
		} else {
			c.offlineCache = append(c.offlineCache, getUserMsg(msg))
		}
	}
	return
}
func (c *client) sendOfflineMsg(ms *offineMsg) (err error) {
	var buf []byte
	buf, err = getOffineMsg(ms)
	if err != nil {
		return
	}
	// Set a pack
	pack, err := c.setPack(PUSH_OFFLINE, buf)
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

func WriteOnlineMsg(msg *store.Msg) {
	// fix the Expired
	if msg.Expired > 0 {
		msg.Expired += time.Now().UTC().Unix()
	}

	c := Users.Get(msg.To_id)
	if c == nil {
		msg.Typ = OFFLINE
		store.InsertOfflineMsg(msg)
		return
	}

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
