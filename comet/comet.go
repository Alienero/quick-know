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
	glog.Infof("Comet RPC Listen at:%v", Conf.RPC_addr)
	go listenRPC()
	// Open the cliens's server
	glog.Infof("Comet listen at: %v", Conf.Listen_addr)
	go startListen(CLIENT, Conf.Listen_addr)

	signal.HandleSignal(signal.InitSignal())
	glog.Info("System exit.")
}
