// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restful_tool

import (
	"testing"
)

func Init() {
	Url = "http://127.0.0.1:9901"
	ID = "1234"
	Psw = "10086"
}

var id string
var sub_id string

func TestAddUser(t *testing.T) {
	Init()
	var err error
	id, err = AddUser("123456")
	if err != nil {
		t.Error(err)
	}
}

func TestAddSub(t *testing.T) {
	var err error
	if sub_id, err = AddSub(0, 0); err != nil {
		t.Error(err)
	}
	println("Sub id is :", sub_id)
}

func TestAddUser2Sub(t *testing.T) {
	if err := User2Sub(sub_id, id); err != nil {
		t.Error(err)
	}
}

func TestGroupMsg(t *testing.T) {
	if err := GroupMsg(sub_id, 0, []byte("GroupMsg")); err != nil {
		t.Error(err)
	}
}

func TestSendAll(t *testing.T) {
	us := NewUrls("/push/all")
	defer us.ReFresh()
	if err := Broadcast(0, []byte("All msg!")); err != nil {
		t.Error(err)
	}
}

func TestAddPrivateMsg(t *testing.T) {
	if err := AddPrivateMsg(id, 0, []byte("Private msg")); err != nil {
		t.Error(err)
	} else {
		println("Pass")
	}
}

func TestRmUserSub(t *testing.T) {
	if err := RmUserSub(sub_id, id); err != nil {
		t.Error(err)
	}
}

func TestDelSub(t *testing.T) {
	if err := DelSub(sub_id); err != nil {
		t.Error(err)
	}
}

func TestDelUser(t *testing.T) {
	if err := DelUser(id); err != nil {
		t.Error(err)
	}
}
