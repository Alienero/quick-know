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
		}
	}()
	glog.Infof("Get a Write msg RPC,msg is :%v", string(msg.Body))
	msg.Dup = 0
	// fix the Expired
	if msg.Expired > 0 {
		msg.Expired = time.Now().UTC().Add(time.Duration(msg.Expired)).Unix()
	}

	c := Users.Get(msg.To_id)
	if c == nil {
		msg.Typ = OFFLINE
		// Get the offline msg id
		err = store.Manager.InsertOfflineMsg(msg)
		return
	}

	c.lock.Lock()
	if len(c.onlines) == Conf.MaxCacheMsg {
		c.lock.Unlock()
		msg.Typ = OFFLINE
		err = store.Manager.InsertOfflineMsg(msg)
		return
	} else {
		c.lock.Unlock()
	}
	c.lock.Lock()
	if c.isStop {
		c.lock.Unlock()
		msg.Typ = OFFLINE
		err = store.Manager.InsertOfflineMsg(msg)
	} else {
		msg.Typ = ONLINE
		c.onlines <- msg
		c.lock.Unlock()
	}
	return
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
