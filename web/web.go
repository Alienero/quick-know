// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Alienero/quick-know/signal"
	"github.com/Alienero/quick-know/store"

	"github.com/golang/glog"
)

func start() {
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
	go func() {
		err := http.ListenAndServe(Conf.Listen_addr, nil)
		glog.Fatal(err)
	}()
}

func main() {
	glog.Info("Web Server:Loading...")
	b := flag.Bool("benchmark", false, "")
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *b {
		glog.Info("Benchmark Mode")
		// Creat a file
		f, err := os.Create("pprof")
		if err != nil {
			glog.Fatal(err)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			glog.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
	defer glog.Flush()

	glog.Info("Read the config.")
	if err := InitConf(); err != nil {
		glog.Fatal(err)
	}

	glog.Info("Web listener start.")
	go start()

	glog.Info("Web etcd start.")
	if err := Init_etcd(); err != nil {
		glog.Fatal(err)
	}

	glog.Infoln("Loading store")
	sotre_conf := ""
	if Conf.From_etcd {
		// Get the Store conf.
		var err error
		sotre_conf, err = GetStore()
		if err != nil {
			panic(err)
		}
	}
	if err := store.Init(sotre_conf); err != nil {
		glog.Fatal(err)
	}
	signal.HandleSignal(signal.InitSignal())
}
