// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

// {"id":"","psw":"","owner":""}
type User struct {
	Id    string
	Psw   string
	Owner string // Owner
}
type Ctrl struct {
	Id   string
	Auth string
}

func Client_login(id, psw string) bool {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	var u = new(User)
	it := c.Find(bson.M{"id": id, "psw": psw}).Iter()
	defer it.Close()
	if !it.Next(u) {
		return false
	}
	return true
}
func Ctrl_login(id, auth string) (bool, string) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	var u = new(Ctrl)
	it := c.Find(bson.M{"auth": auth, "id": id}).Iter()
	defer it.Close()
	if !it.Next(u) {
		return false, ""
	}

	return true, u.Id
}

func Ctrl_login_alive(id, psw string) bool { return false }

// Add or del user
func AddUser(u *User) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	err := c.Insert(u)
	if err != nil {
		glog.Errorf("Insert a new user error:%v", err)
	}
}
func DelUser(id string, own string) error {
	if !IsUserExist(id, own) {
		return fmt.Errorf("Del a user error:user not found,ID:%v,Own:%v", id, own)
	}
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	err := c.Remove(bson.M{"id": id, "owner": own})
	if err != nil {
		glog.Errorf("Del a user error:%v,ID:%v,Own:%v", err, id, own)
		return err
	}
	return nil
}

func IsUserExist(uid, oid string) bool {
	sei := sei_user.New()
	defer sei.Refresh()
	u := new(User)
	it := sei.DB(Config.UserName).C(Config.Clients).Find(bson.M{"id": uid, "owner": oid}).Iter()
	defer it.Close()
	return it.Next(u)
}

func ChanUserID(own string) <-chan string {
	ch := make(chan string, 100)
	go func() {
		sei := sei_user.New()
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
