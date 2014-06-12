package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
	data, err := ioutil.ReadAll(r.Body)
	if err != io.EOF {
		glog.Errorf("Read body error:%v\n", err)
		return
	}
	msg := new(store.Msg)
	msg.Owner = m.ID
	err = json.Unmarshal(data, msg)
	if err != nil {
		glog.Errorf("Unmarshal json error%v\n", err)
		return
	}
	msg.Msg_id = get_uuid()
	comet.WriteOnlineMsg(m.ID, msg)
	io.WriteString(w, `{msg_id":"`)
	io.WriteString(w, msg.Msg_id)
	io.WriteString(w, `"}`)
}

// Add a new user
type add_user struct {
}

func (this *add_user) Post(w http.ResponseWriter, r *http.Request) {

}

// Delete a user
type del_user struct {
}

func (this *del_user) Post(w http.ResponseWriter, r *http.Request) {

}

// Add sub msg
type add_sub struct {
}

func (this *add_sub) Post(w http.ResponseWriter, r *http.Request) {

}

// Del sub msg
type del_sub struct {
}

func (this *del_sub) Post(w http.ResponseWriter, r *http.Request) {

}

// Add use into msg
type sub_msg struct {
}

func (this *sub_msg) Post(w http.ResponseWriter, r *http.Request) {

}

// Remove ues from the msg
type rm_msg_sub struct {
}

func (this *rm_msg_sub) Post(w http.ResponseWriter, r *http.Request) {

}

// Send msg to all
type broadcast struct {
}

func (this *broadcast) Post(w http.ResponseWriter, r *http.Request) {

}
