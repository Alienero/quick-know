// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/Alienero/quick-know/store"
	// "github.com/Alienero/spp"
	"github.com/Alienero/quick-know/mqtt"

	"github.com/golang/glog"
)

var notAlive = errors.New("Connection was dead")

type client struct {
	queue *PackQueue
	id    string

	offlines <-chan *store.Msg
	onlines  chan *store.Msg
	readChan <-chan *packAndErr

	onlineCache map[int]*store.Msg
	offlineNum  int

	CloseChan   chan byte // Other gorountine Call notice exit
	isSendClose bool

	isStop bool
	lock   *sync.Mutex

	isLetClose bool

	isOfflineClose bool // Is offlines chan has been close

	// Online msg id
	curr_id int
	m       int // Residue msgs' ids

	// Offline msg id
	curr_offline int
	n            int
}

func newClient(r *bufio.Reader, w *bufio.Writer, conn net.Conn, id string) *client {
	return &client{
		queue:     NewPackQueue(r, w, conn),
		id:        id,
		CloseChan: make(chan byte),
		lock:      new(sync.Mutex),

		offlines: make(chan *store.Msg, Conf.MaxCacheMsg),
		onlines:  make(chan *store.Msg, Conf.MaxCacheMsg),

		onlineCache: make(map[int]*store.Msg),

		curr_id: -1,
		m:       65536,

		curr_offline: 65535,
		n:            65536,
	}
}

// Push the msg and response the heart beat
func (c *client) listen_loop() (e error) {

	var (
		err     error
		msg     *store.Msg
		pAndErr *packAndErr
		pack    *mqtt.Pack

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
			// Set the msg id
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
			mid := c.getOnlineMsgId()
			if mid == -1 {
				break
			}
			msg.Msg_id = mid
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
			glog.Infof("Client msg(%v)\n", pAndErr.pack.GetType())

			// Choose the requst type
			switch pAndErr.pack.GetType() {
			case mqtt.PUBACK:
				ack := pack.GetVariable().(*mqtt.Puback)
				if ack.GetMid() > 65535 {
					c.delMsg(OFFLINE, ack.GetMid())
				} else {
					// Online msg
					c.delMsg(ONLINE, ack.GetMid())
				}

			case mqtt.PINGREQ:
				// Reply the heart beat
				glog.Info("hb msg")
				err = c.queue.WritePack(mqtt.GetPingRespPack(1, msg.Dup))
				if err != nil {
					break loop
				}
			default:
				// Not define pack type
				glog.Errorf("The type not define:%v\n", pAndErr.pack.GetType())
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
	// Write the onlines msg to the db
	i, err := store.GetOfflineCount(c.id)
	if err != nil {
		glog.Errorf("Select warning: %v", err)
		goto passInsertOffline
	}
	if i == math.MaxUint16 {
		goto passInsertOffline
	}
	if i == 0 {
		i = 65536
	}
	for _, v := range c.onlineCache {
		// Add the offline msg id
		v.Msg_id = i
		v.Typ = OFFLINE
		store.InsertOfflineMsg(v)
		if i++; i > math.MaxUint16 {
			break
		}
	}

passInsertOffline:
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

// Setting a mqtt pack's id.
func (c *client) getOnlineMsgId() int {
	if c.curr_id == 65535 {
		i := 0
		for i < 65536-len(c.onlineCache) {
			if m := c.onlineCache[i]; m == nil {
				c.curr_id = i
				return c.curr_id
			}
			i++
		}
		return -1
	} else {
		id := c.curr_id
		for i := 0; i < 65535-c.curr_id; i++ {
			id++
			if m := c.onlineCache[id]; m == nil {
				c.curr_id = id
				return c.curr_id
			}
		}
		return -1
	}
}

func (c *client) pushMsg(msg *store.Msg) error {
	pack := mqtt.GetPubPack(1, msg.Dup, msg.Msg_id, &msg.Topic, msg.Body)
	// Write this pack
	err := c.queue.WritePack(pack)
	return err
}

func (c *client) delMsg(typ int, msg_id int) {
	switch typ {
	case OFFLINE:
		// Del the offline msg in the store
		glog.Info("Del a offline msg")
		store.DelOfflineMsg(msg_id, c.id)
	case ONLINE:
		// Del the online msg in the msg cache
		glog.Info("Del a online msg")
		delete(c.onlineCache, msg_id)
	default:
		glog.Errorf("No define the type:%v", typ)
	}
}

func (c *client) pushOfflineMsg(msg *store.Msg) (err error) {
	// The max cache size is 20
	if c.isOfflineClose && msg == nil {
		// Send
		err = c.queue.Flush()
	} else {
		// Wait
		if c.offlineNum > 20 {
			// Send
			pack := mqtt.GetPubPack(1, msg.Dup, msg.Msg_id, &msg.Topic, msg.Body)
			err = c.queue.WritePack(pack)
			// Clean the cache
			c.offlineNum = 0
		} else {
			c.offlineNum++
			pack := mqtt.GetPubPack(1, msg.Dup, msg.Msg_id, &msg.Topic, msg.Body)
			err = c.queue.WriteDelayPack(pack)
		}
	}
	return
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
