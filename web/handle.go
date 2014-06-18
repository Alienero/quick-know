// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"io"
	"net/http"

	"github.com/Alienero/quick-know/comet"
	"github.com/Alienero/quick-know/store"

	"github.com/golang/glog"
)

// Push a private msg
func private_msg(w http.ResponseWriter, r *http.Request, u *user) {
	msg := new(store.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.IsUserExist(msg.To_id, u.ID) {
		msg.Owner = u.ID
		msg.Msg_id = get_uuid()
		comet.WriteOnlineMsg(msg)
		io.WriteString(w, `{msg_id":"`)
		io.WriteString(w, msg.Msg_id)
		io.WriteString(w, `"}`)
	} else {
		io.WriteString(w, `{msg_id":"`)
		// io.WriteString(w, msg.Msg_id)
		io.WriteString(w, `"}`)
	}
}

// Add a new user
func add_user(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Add a new user(Client)")
	u := new(store.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	u.Owner = uu.ID
	u.Id = get_uuid()
	store.AddUser(u)
	io.WriteString(w, `{id":"`)
	io.WriteString(w, u.Id)
	io.WriteString(w, `"}`)
}

// Delete a user
func del_user(w http.ResponseWriter, r *http.Request, uu *user) {
	u := new(store.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	if err := store.DelUser(u.Id, uu.ID); err != nil {
		glog.Errorf("Del user in the web error:%v", err)
		io.WriteString(w, `{id":"`)
		io.WriteString(w, u.Id)
		io.WriteString(w, `"}`)
	} else {
		io.WriteString(w, `{Status":"Fail"}`)
	}
}

// Add sub group
func add_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sub := new(store.Sub)
	err := readAdnGet(r.Body, sub)
	if err != nil {
		glog.Errorf("Get a new sub error%v\n", err)
		return
	}
	sub.Id = get_uuid()
	store.AddSub(sub)
	// Write the response
	io.WriteString(w, `{sub_id":"`)
	io.WriteString(w, sub.Id)
	io.WriteString(w, `"}`)
}

// Del sub msg
func del_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sub := new(store.Sub)
	err := readAdnGet(r.Body, sub)
	if err != nil {
		glog.Errorf("Get a new sub error%v\n", err)
		return
	}

	if err := store.DelSub(sub.Id, uu.ID); err != nil {
		// Write the response
		glog.Errorf("Del sub in the web error:%v", err)
		io.WriteString(w, `{sub_id":"`)
		io.WriteString(w, sub.Id)
		io.WriteString(w, `"}`)
	} else {
		io.WriteString(w, `{Status":"Fail"}`)
	}
}

// Add use into msg's sub group
func user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sm := new(store.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the add user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, uu.ID, sm.Sub_id, err)
		return
	}
	// Write the response
	if err = store.AddUserToSub(sm, uu.ID); err != nil {
		glog.Errorf("Store the sub_map error:", err)
		io.WriteString(w, `{Status":"Fail"}`)
	} else {
		io.WriteString(w, `{Status":"OK"}`)
	}
}

// Remove user from the sub group
func rm_user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sm := new(store.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the remove user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, uu.ID, sm.Sub_id, err)
		return
	}
	if err = store.DelUserFromSub(sm.Sub_id, sm.User_id, uu.ID); err != nil {
		glog.Errorf("Del the user o error:%v\n", err)
		io.WriteString(w, `{Status":"Fail"}`)
	} else {
		io.WriteString(w, `{Status":"OK"}`)
	}
}

// Send msg to all
func broadcast(w http.ResponseWriter, r *http.Request, uu *user) {
	msg := new(store.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.IsUserExist(msg.To_id, uu.ID) && msg.To_id == "" {
		msg.Owner = uu.ID
		msg.Msg_id = get_uuid()

		ch := store.ChanUserID(uu.ID)
		go func() {
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				msg.To_id = s
				comet.WriteOnlineMsg(msg)
			}
		}()
		io.WriteString(w, `{msg_id":"`)
		io.WriteString(w, msg.Msg_id)
		io.WriteString(w, `"}`)
	} else {
		io.WriteString(w, `{Status":"Fail"}`)
		// io.WriteString(w, msg.Msg_id)
		// io.WriteString(w, `"}`)
	}
}

// Send msg to sub group
func group_msg(w http.ResponseWriter, r *http.Request, uu *user) {
	type multi_cast struct {
		Sub_id string
		Msg    *store.Msg
	}
	mc := new(multi_cast)
	err := readAdnGet(r.Body, mc)
	if err != nil {
		glog.Errorf("Get sub msg error:%v\n", err)
		return
	}
	// Check the sub group belong to the user
	if store.IsSubExist(mc.Sub_id, uu.ID) {
		mc.Msg.Msg_id = get_uuid()
		// Submit msg
		go func() {
			ch := store.ChanSubUsers(mc.Sub_id)
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				mc.Msg.To_id = s
				comet.WriteOnlineMsg(mc.Msg)
			}
		}()
		io.WriteString(w, `{msg_id":"`)
		io.WriteString(w, mc.Msg.Msg_id)
		io.WriteString(w, `"}`)
	} else {
		io.WriteString(w, `{Status":"Fail"}`)
	}
}
