// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
)

const (
	WsPath = "/pub"
)

func WsHandle(ws *websocket.Conn) {

}

// If http.ListenandServe return an error,
// it will throws a panic.
func wsListener() {
	if err := func() error {
		httpServeMux := http.NewServeMux()
		httpServeMux.Handle("/pub", websocket.Handler(WsHandle))
		var (
			l   net.Listener
			err error
		)
		if Conf.Tls {
			tlsConf := new(tls.Config)
			tlsConf.Certificates = make([]tls.Certificate, 1)
			tlsConf.Certificates[0], err = tls.X509KeyPair(Conf.Cert, Conf.Key)
			if err != nil {
				return err
			}
			l, err = tls.Listen("tcp", Conf.WebSocket_addr, tlsConf)
			if err != nil {
				return err
			}
			return http.Serve(l, httpServeMux)
		} else {
			return http.ListenAndServe(Conf.WebSocket_addr, httpServeMux)
		}
	}(); err != nil {
		panic(err)
	}
}
