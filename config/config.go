// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

var (
	path     = flag.String("path", "", "-path=qk.conf")
	etcds    []string
	etcd_tmp string

	logger = log.New(os.Stdout, "qk_conf", log.Ltime|log.Lshortfile|log.LstdFlags)
)

func init() {
	flag.Parse()
	etcds = strings.Split(etcd_tmp, ",")
	data, err := ioutil.ReadFile(*path)
	if err != nil {
		logger.Panic(err)
	}
}

func main() {

}
