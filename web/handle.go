// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"io"
	"net/http"

	"github.com/Alienero/quick-know/comet"
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
	glog.Info("Add a private msg")
	msg := new(define.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.Manager.IsUserExist(msg.To_id, u.ID) {
		msg.Owner = u.ID
		// msg.Msg_id = get_uuid()
		msg.Topic = Private + msg.Owner
		comet.WriteOnlineMsg(msg)
		io.WriteString(w, `{"status":"success"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Add a new user
func add_user(w http.ResponseWriter, r *http.Request, uu *user) {
	glog.Info("Add a new user(Client)")
	u := new(define.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	u.Owner = uu.ID
	u.Id = get_uuid()
	if err := store.Manager.AddUser(u); err != nil {
		badReaquest(w, `{"status":"fail"}`)
	} else {
		io.WriteString(w, `{"id":"`)
		io.WriteString(w, u.Id)
		io.WriteString(w, `"}`)
	}
}

// Delete a user
func del_user(w http.ResponseWriter, r *http.Request, uu *user) {
	u := new(define.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	if err := store.Manager.DelUser(u.Id, uu.ID); err != nil {
		glog.Errorf("Del user in the web error:%v", err)
		io.WriteString(w, `{"id":"`)
		io.WriteString(w, u.Id)
		io.WriteString(w, `"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Add sub group
func add_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sub := new(define.Sub)
	err := readAdnGet(r.Body, sub)
	if err != nil {
		glog.Errorf("Get a new sub error%v\n", err)
		return
	}
	sub.Id = get_uuid()
	sub.Own = uu.ID
	if err := store.Manager.AddSub(sub); err != nil {
		badReaquest(w, `{"status":"fail"}`)
	} else {
		io.WriteString(w, `{"sub_id":"`)
		io.WriteString(w, sub.Id)
		io.WriteString(w, `"}`)
	}
}

// Del sub msg
func del_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sub := new(define.Sub)
	err := readAdnGet(r.Body, sub)
	if err != nil {
		glog.Errorf("Get a new sub error%v\n", err)
		return
	}

	if err := store.Manager.DelSub(sub.Id, uu.ID); err != nil {
		// Write the response
		glog.Errorf("Del sub in the web error:%v", err)
		io.WriteString(w, `{"sub_id":"`)
		io.WriteString(w, sub.Id)
		io.WriteString(w, `"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Add use into msg's sub group
func user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sm := new(define.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the add user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, uu.ID, sm.Sub_id, err)
		return
	}
	// Write the response
	if err = store.Manager.AddUserToSub(sm, uu.ID); err != nil {
		glog.Errorf("Store the sub_map error:", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		io.WriteString(w, `{"status":"success"}`)
	}
}

// Remove user from the sub group
func rm_user_sub(w http.ResponseWriter, r *http.Request, uu *user) {
	sm := new(define.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the remove user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, uu.ID, sm.Sub_id, err)
		return
	}
	if err = store.Manager.DelUserFromSub(sm, uu.ID); err != nil {
		glog.Errorf("Del the user o error:%v\n", err)
		badReaquest(w, `{"status":"fail"}`)
	} else {
		io.WriteString(w, `{"status":"success"}`)
	}
}

// Send msg to all
func broadcast(w http.ResponseWriter, r *http.Request, uu *user) {
	msg := new(define.Msg)
	msg.Topic = Inf_All
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push inform msg error%v\n", err)
		return
	}
	if store.Manager.IsUserExist(msg.To_id, uu.ID) && msg.To_id == "" {
		msg.Owner = uu.ID
		// msg.Msg_id = get_uuid()

		ch := store.Manager.ChanUserID(uu.ID)
		go func() {
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				m := *msg
				m.To_id = s
				comet.WriteOnlineMsg(&m)
			}
		}()
		io.WriteString(w, `{"status":"success"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

// Send msg to sub group
func group_msg(w http.ResponseWriter, r *http.Request, uu *user) {
	type multi_cast struct {
		Sub_id string
		Msg    *define.Msg
	}
	mc := new(multi_cast)
	err := readAdnGet(r.Body, mc)
	if err != nil {
		glog.Errorf("Get sub msg error:%v\n", err)
		return
	}
	// Check the sub group belong to the user
	if store.Manager.IsSubExist(mc.Sub_id, uu.ID) {
		// mc.Msg.Msg_id = get_uuid()
		mc.Msg.Topic = Group + mc.Sub_id
		// Submit msg
		go func() {
			ch := store.Manager.ChanSubUsers(mc.Sub_id)
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				msg := *mc.Msg
				msg.To_id = s
				comet.WriteOnlineMsg(&msg)
			}
		}()
		io.WriteString(w, `{"status":"success"}`)
	} else {
		badReaquest(w, `{"status":"fail"}`)
	}
}

func badReaquest(w http.ResponseWriter, err string) {
	http.Error(w, err, http.StatusBadRequest)
}
