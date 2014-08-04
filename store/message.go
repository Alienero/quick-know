// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"time"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

type Msg struct {
	Msg_id int    // Msg ID
	Owner  string // Owner
	To_id  string // Topic
	Body   []byte
	Typ    int

	Dup byte // mqtt dup

	Expired int64
}

func GetOfflineMsg(id string, fin <-chan byte) (<-chan *Msg, <-chan byte) {
	// defer recover()
	// Find in the db
	ch := make(chan *Msg, Config.OfflineMsgs)
	ch2 := make(chan byte, 1)
	go func() {
		sei := sei_msg.New()
		c := sei.DB(Config.MsgName).C(Config.OfflineName)
		iter := c.Find(bson.M{"to_id": id}).Limit(Config.OfflineMsgs).Iter()
		msg := new(Msg)
		flag := false
		// Check time expired
	loop:
		for iter.Next(msg) {
			if msg.Expired > 0 {
				if time.Now().UTC().Unix() > msg.Expired {
					// Delet the offline msg in the BD
					DelOfflineMsg(msg.Msg_id, id)
					continue
				}
			}
			select {
			case ch <- msg:
				msg = new(Msg)
			case <-fin:
				// No read the all offline msg, notice close
				flag = true
				break loop
			}

		}

		iter.Close()
		sei.Refresh()

		ch2 <- 1
		if !flag {
			// not recive the fin. must wait the fin
			<-fin
		}
		close(ch)
		close(ch2)
	}()
	return ch, ch2
}

// Del the offile msg
func DelOfflineMsg(msg_id int, id string) {
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
