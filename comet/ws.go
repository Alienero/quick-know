// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/net/websocket"
	"net/http"
)

const (
	WsPath = "/pub"
)

func WsHandle(ws *websocket.Conn) {

}

// If http.ListenandServe return an error,
// it will throws a panic.
func wsListener(addr string, tls bool) {
	if err := func() error {
		httpServeMux := http.NewServeMux()
		httpServeMux.Handle("/pub", websocket.Handler(WsHandle))
		// if tls {
		// 	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
		// }
		return http.ListenAndServe(addr, httpServeMux)
	}(); err != nil {
		panic(err)
	}
}
