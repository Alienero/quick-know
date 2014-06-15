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
type private_msg struct {
	handle
}

func (this *private_msg) Post(w http.ResponseWriter, r *http.Request) {
	msg := new(store.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.IsUserExist(msg.ToID, this.ID) {
		msg.Owner = this.ID
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
type add_user struct {
	handle
}

func (this *add_user) Post(w http.ResponseWriter, r *http.Request) {
	u := new(store.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	u.Id = get_uuid()
	store.AddUser(u)
	io.WriteString(w, `{id":"`)
	io.WriteString(w, u.Id)
	io.WriteString(w, `"}`)
}

// Delete a user
type del_user struct {
	handle
}

func (this *del_user) Post(w http.ResponseWriter, r *http.Request) {
	u := new(store.User)
	err := readAdnGet(r.Body, u)
	if err != nil {
		glog.Errorf("add user error%v\n", err)
		return
	}
	store.DelUser(u.Id, this.ID)
	io.WriteString(w, `{id":"`)
	io.WriteString(w, u.Id)
	io.WriteString(w, `"}`)
}

// Add sub msg
type add_sub struct {
	handle
}

func (this *add_sub) Post(w http.ResponseWriter, r *http.Request) {
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
type del_sub struct {
	handle
}

func (this *del_sub) Post(w http.ResponseWriter, r *http.Request) {
	sub := new(store.Sub)
	err := readAdnGet(r.Body, sub)
	if err != nil {
		glog.Errorf("Get a new sub error%v\n", err)
		return
	}
	store.DelSub(sub.Id, this.ID)
	// Write the response
	io.WriteString(w, `{sub_id":"`)
	io.WriteString(w, sub.Id)
	io.WriteString(w, `"}`)
}

// Add use into msg's sub group
type sub_msg struct {
	handle
}

func (this *sub_msg) Post(w http.ResponseWriter, r *http.Request) {
	sm := new(store.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the add user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, this.ID, sm.Sub_id, err)
		return
	}
	// Write the response
	if err = store.AddUserToSub(sm, this.ID); err != nil {
		glog.Errorf("Store the sub_map error:", err)
		io.WriteString(w, `{Status":"Fail"}`)
	} else {
		io.WriteString(w, `{Status":"OK"}`)
	}
}

// Remove ues from the msg
type rm_msg_sub struct {
	handle
}

func (this *rm_msg_sub) Post(w http.ResponseWriter, r *http.Request) {
	sm := new(store.Sub_map)
	err := readAdnGet(r.Body, sm)
	if err != nil {
		glog.Errorf("Get the remove user(%v) of id(%v) to sub(%v) error:%v\n", sm.User_id, this.ID, sm.Sub_id, err)
		return
	}
	if err = store.DelUserFromSub(sm.Sub_id, sm.User_id, this.ID); err != nil {
		glog.Errorf("Del the user o error:%v\n", err)
		io.WriteString(w, `{Status":"Fail"}`)
	} else {
		io.WriteString(w, `{Status":"OK"}`)
	}
}

// Send msg to all
type broadcast struct {
	handle
}

func (this *broadcast) Post(w http.ResponseWriter, r *http.Request) {
	msg := new(store.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	if store.IsUserExist(msg.ToID, this.ID) && msg.ToID == "" {
		msg.Owner = this.ID
		msg.Msg_id = get_uuid()

		ch := store.ChanUserID(this.ID)
		go func() {
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				msg.ToID = s
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
type group_msg struct {
	handle
}

func (this *group_msg) Post(w http.ResponseWriter, r *http.Request) {
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
	if store.IsSubExist(mc.Sub_id, this.ID) {
		mc.Msg.Msg_id = get_uuid()
		// Submit msg
		go func() {
			ch := store.ChanSubUsers(mc.Sub_id)
			for {
				s, ok := <-ch
				if !ok {
					break
				}
				mc.Msg.ToID = s
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
