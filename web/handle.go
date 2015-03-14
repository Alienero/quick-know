// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/store/define"

	"github.com/golang/glog"
)

var (
	// [type]/|[.../][who send]
	Private = "m/"
	Group   = "ms/"
	Inf_All = "inf/all"
)

// Push a private msg
func private_msg(w http.ResponseWriter, r *http.Request, u *user) {
	glog.Info("Add a private msg.")
	msg := newMsg()
	r.ParseForm()
	msg.To_id = r.FormValue("to_id")
	glog.Infof("To_id is :%v", msg.To_id)
	if s := r.FormValue("expired"); s != "" {
		msg.Expired, _ = strconv.ParseInt(s, 10, 64)
	}
	var err error
	if msg.Body, err = ioutil.ReadAll(r.Body); err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.Manager.IsUserExist(msg.To_id, u.ID) {
		msg.Owner = u.ID
		// msg.Msg_id = get_uuid()
		msg.Topic = Private + msg.Owner
		if err = write_msg(msg); err != nil {
			glog.Error(err)
			badReaquest(w, `{"status":"fail"}`)
			return
		}
		io.WriteString(w, `{"status":"success"}`)
		u.isOK = true
	} else {
		glog.Info("push private msg error: user not exist.")
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Add a new user
func add_user(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Add a new user(Client)")
	u := new(define.User)
	r.ParseForm()
	u.Psw = r.FormValue("psw")
	u.Owner = uu.ID
	if err := store.Manager.AddUser(u); err != nil {
		badReaquest(w, `{"status":"fail"}`)
	} else {
		uu.isOK = true
		io.WriteString(w, `{"id":"`)
		io.WriteString(w, u.Id)
		io.WriteString(w, `"}`)
	}
}

// Delete a user
func del_user(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Del a user.")
	u := new(define.User)
	r.ParseForm()
	u.Id = r.FormValue("id")
	if err := store.Manager.DelUser(u.Id, uu.ID); err != nil {
		glog.Errorf("Del user in the web error:%v", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		uu.isOK = true
		io.WriteString(w, `{"id":"`)
		io.WriteString(w, u.Id)
		io.WriteString(w, `"}`)
	}
}

// Add sub group
func add_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Add a new sub.")
	sub := new(define.Sub)
	r.ParseForm()
	if s := r.FormValue("max"); s != "" {
		sub.Max, _ = strconv.Atoi(s)
	}
	if s := r.FormValue("type"); s != "" {
		sub.Typ, _ = strconv.Atoi(s)
	}
	// sub.Id = get_uuid()
	sub.Own = uu.ID
	if err := store.Manager.AddSub(sub); err != nil {
		glog.Errorf("Add new sub error(%v)", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		glog.Info("Add sub ok.")
		uu.isOK = true
		io.WriteString(w, `{"sub_id":"`)
		io.WriteString(w, sub.Id)
		io.WriteString(w, `"}`)
	}
}

// Del sub msg
func del_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Del a sub.")
	sub := new(define.Sub)
	r.ParseForm()
	sub.Id = r.FormValue("id")
	if err := store.Manager.DelSub(sub.Id, uu.ID); err != nil {
		// Write the response
		glog.Errorf("Del sub in the web error:%v", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		uu.isOK = true
		io.WriteString(w, `{"status":"success"}`)
	}
}

// Add use into msg's sub group
func user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Del a user sub.")
	sm := new(define.Sub_map)
	r.ParseForm()
	sm.Sub_id = r.FormValue("sub_id")
	sm.User_id = r.FormValue("user_id")
	// Write the response
	if err := store.Manager.AddUserToSub(sm, uu.ID); err != nil {
		glog.Errorf("Store the sub_map error:", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		uu.isOK = true
		io.WriteString(w, `{"status":"success"}`)
	}
}

// Remove user from the sub group
func rm_user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("rm a user from a sub.")
	sm := new(define.Sub_map)
	r.ParseForm()
	sm.Sub_id = r.FormValue("sub_id")
	sm.User_id = r.FormValue("user_id")
	if err := store.Manager.DelUserFromSub(sm, uu.ID); err != nil {
		glog.Errorf("Del the user o error:%v\n", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		uu.isOK = true
		io.WriteString(w, `{"status":"success"}`)
	}
}

// Send msg to all
func broadcast(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Send a msg to all.")
	msg := newMsg()
	msg.Topic = Inf_All
	r.ParseForm()
	if s := r.FormValue("expired"); s != "" {
		msg.Expired, _ = strconv.ParseInt(s, 10, 64)
	}
	var err error
	if msg.Body, err = ioutil.ReadAll(r.Body); err != nil {
		glog.Errorf("push inform msg error%v\n", err)
		return
	}
	if msg.To_id == "" {
		msg.Owner = uu.ID
		ch := store.Manager.ChanUserID(uu.ID)
		go func() {
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				m := *msg
				m.To_id = s
				write_msg(&m)
			}
		}()
		uu.isOK = true
		io.WriteString(w, `{"status":"success"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Send msg to sub group
func group_msg(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("add a msg to a group.")
	msg := newMsg()
	r.ParseForm()
	sub_id := r.FormValue("sub_id")
	if s := r.FormValue("expired"); s != "" {
		msg.Expired, _ = strconv.ParseInt(s, 10, 64)
	}
	var err error
	if msg.Body, err = ioutil.ReadAll(r.Body); err != nil {
		glog.Errorf("push inform msg error%v\n", err)
		return
	}
	// Check the sub group belong to the user
	if store.Manager.IsSubExist(sub_id, uu.ID) {
		// mc.Msg.Msg_id = get_uuid()
		msg.Topic = Group + sub_id
		// Submit msg
		go func() {
			ch := store.Manager.ChanSubUsers(sub_id)
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				m := *msg
				m.To_id = s
				write_msg(&m)
			}
		}()
		uu.isOK = true
		io.WriteString(w, `{"status":"success"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

func badReaquest(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}

func newMsg() *define.Msg {
	id, err := getID()
	if err != nil {
		panic(err)
	}
	return &define.Msg{
		Id: id,
	}
}
