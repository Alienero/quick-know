// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	. "github.com/Alienero/quick-know/store/define"

	"labix.org/v2/mgo"
)

type Mongodb struct {
	sei_user *mgo.Session
	sei_msg  *mgo.Session
}

func NewMongo() (mongo *Mongodb, err error) {
	mongo = new(Mongodb)
	if err = connect(&mongo.sei_user, config.UserAddr); err != nil {
		return
	}
	err = connect(&mongo.sei_msg, config.MsgAddr)
	return
}

func connect(sei **mgo.Session, addr string) (err error) {
	if *sei != nil {
		sei.Close()
	}
	*sei, err = mgo.Dial(addr)
	sei.EnsureSafe(&mgo.Safe{})
	sei.SetMode(mgo.Monotonic, true)
	sei.Refresh()
	return
}
