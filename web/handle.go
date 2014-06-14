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

func (m *private_msg) Post(w http.ResponseWriter, r *http.Request) {
	msg := new(store.Msg)
	err := readAdnGet(r.Body, msg)
	if err != nil {
		glog.Errorf("push private msg error%v\n", err)
		return
	}
	msg.Owner = m.ID
	msg.Msg_id = get_uuid()
	comet.WriteOnlineMsg(m.ID, msg)
	io.WriteString(w, `{msg_id":"`)
	io.WriteString(w, msg.Msg_id)
	io.WriteString(w, `"}`)
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

// Add use into msg
type sub_msg struct {
	handle
}

func (this *sub_msg) Post(w http.ResponseWriter, r *http.Request) {

}

// Remove ues from the msg
type rm_msg_sub struct {
	handle
}

func (this *rm_msg_sub) Post(w http.ResponseWriter, r *http.Request) {

}

// Send msg to all
type broadcast struct {
	handle
}

func (this *broadcast) Post(w http.ResponseWriter, r *http.Request) {

}

// Send msg to sub group
type group_msg struct {
	handle
}

func (this *group_msg) Post(w http.ResponseWriter, r *http.Request) {

}
