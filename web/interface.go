// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

// The pack use default mux
import (
	"net/http"

	"github.com/Alienero/quick-know/store"
)

func Init() error {
	// Init handle into mux
	return nil
}

type user struct {
	ID      string
	isBreak bool
}

type handle struct {
}

func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := new(user)
	h.prepare(w, r, u)
	if u.isBreak {
		http.Error(w, "", http.StatusForbidden)
		return
	}
	switch r.Method {
	case "POST":
		h.Post(w, r, u)
	default:
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
	auth := temp[7:]
	if b, id := store.Ctrl_login(auth); !b {
		u.isBreak = true
	} else {
		u.ID = id
	}
}
func (h *handle) Post(w http.ResponseWriter, r *http.Request, u *user) {
}
