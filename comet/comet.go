// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Alienero/quick-know/signal"
	"github.com/Alienero/quick-know/store"

	"github.com/golang/glog"
)

func main() {
	glog.Info("Comet Server:Loading...")
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

	// Open the RPC listener.
	glog.Info("Comet RPC Listener Start.")
	go listenRPC()
	// Open the cliens's server
	glog.Info("Comet client Listener Start.")
	go startListen(CLIENT, Conf.Listen_addr)
	// Reg the comet server, and open the etcd keeper.
	glog.Info("Comet etcd Start.")
	if err := Init_etcd(); err != nil {
		glog.Fatal(err)
	}

	sotre_conf := ""
	if Conf.From_etcd {
		// Get the Store conf.
		var err error
		sotre_conf, err = GetStore()
		if err != nil {
			panic(err)
		}
	}
	glog.Infoln("Loading store")
	if err := store.Init(sotre_conf); err != nil {
		glog.Fatal(err)
	}

	signal.HandleSignal(signal.InitSignal())
	glog.Info("System exit.")
}
