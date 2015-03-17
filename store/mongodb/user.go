// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"crypto/sha512"
	"fmt"
	"io"
	"strings"

	. "github.com/Alienero/quick-know/store/define"

	"labix.org/v2/mgo/bson"
)

func (mongo *Mongodb) Client_login(id, psw string) bool {
	sei := mongo.sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	var u = new(User)
	psw = getSalt(psw)
	it := c.Find(bson.M{"id": id, "psw": psw}).Iter()
	defer it.Close()
	if !it.Next(u) {
		return false
	}
	return true
}
func (mongo *Mongodb) Ctrl_login(id, auth string) (bool, string) {
	sei := mongo.sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	var u = new(Ctrl)
	auth = getSalt(auth)
	it := c.Find(bson.M{"auth": auth, "id": id}).Iter()
	defer it.Close()
	if !it.Next(u) {
		return false, ""
	}

	return true, u.Id
}

// Add or del user
func (mongo *Mongodb) AddUser(u *User) error {
	sei := mongo.sei_user.New()
	defer sei.Refresh()
	u.Id = Get_uuid()
	c := sei.DB(Config.UserName).C(Config.Clients)
	u.Psw = getSalt(u.Psw)
	return c.Insert(u)
}

func (mongo *Mongodb) DelUser(id string, own string) error {
	if !mongo.IsUserExist(id, own) {
		return fmt.Errorf("Del a user error:user not found,ID:%v,Own:%v", id, own)
	}
	sei_m := mongo.sei_msg.New()
	defer sei_m.Refresh()
	_, err := sei_m.DB(Config.MsgName).C(Config.SubsName).RemoveAll(bson.M{"user_id": id})
	if err != nil {
		return err
	}

	sei_msg := mongo.sei_msg.New()
	c := sei_msg.DB(Config.MsgName).C(Config.OfflineName)
	defer sei_msg.Refresh()
	if _, err = c.RemoveAll(bson.M{"m.to_id": id}); err != nil {
		return err
	}

	sei := mongo.sei_user.New()
	defer sei.Refresh()
	c = sei.DB(Config.UserName).C(Config.Clients)
	err = c.Remove(bson.M{"id": id, "owner": own})
	if err != nil {
		return err
	}
	return nil
}

func (mongo *Mongodb) IsUserExist(uid, oid string) bool {
	sei := mongo.sei_user.New()
	defer sei.Refresh()
	u := new(User)
	it := sei.DB(Config.UserName).C(Config.Clients).Find(bson.M{"id": uid, "owner": oid}).Iter()
	defer it.Close()
	return it.Next(u)
}

// Get the All use's id.
func (mongo *Mongodb) ChanUserID(own string) <-chan string {
	ch := make(chan string, 100)
	go func() {
		sei := mongo.sei_user.New()
		defer sei.Refresh()
		it := sei.DB(Config.UserName).C(Config.Clients).Find(bson.M{"owner": own}).Iter()
		u := new(User)
		for it.Next(u) {
			ch <- u.Id
		}
		it.Close()
		sei.Refresh()
		close(ch)
	}()

	return ch
}

func getSalt(s string) string {
	h := sha512.New()
	io.WriteString(h, s+Config.Salt)
	return strings.Replace(fmt.Sprintf("% x", h.Sum(nil)), " ", "", -1)
}
