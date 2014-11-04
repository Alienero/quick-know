// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"time"

	myrpc "github.com/Alienero/quick-know/rpc"

	"github.com/golang/glog"
)

type Comet_RPC struct {
}

// func (*comet_RPC) Lock(id string, r *rpc.Reply) error {
// 	c := Users.Get(id)
// 	if c != nil {
// 		c.lock.Lock()
// 	}
// 	return nil
// }

// func (*comet_RPC) Unlock(id string, r *rpc.Reply) error {
// 	c := Users.Get(id)
// 	if c != nil {
// 		c.lock.Unlock()
// 	}
// 	return nil
// }

func (*Comet_RPC) Relogin(id string, r *myrpc.Reply) error {
	c := Users.Get(id)
	if c == nil {
		r.IsRe = true
	} else {
		c.lock.Lock()
		if !c.isLetClose {
			c.isLetClose = true
			c.lock.Unlock()
			select {
			case c.CloseChan <- 1:
				<-c.CloseChan
				r.IsOk = true
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
