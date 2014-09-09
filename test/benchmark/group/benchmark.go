// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"runtime"
	"strconv"
	"time"

	"github.com/Alienero/quick-know/comet"
	tool "github.com/Alienero/quick-know/restful_tool"
	// "github.com/Alienero/quick-know/signal"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/web"

	"github.com/golang/glog"
)

var (
	psw = "123"
	ids = make([]string, 2000)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer glog.Flush()
	// Start server
	// run_server()
	// Add user
	if err := add_user(); err != nil {
		glog.Fatal(err)
	}
	// Client Login
	for i := 0; i < cap(ids); i++ {
		go client(ids[i])
	}
	// Push msg
	go func() {
		for i := 0; i < 100; i++ {
			glog.Infof("ready for push %v times\n", i)
			if err := tool.PushMsg2All([]byte("benchmark test:"+strconv.Itoa(i)), 0); err != nil {
				glog.Info("push API get error:%v", err)
				continue
			}
			glog.Infof("finish for push %v times\n", i)
		}
	}()
	time.Sleep(10 * time.Minute)
	// exit
	glog.Info("Sysytem exit")
}

func add_user() error {
	tool.ID = "1234"
	tool.Psw = "10086"
	for i := 0; i < cap(ids); i++ {
		if id, err := tool.AddUser(psw); err != nil {
			return err
		} else {
			ids[i] = id
		}
	}
	return nil
}

func run_server() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	defer glog.Flush()
	glog.Infoln("Server loding!")
	// Init the DB conf
	glog.Infoln("Loading store")
	if err := store.Init(); err != nil {
		glog.Fatal(err)
	}
	glog.Infoln("Loading comet server")
	if err := comet.InitConf(); err != nil {
		glog.Fatal(err)
	}
	glog.Infoln("Loading web server")
	if err := web.InitConf(); err != nil {
		glog.Fatal(err)
	}

	go web.Start()
	glog.Infoln("Web server has been started!")

	go comet.Start()
	glog.Infoln("Comet server has been started!")

}
