// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"fmt"
	"testing"

	"labix.org/v2/mgo/bson"
)

func insertAdmin() {
	u := &Ctrl{"1234", "10086"}
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	err := c.Insert(u)
	if err != nil {
		println(err.Error())
	}
}
func TestInit(t *testing.T) {
	if err := Init(); err != nil {
		t.Errorf(err.Error())
	} else {
		if sei_user == nil {
			t.Logf("sei_user has no vaule")
		}
	}
	if sei_user == nil {
		t.Logf("sei_user is empty")
	}
	fmt.Println(Config)
	insertAdmin()
}

func TestInsertUser(t *testing.T) {
	AddUser(&User{
		Id:    "1001",
		Owner: "1234",
		Psw:   "102",
	})
}

func TestLogin(t *testing.T) {
	if b := Client_login("1001", "102"); !b {
		t.Errorf("Error login")
	}
}

func TestCtrlLogin(t *testing.T) {
	if b, id := Ctrl_login("1234", "10086"); !b || id != "1234" {
		t.Errorf("Error login")
	}
}

// Msg test -------start
func TestInsertMsg(t *testing.T) {
	// TODO To_id belong
	InsertOfflineMsg(&Msg{
		Msg_id: 10081,
		Owner:  "1234",
		To_id:  "1001",
		Body:   []byte("Hello word!"),
	})
}

func TestDelMsg(t *testing.T) {
	// TODO To_id belong
	DelOfflineMsg(10081, "1001")
}

func TestIsUserExist(t *testing.T) {
	if !IsUserExist("1001", "1234") {
		t.Errorf("Exist Error")
	}
	if IsUserExist("1001", "12") {
		t.Errorf("Exist Error")
	}
}

// Msg test --------stop

// Subs test --------start
func TestSub(t *testing.T) {
	AddSub(&Sub{
		Id:  "12",
		Own: "1234",
	})
	// AddUserToSub(sub_id, id)
	if err := AddUserToSub(&Sub_map{"12", "1001"}, "1234"); err != nil {
		t.Error(err)
	}
	if err := DelUserFromSub(&Sub_map{"12", "1001"}, "1234"); err != nil {
		t.Error(err)
	}
	if err := DelSub("12", "1234"); err != nil {
		t.Error(err)
	}
}

// Subs test ---------stop

func TestDelUser(t *testing.T) {
	DelUser("1001", "1234")
}

func TestDel(t *testing.T) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	err := c.Remove(bson.M{"id": "1234"})
	if err != nil {
		t.Error(err)
	}
}
