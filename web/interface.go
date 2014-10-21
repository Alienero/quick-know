// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

// The pack use default mux
import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Alienero/quick-know/store"
	"github.com/golang/glog"
)

type user struct {
	ID      string
	isBreak bool

	isOK bool
}

type handler func(w http.ResponseWriter, r *http.Request, uu *user)

type handle struct {
	f      handler
	method string
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != h.method {
		glog.Info("Undefine Method.")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
	u := new(user)
	h.prepare(w, r, u)
	if u.isBreak {
		glog.Info("Login fail")
		http.Error(w, "", http.StatusForbidden)
		return
	}
	glog.Infof("Do the %v method", r.Method)
	h.f(w, r, u)
}
func (h *handle) prepare(w http.ResponseWriter, r *http.Request, u *user) {
	// Check the use name and password
	temp := r.Header.Get("Authorization")
	if len(temp) < 7 {
		u.isBreak = true
		return
	}
	auth := strings.TrimLeft(temp, " ")[6:]
	buf, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		u.isBreak = true
		glog.Errorf("Decode string(%v) from base64 error:%v", auth, err)
		return
	}
	kp := string(buf)
	index := strings.Index(kp, ":")
	if index < 0 || index > len(kp) {
		u.isBreak = true
		glog.Error("auth basic string out of range")
		return
	}
	if b, id := store.Manager.Ctrl_login(kp[:index], kp[index+1:]); !b {
		u.isBreak = true
	} else {
		u.ID = id
	}
}

func Handle(path, method string, h handler) {
	http.Handle(path, &handle{
		f:      h,
		method: method})
}
