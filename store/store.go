// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

type Msg struct {
	Msg_id int    // Msg ID
	Owner  string // Owner
	To_id  string
	Topic  string
	Body   []byte
	Typ    int // Online or Oflline msg

	Dup byte // mqtt dup

	Expired int64
}

type User struct {
	Id    string
	Psw   string
	Owner string // Owner
}

type Ctrl struct {
	Id   string
	Auth string
}

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

type DataStorer interface {
	// Msg's methods.

	// Get the offline msg from the database which you choosed.
	GetOfflineMsg(id string, fin <-chan byte) (<-chan *Msg, <-chan byte)
	// Get the muber of the offline msgs.
	GetOfflineCount(id string) (int, error)
	// Delete the offline msg.
	DelOfflineMsg(msg_id int, id string)
	// Insert a msg to the database.
	InsertOfflineMsg(msg *Msg)

	// User

	// Check the ctrl use login.
	Ctrl_login(id, auth string) (bool, string)
	// Check the user login.
	Client_login(id, psw string) bool
	// Check the Ctrl is logon.
	Ctrl_login_alive(id, psw string) bool
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
