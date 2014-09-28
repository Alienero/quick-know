// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Alienero/quick-know/comet"
	"github.com/Alienero/quick-know/signal"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/web"

	"github.com/golang/glog"
)

func main() {
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

	// init signals, block wait signals
	signal.HandleSignal(signal.InitSignal())
	// exit
	glog.Info("Sysytem exit")
}
