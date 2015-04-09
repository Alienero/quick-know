// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/Alienero/quick-know/mqtt"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/store/define"

	"github.com/golang/glog"
)

var notAlive = errors.New("Connection was dead")

type client struct {
	queue *PackQueue
	id    string

	offlines <-chan *define.Msg
	onlines  chan *define.Msg
	readChan <-chan *packAndErr

	onlineCache map[int]*define.Msg
	offlineNum  int // The max of once send the offline msgs.
	offline_map map[int]string

	closeChan   chan byte // Other gorountine Call notice exit
	isSendClose bool      // Wheather has a new login user.
	isLetClose  bool      // Wheather has relogin.

	isStop bool
	lock   *sync.Mutex

	isOfflineClose bool // Is offlines chan has been close

	// Online msg id
	curr_id int
	// The all not check msgs's number.
	counter int
}

func newClient(r *bufio.Reader, w *bufio.Writer, conn net.Conn, id string, alive int) *client {
	return &client{
		queue:     NewPackQueue(r, w, conn, alive),
		id:        id,
		closeChan: make(chan byte),
		lock:      new(sync.Mutex),

		offlines:    make(chan *define.Msg, Conf.MaxCacheMsg),
		onlines:     make(chan *define.Msg, Conf.MaxCacheMsg),
		offline_map: make(map[int]string),

		onlineCache: make(map[int]*define.Msg),

		curr_id: 0,
	}
}

// Push the msg and response the heart beat
func (c *client) listen_loop() (e error) {
	defer func() {
		err := redis.Logout(c.id, Conf.RPC_addr)
		if err != nil {
			err = redis.Logout(c.id, Conf.RPC_addr)
			if err != nil {
				glog.Errorf("Redis conn error:%v", err)
			}
			Users.Del(c.id)
			if c.isSendClose {
				c.closeChan <- 0
			}
		}
	}()
	var (
		err     error
		msg     *define.Msg
		pAndErr *packAndErr

		noticeFin = make(chan byte, 2)
		// wg        = new(sync.WaitGroup)

		findEndChan <-chan byte
	)

	// Start the write queue
	go c.queue.writeLoop()

	c.offlines, findEndChan = store.Manager.GetOfflineMsg(c.id, noticeFin)
	c.readChan = c.queue.ReadPackInLoop(noticeFin)

	// Start push
loop:
	for {
		if c.counter < math.MaxUint16 {
			select {
			case <-findEndChan:
				if err = c.offlineEnd(); err != nil {
					glog.Error(err)
					break loop
				}
			case msg_id := <-c.offlines:
				if err = c.waitOffline(msg_id); err != nil {
					glog.Error(err)
					break loop
				}
			case msg = <-c.onlines:
				if err = c.waitOnline(msg); err != nil {
					glog.Error(err)
					break loop
				}
			case pAndErr = <-c.readChan:
				if err = c.waitPack(pAndErr); err != nil {
					glog.Error("Get a connection error , will break(%v)", err)
					break loop
				}
			case <-c.closeChan:
				c.waitQuit()
				break loop
			}
		} else {
			select {
			case <-findEndChan:
				if err = c.offlineEnd(); err != nil {
					glog.Error(err)
					break loop
				}
			case pAndErr = <-c.readChan:
				if err = c.waitPack(pAndErr); err != nil {
					glog.Error(err)
					break loop
				}
			case <-c.closeChan:
				c.waitQuit()
				break loop
			}
		}

	}

	c.lock.Lock()
	c.isStop = true
	c.lock.Unlock()
	// Wrte the onlines msg to the db
	// Free resources
	// Close channels
	for i := 0; i < 2; i++ {
		noticeFin <- 1
	}
	close(noticeFin)
	// Write the onlines msg to the db
	for _, v := range c.onlineCache {
		// Add the offline msg id
		v.Typ = define.OFFLINE
		store.Manager.InsertOfflineMsg(v, Conf.RPC_addr, Conf.Etcd_addr)
	}
	// Flush the channel.
	for len(c.onlines) > 0 {
		msg := <-c.onlines
		msg.Typ = define.OFFLINE
		store.Manager.InsertOfflineMsg(msg, Conf.RPC_addr, Conf.Etcd_addr)
	}
	glog.Info("Cleaned the online msgs channel.")
	// Close the online msg channel
	close(c.onlines)
	close(c.closeChan)
	glog.Info("Groutine will esc.")
	return
}

