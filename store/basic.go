// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"labix.org/v2/mgo"
	// "labix.org/v2/mgo/bson"
)

var sei_user *mgo.Session
var sei_msg *mgo.Session

func Init() (err error) {
	err = initConfig()
	if err != nil {
		return
	}
	err = connect(sei_user, Config.UserAddr)
	if err != nil {
		return
	}
	err = connect(sei_msg, Config.MsgAddr)
	return
}

func connect(sei *mgo.Session, addr string) (err error) {
	if sei != nil {
		sei.Close()
	}
	sei, err = mgo.Dial(addr)
	sei.EnsureSafe(&mgo.Safe{})
	sei.SetMode(mgo.Monotonic, true)
	sei.Refresh()
	return
}
