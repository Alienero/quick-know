// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"fmt"

	"github.com/Alienero/spp"

	// "github.com/golang/glog"
)

var ctrl_handles = make(map[int]handle)

// For call the serve
type handle func(c interface{}, pack *spp.Pack) error

func (h handle) serve(c interface{}, pack *spp.Pack) error {
	return h(c, pack)
}
func addHandle(typ int, f handle) {
	ctrl_handles[typ] = f
}

// Control server
type ControlServer struct {
	// queue *PackQueue
	id string
}

func newCServer(rw *spp.Conn, id string) *ControlServer {
	return &ControlServer{
		// queue: NewPackQueue(rw),
		id: id,
	}
}

func (cs *ControlServer) listen_loop() (err error) {
	defer func() {
		// Close the res
		// cs.queue.Close()
	}()
	var pack *spp.Pack
	for {
		// Listen
		// pack, err = cs.queue.ReadPack()
		// if err != nil {
		// glog.Errorf("clientLoop read pack error:%v\n", err)
		// break
		// }
		f := ctrl_handles[pack.Typ]
		if f == nil {
			err = fmt.Errorf("No such pack type:%v", pack.Typ)
			break
		}
		// Call function f
		err = f.serve(cs, pack)
		if err != nil {
			break
		}
	}
	return
}
