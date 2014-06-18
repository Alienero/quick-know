package comet

import (
	"github.com/golang/glog"
)

func Start() {
	// Open the cliens's server
	if err := startListen(CLIENT, Conf.Listen_addr); err != nil {
		glog.Fatal(err)
	}
}