// Select methods.
func (c *client) offlineEnd() (err error) {
	glog.Info("Offline has been finish send.")
	c.isOfflineClose = true
	err = c.pushOfflineMsg(nil)
	return
}

// Setting a mqtt pack's id.
func (c *client) getOnlineMsgId() int {
	if c.curr_id == math.MaxUint16 {
		if c.onlineCache[1] == nil {
			c.curr_id = 1
			return c.curr_id
		}
		return -1
	} else {
		if id := c.curr_id + 1; c.onlineCache[id] == nil {
			c.curr_id = id
			return c.curr_id
		}
		return -1
	}
}

func (c *client) waitOffline(msg *define.Msg) (err error) {
	// Check the msg time
	if msg.Expired > 0 {
		if time.Now().UTC().Unix() > msg.Expired {
			// cancel send the msg
			glog.Info("Out of time.")
			// TODO Del the offline msg
			store.Manager.DelOfflineMsg(msg.Id)
			return
		}
	}
	msg.Msg_id = c.getOnlineMsgId()
	if msg.Msg_id == -1 {
		return
	}
	c.offline_map[msg.Msg_id] = msg.Id
	// Push the offline msg
	glog.Info("Get a offline msg")
	// Set the msg id
	err = c.pushOfflineMsg(msg)
	c.counter++
	return
}

func (c *client) waitOnline(msg *define.Msg) (err error) {
	// Push the online msg
	glog.Info("Get a online msg")
	// Check the msg time
	if msg.Expired > 0 {
		if time.Now().UTC().Unix() > msg.Expired {
			// cancel send the msg
			glog.Info("Out of time")
			return
		}
	}
	// Add the msg into cache
	if msg == nil {
		err = errors.New("onlines has been close")
		return
	}
	if len(c.onlineCache) > Conf.MaxCacheMsg && Conf.MaxCacheMsg != 0 {
		err = fmt.Errorf("Online msg is out of range:%v", len(c.onlineCache))
		return
	}
	mid := c.getOnlineMsgId()
	if mid == -1 {
		return
	}
	msg.Msg_id = mid
	msg.Typ = define.ONLINE
	c.onlineCache[msg.Msg_id] = msg
	err = c.pushMsg(msg)
	c.counter++
	return
}

func (c *client) waitPack(pAndErr *packAndErr) (err error) {
	// If connetion has a error, should break
	// if it return a timeout error, illustrate
	// hava not recive a heart beat pack at an
	// given time.
	if pAndErr.err != nil {
		err = pAndErr.err
		return
	}
	glog.Infof("Client msg(%v)\n", pAndErr.pack.GetType())

	// Choose the requst type
	switch pAndErr.pack.GetType() {
	case mqtt.PUBACK:
		ack := pAndErr.pack.GetVariable().(*mqtt.Puback)
		// Del the msg
		c.delMsg(ack.GetMid())

	case mqtt.PINGREQ:
		// Reply the heart beat
		glog.Info("hb msg")
		err = c.queue.WritePack(mqtt.GetPingResp(1, pAndErr.pack.GetDup()))
	default:
		// Not define pack type
		err = fmt.Errorf("The type not define:%v\n", pAndErr.pack.GetType())
	}
	return
}

func (c *client) waitQuit() {
	// Start close
	glog.Info("Will break new relogin")
	c.isSendClose = true
}

func (c *client) delMsg(msg_id int) {
	if c.onlineCache[msg_id] != nil {
		// Del the online msg in the msg cache
		glog.Infof("Del a online msg id:%v", msg_id)
		delete(c.onlineCache, msg_id)
		c.counter--
	} else if c.offline_map[msg_id] != "" {
		// Del the offline msg in the store
		glog.Info("Del a offline msg")
		if err := store.Manager.DelOfflineMsg(c.offline_map[msg_id]); err != nil {
			glog.Error(err)
		}
		delete(c.offline_map, msg_id)
		c.counter--
	}
}

func (c *client) pushMsg(msg *define.Msg) error {
	pack := mqtt.GetPubPack(1, msg.Dup, msg.Msg_id, &msg.Topic, msg.Body)
	// Write this pack
	err := c.queue.WritePack(pack)
	return err
}

func (c *client) pushOfflineMsg(msg *define.Msg) (err error) {
	// The max cache size is 20
	if msg == nil && c.isOfflineClose {
		// Send
		err = c.queue.Flush()
	} else {
		// Wait
		if c.offlineNum > 20 || c.isOfflineClose {
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
