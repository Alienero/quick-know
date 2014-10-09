// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"

	. "github.com/Alienero/quick-know/store/define"

	"labix.org/v2/mgo/bson"
)

// id is to_id (client id)
func (mongo *Mongodb) GetOfflineMsg(id string, fin <-chan byte) (<-chan *Msg_id, <-chan byte) {
	// defer recover()
	// Find in the db
	ch := make(chan *Msg_id, Config.OfflineMsgs)
	ch2 := make(chan byte, 1)
	go func() {
		sei := mongo.sei_msg.New()
		c := sei.DB(Config.MsgName).C(Config.OfflineName)
		iter := c.Find(bson.M{"m.to_id": id}).Limit(Config.OfflineMsgs).Iter()
		msg_id := new(Msg_id)
		flag := false
	loop:
		for iter.Next(msg_id) {
			select {
			case ch <- msg_id:
				msg_id = new(Msg_id)
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

// Should delete.
func (mongo *Mongodb) GetOfflineCount(id string) (int, error) {
	c := mongo.sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer mongo.sei_msg.Refresh()
	msg := new(Msg)
	if err := c.Find(bson.M{"to_id": id}).Sort("msg_id", "-1").One(&msg); err != nil {
		return 0, err
	}
	return msg.Msg_id, nil
}

// Del the offile msg
func (mongo *Mongodb) DelOfflineMsg(id string) error {
	c := mongo.sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer mongo.sei_msg.Refresh()
	err := c.Remove(bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("Remove a offline msg(id:%v) error:%v", id, err)
	}
	return nil
}

// Intert a new offilne msg
// Before should check the to_id belong the user
func (mongo *Mongodb) InsertOfflineMsg(msg *Msg) error {
	c := mongo.sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer mongo.sei_msg.Refresh()
	id := Get_uuid()
	err := c.Insert(&Msg_id{
		M:  msg,
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("Intert a offline msg(id:%v) error:%v", id)
	}
	return nil
}
