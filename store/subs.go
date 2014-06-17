// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

type Sub struct {
	Id  string
	Own string

	Max int // Max of the group members
	Typ int
}
type Sub_map struct {
	Sub_id  string
	User_id string
}

func AddSub(sub *Sub) {
	sei := sei_msg.New()
	defer sei.Refresh()
	c := sei.DB(Config.MsgName).C(Config.SubName)
	err := c.Insert(sub)
	if err != nil {
		glog.Errorf("Insert a new sub group error:%v\n", err)
	}
}

func DelSub(sub_id, id string) error {
	if IsSubExist(sub_id, id) {
		sei := sei_msg.New()
		defer sei.Refresh()
		c := sei.DB(Config.MsgName).C(Config.SubsName)
		_, err := c.RemoveAll(bson.M{"sub_id": sub_id})
		if err != nil {
			glog.Errorf("Del a sub group's all members error:%v\n", err)
			return err
		}
		c = sei.DB(Config.MsgName).C(Config.SubName)
		err = c.Remove(bson.M{"id": sub_id})
		if err != nil {
			glog.Errorf("Del a sub group error:%v\n", err)
			return err
		}
		return nil
	}
	return fmt.Errorf("the sub(%v) of the user(%v) is not exist", sub_id, id)
}

func AddUserToSub(sm *Sub_map, id string) error {
	// Check the sub has been exist and belong to the use
	sei := sei_msg.New()
	defer sei.Refresh()

	if IsSubExist(sm.Sub_id, id) && IsUserExist(sm.User_id, id) {
		if err := sei.DB(Config.MsgName).C(Config.SubsName).Insert(sm); err != nil {
			glog.Errorf("Insert a new sub's user error:%v\n", err)
		}
		return nil
	} else {
		return fmt.Errorf("sub(%v) or user(%v) of id(%v) not exist ", sm.Sub_id, sm.User_id, id)
	}
}

func DelUserFromSub(sub_id, uid, id string) error {
	if IsUserExist(uid, id) && IsSubExist(sub_id, id) {
		sei := sei_msg.New()
		defer sei.Refresh()
		err := sei.DB(Config.MsgName).C(Config.SubsName).Remove(bson.M{"sub_id": sub_id, "user_id": uid})
		if err != nil {
			glog.Errorf("Remove the user(%v) from the sub(%v) error:%v\n", uid, sub_id, err)
		}
		return nil
	}
	return fmt.Errorf("sub(%v) or user(%v) of the id(%v) not exist", sub_id, uid, id)
}

func IsSubExist(sub_id, id string) bool {
	sei := sei_msg.New()
	defer sei.Refresh()
	it := sei.DB(Config.MsgName).C(Config.SubName).Find(bson.M{"id": sub_id, "own": id}).Iter()
	defer it.Close()
	sub := new(Sub)
	return it.Next(sub)
}

func ChanSubUsers(sub_id string) <-chan string {
	ch := make(chan string, 100)

	go func() {
		sei := sei_msg.New()
		it := sei.DB(Config.MsgName).C(Config.SubsName).Find(bson.M{"sub_id": sub_id}).Iter()
		sm := new(Sub_map)
		for it.Next(sm) {
			ch <- sm.User_id
		}
		close(ch)
		it.Close()
		sei.Refresh()
	}()

	return ch
}
