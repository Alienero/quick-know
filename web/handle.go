package web

// The pack use default mux
import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Alienero/quick-know/comet"
	"github.com/Alienero/quick-know/store"

	"github.com/golang/glog"
)

func Init() error {
	// Init handle into mux
}

type handle struct {
	ID      string
	isBreak bool
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.prepare(w, r)
	if h.isBreak {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	switch r.Method {
	case "POST":
		h.Post(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
func (h *handle) prepare(w http.ResponseWriter, r *http.Request) {
	// Check the use name and password
	temp := r.Header.Get("Authorization")
	if len(temp) < 7 {
		h.isBreak = true
		return
	}
	auth := temp[7:]
	if b, id := store.Ctrl_login(auth); !b {
		h.isBreak = true
	} else {
		h.ID = id
	}
}
func (h *handle) Post(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r)
	if err != io.EOF {
		glog.Errorf("Read body error:%v\n", err)
		return
	}
	msg := new(store.Msg)
	msg.Owner = h.id
	err = json.Unmarshal(data, msg)
	if err != nil {
		glog.Errorf("Unmarshal json error%v\n", err)
		return
	}
	comet.WriteOnlineMsg(h.ID, msg)
}
