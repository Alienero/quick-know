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
	Handle("/push/private", "PUT", private_msg)
	Handle("/push/add_user", "PUT", add_user)
	Handle("/push/del_user", "DELETE", del_user)
	Handle("/push/add_sub", "PUT", add_sub)
	Handle("/push/del_sub", "DELETE", del_sub)
	Handle("/push/user_sub", "PUT", user_sub)
	Handle("/push/rm_user_sub", "DELETE", rm_user_sub)
	Handle("/push/group_msg", "PUT", group_msg)
	Handle("/push/all", "PUT", broadcast)

	glog.Infof("Listen at port :%v", Conf.Listen_addr)
	glog.Error(http.ListenAndServe(Conf.Listen_addr, nil))
}
