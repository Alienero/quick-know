// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"net/http"

	"github.com/golang/glog"
)

func Start() {
	// Add the mux into the web server
	http.Handle("/push/private", &handle{Post: private_msg})
	http.Handle("/push/add_user", &handle{Post: add_user})
	http.Handle("/push/del_user", &handle{Post: del_user})
	http.Handle("/push/add_sub", &handle{Post: add_sub})
	http.Handle("/push/del_sub", &handle{Post: del_sub})
	http.Handle("/push/user_sub", &handle{Post: user_sub})
	http.Handle("/push/rm_user_sub", &handle{Post: rm_user_sub})
	http.Handle("/push/all", &handle{Post: broadcast})
	http.Handle("/push/group_msg", &handle{Post: group_msg})

	glog.Infof("Listen at port :%v", Conf.Listen_addr)
	glog.Error(http.ListenAndServe(Conf.Listen_addr, nil))
}
