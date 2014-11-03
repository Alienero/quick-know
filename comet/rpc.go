// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"net/rpc"
	"time"

	"github.com/Alienero/quick-know/rpc"
)

type comet_RPC struct {
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

func (*comet_RPC) Relogin(id string, r *rpc.Reply) error {
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
					return "Close the logon user timeout"
				}
				// Has been esc.
			}
		} else {
			c.lock.Unlock()
		}
	}
	return nil
}
