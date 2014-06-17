// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"flag"

	"github.com/Alienero/quick-know/store"
)

func Start() {
	flag.Parse()
	// Init the DB conf
	if err := store.Init(); err != nil {
		panic(err)
	}
	if err := InitConf(); err != nil {
		panic(err)
	}
	// Open the cliens's server
	if err := startListen(CLIENT, Conf.Listen_addr); err != nil {
		panic(err)
	}
}
