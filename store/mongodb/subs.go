// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"

	. "github.com/Alienero/quick-know/store/define"

	"labix.org/v2/mgo/bson"
)

func (mongo *Mongodb) AddSub(sub *Sub) error {
	sei := mongo.sei_msg.New()
	defer sei.Refresh()
	sub.Id = Get_uuid()
	c := sei.DB(Config.MsgName).C(Config.SubName)
	return c.Insert(sub)
}

func (mongo *Mongodb) DelSub(sub_id, id string) error {
	if mongo.IsSubExist(sub_id, id) {
		sei := mongo.sei_msg.New()
		defer sei.Refresh()
		c := sei.DB(Config.MsgName).C(Config.SubsName)
		_, err := c.RemoveAll(bson.M{"sub_id": sub_id})
		if err != nil {
			return err
		}
		c = sei.DB(Config.MsgName).C(Config.SubName)
		err = c.Remove(bson.M{"id": sub_id})
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("the sub(%v) of the user(%v) is not exist", sub_id, id)
}

func (mongo *Mongodb) AddUserToSub(sm *Sub_map, id string) error {
	// Check the sub has been exist and belong to the use
	sei := mongo.sei_msg.New()
	defer sei.Refresh()

	if mongo.IsSubExist(sm.Sub_id, id) && mongo.IsUserExist(sm.User_id, id) {
		return sei.DB(Config.MsgName).C(Config.SubsName).Insert(sm)
	} else {
		return fmt.Errorf("sub(%v) or user(%v) of id(%v) not exist ", sm.Sub_id, sm.User_id, id)
	}
}

func (mongo *Mongodb) DelUserFromSub(sm *Sub_map, id string) error {
	if mongo.IsUserExist(sm.User_id, id) && mongo.IsSubExist(sm.Sub_id, id) {
		sei := mongo.sei_msg.New()
		defer sei.Refresh()
		return sei.DB(Config.MsgName).C(Config.SubsName).Remove(bson.M{"sub_id": sm.Sub_id, "user_id": sm.User_id})
	}
	return fmt.Errorf("sub(%v) or user(%v) of the id(%v) not exist", sm.Sub_id, sm.User_id, id)
}

func (mongo *Mongodb) IsSubExist(sub_id, id string) bool {
	sei := mongo.sei_msg.New()
	defer sei.Refresh()
	it := sei.DB(Config.MsgName).C(Config.SubName).Find(bson.M{"id": sub_id, "own": id}).Iter()
	defer it.Close()
	sub := new(Sub)
	return it.Next(sub)
}

func (mongo *Mongodb) ChanSubUsers(sub_id string) <-chan string {
	ch := make(chan string, 100)

	go func() {
		sei := mongo.sei_msg.New()
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
