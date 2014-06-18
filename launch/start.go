package main

import (
	"flag"
	"runtime"

	"github.com/Alienero/quick-know/comet"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/web"

	"github.com/golang/glog"
)

func main() {
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
	if err := web.Init(); err != nil {
		glog.Fatal(err)
	}

	go web.Start()
	glog.Infoln("Web server has been started!")

	comet.Start()
	glog.Infoln("Comet server has been started!")
}
