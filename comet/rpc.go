// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"time"

	myrpc "github.com/Alienero/quick-know/rpc"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/store/define"

	"github.com/golang/glog"
)

type Comet_RPC struct {
}

func (*Comet_RPC) Relogin(id string, r *myrpc.Reply) error {
	c := Users.Get(id)
	if c == nil {
		// r.IsRe = true
	} else {
		c.lock.Lock()
		if !c.isLetClose {
			c.isLetClose = true
			c.lock.Unlock()
			select {
			case c.CloseChan <- 1:
				<-c.CloseChan
				r.IsOk = true
				glog.Info("RPC: Ok will be relogin.")
			case <-time.After(2 * time.Second):
				// Timeout.
				if c := Users.Get(id); c != nil {
					return errors.New("Close the logon user timeout")
				}
				// Has been esc.
			}
		} else {
			c.lock.Unlock()
		}
	}
	return nil
}

func (*Comet_RPC) WriteOnlineMsg(msg *define.Msg, r *myrpc.Reply) (err error) {
	defer func() {
		if err == nil {
			r.IsOk = true
		} else {
			glog.Error(err)
		}
	}()
	glog.Info("Get a Write Online msg RPC")

	c := Users.Get(msg.To_id)
	if c == nil {
		msg.Typ = define.OFFLINE
		// Get the offline msg id
		err = store.Manager.InsertOfflineMsg(msg, Conf.RPC_addr, Conf.Etcd_addr)
		return
	}

	c.lock.Lock()
	if len(c.onlines) == Conf.MaxCacheMsg {
		c.lock.Unlock()
		msg.Typ = define.OFFLINE
		err = store.Manager.InsertOfflineMsg(msg, Conf.RPC_addr, Conf.Etcd_addr)
		return
	} else {
		c.lock.Unlock()
	}
	c.lock.Lock()
	if c.isStop {
		c.lock.Unlock()
		msg.Typ = define.OFFLINE
		err = store.Manager.InsertOfflineMsg(msg, Conf.RPC_addr, Conf.Etcd_addr)
	} else {
		msg.Typ = define.ONLINE
		c.onlines <- msg
		c.lock.Unlock()
	}
	return
}

func (c *Comet_RPC) WriteMsg(msg *define.Msg, r *myrpc.Reply) (err error) {
	msg.Dup = 0
	// Fix the Expired
	if msg.Expired > 0 {
		msg.Expired = time.Now().UTC().Add(time.Duration(msg.Expired)).Unix()
	}
	// Check the user whether online.
	b, addr, err := redis.IsExist(msg.To_id)
	if !b {
		if err != nil {
			glog.Error(err)
		}
		// User is offline.
		err = store.Manager.InsertOfflineMsg(msg, Conf.RPC_addr, Conf.Etcd_addr)
		if err != nil {
			glog.Error(err)
		} else {
			// nil error.
			r.IsOk = true
		}
		return err
	}
	// User is online.
	if addr == Conf.RPC_addr {
		// Call local comet.
		return c.WriteOnlineMsg(msg, r)
	} else {
		// RPC.
		c, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			return err
		}
		reply := new(myrpc.Reply)
		if err = c.Call("Comet_RPC.WriteOnlineMsg", msg, reply); err != nil {
			return err
		}
		r.IsRe = reply.IsRe
		r.IsOk = reply.IsOk
		return nil
	}
}

func (*Comet_RPC) Ping(total *int, r *myrpc.Reply) (err error) {
	r.IsOk = true
	l := Users.Len()
	total = &l
	return nil
}

func listenRPC() {
	comet := new(Comet_RPC)
	if err := rpc.Register(comet); err != nil {
		glog.Fatal(err)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", Conf.RPC_addr)
	if err != nil {
		glog.Fatal(err)
	}
	if err = http.Serve(l, nil); err != nil {
		glog.Error(err)
	}
}
