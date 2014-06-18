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
}

type handle struct {
	Post func(w http.ResponseWriter, r *http.Request, uu *user)
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := new(user)
	h.prepare(w, r, u)
	if u.isBreak {
		glog.Info("Login fail")
		http.Error(w, "", http.StatusForbidden)
		return
	}
	switch r.Method {
	case "POST":
		h.Post(w, r, u)
	default:
		glog.Info("No define method")
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
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
	if b, id := store.Ctrl_login(string(buf)); !b {
		u.isBreak = true
	} else {
		u.ID = id
	}
}

// func (h *handle) Post(w http.ResponseWriter, r *http.Request, uu *user) {
// }
