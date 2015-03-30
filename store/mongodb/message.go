// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"

	. "github.com/Alienero/quick-know/store/define"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

// id is to_id (client id)
func (mongo *Mongodb) GetOfflineMsg(id string, fin <-chan byte) (<-chan *Msg, <-chan byte) {
	// defer recover()
	// Find in the db
	ch := make(chan *Msg, Config.OfflineMsgs)
	ch2 := make(chan byte, 1)
	go func() {
		sei := mongo.sei_msg.New()
		c := sei.DB(Config.MsgName).C(Config.OfflineName)
		iter := c.Find(bson.M{"to_id": id}).Limit(Config.OfflineMsgs).Iter()
		msg := new(Msg)
		flag := false
	loop:
		for iter.Next(msg) {
			// Get sub's msg body.
			if msg.IsSub {
				sei := mongo.sei_msg.New()
				defer sei.Refresh()
				sc := sei.DB(Config.MsgName).C(Config.SubMsgName)
				subM := new(SubMsgs)
				if err := sc.Find(bson.M{"id": msg.Id}).One(&subM); err != nil {
					glog.Error(err)
					continue
				}
				msg.Body = subM.Body
			}
			select {
			case ch <- msg:
				msg = new(Msg)
			case <-fin:
				// No read the all offline msg, notice close
				glog.Info("DB FIN")
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
			glog.Info("DB FIN,2")
		}
		close(ch)
		close(ch2)
	}()
	return ch, ch2
}

// Should delete.
// func (mongo *Mongodb) GetOfflineCount(id string) (int, error) {
// 	c := mongo.sei_msg.DB(Config.MsgName).C(Config.OfflineName)
// 	defer mongo.sei_msg.Refresh()
// 	msg := new(Msg)
// 	if err := c.Find(bson.M{"to_id": id}).Sort("msg_id", "-1").One(&msg); err != nil {
// 		return 0, err
// 	}
// 	return msg.Msg_id, nil
// }

// Del the offile msg
func (mongo *Mongodb) DelOfflineMsg(id string) error {
	sei := mongo.sei_msg.New()
	c := sei.DB(Config.MsgName).C(Config.OfflineName)
	defer sei.Refresh()
	result := new(Msg)
	if err := c.Find(bson.M{"id": id}).One(result); err != nil {
		return err
	}
	if result.IsSub {
		// Delete offline sub msg.
		sei := mongo.sei_msg.New()
		sc := mongo.sei_msg.DB(Config.MsgName).C(Config.SubMsgName)
		defer sei.Refresh()
		count := new(SubMsgs)
		q := bson.M{"id": id}
		if err := sc.Find(q).One(count); err != nil {
			return err
		}
		if count.Count == 1 {
			// Remove this doc.
			if err := sc.Remove(q); err != nil {
				return err
			}
		} else {
			if err := sc.Update(q, bson.M{"count": "$dec"}); err != nil {
				return err
			}
		}

	}
	err := c.Remove(bson.M{"m.id": id})
	if err != nil {
		return fmt.Errorf("Remove a offline msg(id:%v) error:%v", id, err)
	}
	return nil
}

// Intert a new offilne msg
// Before should check the to_id belong the user
func (mongo *Mongodb) InsertOfflineMsg(msg *Msg) error {
	return mongo.insert(msg, false)
}

// TODO:
func (mongo *Mongodb) InsertSubOfflineMsg(msg *Msg, subId string) error {
	return nil
}

func (mongo *Mongodb) insert(msg *Msg, isSub bool) error {
	sei := mongo.sei_msg.New()
	c := sei.DB(Config.MsgName).C(Config.OfflineName)
	defer sei.Refresh()
	if msg.IsSub {
		//
	}

	err := c.Insert(&msg)
	if err != nil {
		return fmt.Errorf("Intert a offline msg(id:%v) error:%v", msg.Id)
	}
	return nil
}

// TODO sub msg insert and find.
