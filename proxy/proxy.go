// Copyright © 2015 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
	"runtime"
	"sync"

	"github.com/Alienero/quick-know/signal"

	"github.com/golang/glog"
)

func main() {
	p := newProxy()
	go p.listen()
	signal.HandleSignal(signal.InitSignal())
	p.Wait()
}

type Proxy struct {
	user_map map[string][]string
	conf     *config
	*sync.WaitGroup
}

func newProxy() Proxy {
	return Proxy{}
}

func (proxy *Proxy) listen() {
	l, err := net.Listen("tcp", proxy.conf.ListenAddr)
	if err != nil {
		glog.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Temporary() {
					continue
				}
			}
			glog.Errorf("http: Accept error: %v", err)
			return
		}
		proxy.Add(1)
		go proxy.serve(conn)
	}
}

func (proxy *Proxy) serve(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 4096)
			runtime.Stack(buff, false)
			glog.Errorf("conn.serve() panic(%v)\n info:%s", err, string(buff))
		}
		conn.Close()
		proxy.Done()
	}()
}
