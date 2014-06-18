package web

import (
	"net/http"
)

func Start() {
	// Add the mux into the web server
	http.Handle("/push/private", &private_msg{})
	http.Handle("/push/add_user", &add_user{})
	http.Handle("/push/del_user", &del_user{})
	http.Handle("/push/add_sub", &add_sub{})
	http.Handle("/push/del_sub", &del_sub{})
	http.Handle("/push/user_sub", &user_sub{})
	http.Handle("/push/rm_user_sub", &rm_user_sub{})
	http.Handle("/push/all", &broadcast{})
	http.Handle("/push/group_msg", &group_msg{})

	http.ListenAndServe(Conf.Listen_addr, nil)
	// http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}
