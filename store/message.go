// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	// "time"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

const (
	OFFLINE = 11
	ONLINE  = 12
)

type Msg struct {
	Msg_id string // Msg ID
	Owner  string // Owner
	To_id  string
	Body   []byte
	Typ    int

	Expired int64
}

// BUG[#1]: if channel has been clsoe , but client server has benn not recive the nil msg
// it will send the fin(1) to the channel and wait. The stack has been not exist, so it will
// become a dead lock.
func GetOfflineMsg(id string, fin <-chan byte) <-chan *Msg {
	// defer recover()
	// Find in the db
	ch := make(chan *Msg, Config.OfflineMsgs)
	go func() {
		sei := sei_msg.New()
		defer sei.Refresh()
		c := sei.DB(Config.MsgName).C(Config.OfflineName)
		iter := c.Find(bson.M{"to_id": id}).Limit(Config.OfflineMsgs).Iter()
		defer iter.Close()
		msg := new(Msg)
		flag := false
	loop:
		for iter.Next(msg) {
			select {
			case ch <- msg:
				msg = new(Msg)
			case <-fin:
				flag = true
				break loop
			}

		}
		if !flag {
			// not recive the fin. wait the fin
			<-fin
		}
		close(ch)
	}()
	return ch
}

// Del the offile msg
func DelOfflineMsg(msg_id string, id string) {
	c := sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer sei_msg.Refresh()
	err := c.Remove(bson.M{"msg_id": msg_id, "to_id": id})
	if err != nil {
		glog.Errorf("Remove a offline msg(id:%v,to_id:%v) error:%v", msg_id, id, err)
	}
}

// Intert a new offilne msg
// Before should check the to_id belong the user
func InsertOfflineMsg(msg *Msg) {
	c := sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer sei_msg.Refresh()
	err := c.Insert(msg)
	if err != nil {
		glog.Errorf("Intert a offline msg(id:%v) error:%v", msg.Msg_id)
	}
}
