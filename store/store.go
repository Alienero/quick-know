// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	. "github.com/Alienero/quick-know/store/define"
	"github.com/Alienero/quick-know/store/mongodb"
)

// A db manage
var Manager DataStorer

type DataStorer interface {
	// Msg's methods.

	// Get the offline msg from the database which you choosed.
	GetOfflineMsg(id string, fin <-chan byte) (<-chan *Msg_id, <-chan byte)
	// Get the muber of the offline msgs.
	// GetOfflineCount(id string) (int, error)
	// Delete the offline msg.
	DelOfflineMsg(id string) error
	// Insert a msg to the database.
	InsertOfflineMsg(msg *Msg) error

	// User

	// Check the ctrl use login.
	Ctrl_login(id, auth string) (bool, string)
	// Check the user login.
	Client_login(id, psw string) bool
	// Add a client user.
	AddUser(u *User) error
	// Del a client user.
	DelUser(id string, own string) error
	// Check a user wheather belong a ctrl user.
	IsUserExist(uid, oid string) bool
	// Get the whole of the ctrl's users.
	ChanUserID(own string) <-chan string

	// Sub

	// Add a sub.
	AddSub(sub *Sub) error
	// Del a sub.
	DelSub(sub_id, id string) error
	// Add a user into a sub
	AddUserToSub(sm *Sub_map, id string) error
	// Del a user from a sub
	DelUserFromSub(sm *Sub_map, id string) error
	// Check a sub weather exist.
	IsSubExist(sub_id, id string) bool
	// Get a sub all users
	ChanSubUsers(sub_id string) <-chan string
}

func Init(s string) (err error) {
	err = InitConfig(s)
	if err != nil {
		return
	}
	switch Config.DBType {
	case "mongodb":
		Manager, err = mongodb.NewMongo()
	default:
		err = fmt.Errorf("no such db type(%v) support ", Config.DBType)
	}
	return
}
