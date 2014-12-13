// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// Import config define
	comet "github.com/Alienero/quick-know/comet/config"
	store "github.com/Alienero/quick-know/store/define"
	web "github.com/Alienero/quick-know/web/config"

	"github.com/coreos/go-etcd/etcd"
)

var (
	path     = flag.String("path", "", "-path=qk.conf")
	etcd_tmp = flag.String("etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")

	logger = log.New(os.Stdout, "qk_conf", log.Ltime|log.Lshortfile|log.LstdFlags)

	Conf = config{}

	etcdClient = *etcd.Client
)

func init() {
	flag.Parse()
	// Init etcd.
	etcdClient = etcd.NewClient(strings.Split(*etcd_tmp, ","))
}

type config struct {
	Comet struct {
		comet.Etcd
		comet.Redis
		comet.Restriction
	}
	Web   web.Config
	Store store.DBConfig
}

func main() {
	// Read config.
	if err := readFile(*path); err != nil {
		logger.Panic(err)
	}
	// Share config.
}

func readFile(path string) error {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	return json.Unmarshal(data, &Conf)
}

func setNode(node string, vaule string) error {

}
