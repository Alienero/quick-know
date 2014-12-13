// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	// "bufio"
	// "bytes"
	"encoding/json"
	"flag"
	// "os"
	"strings"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/web/config"
)

var Conf = config.Config{}

// func InitConf() error {
// 	buf := new(bytes.Buffer)
// 	f, err := os.Open("web.conf")
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	r := bufio.NewReader(f)
// 	for {
// 		line, err := r.ReadSlice('\n')
// 		if err != nil {
// 			if len(line) > 0 {
// 				buf.Write(line)
// 			}
// 			break
// 		}
// 		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
// 			buf.Write(line)
// 		}
// 	}
// 	return json.Unmarshal(buf.Bytes(), Conf)
// }

var etcd_addr_temp string

func init() {
	flag.StringVar(&Conf.Listen_addr, "listen", "", "-listen=127.0.0.1:9002")
	flag.BoolVar(&Conf.Tls, "tls", false, "-tls=true")
	flag.StringVar(&etcd_addr_temp, "etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")
}

func InitConf() error {
	// Read config frome etcd.
	// Init etcd client.
	Conf.Etcd_addr = strings.Split(etcd_addr_temp, ",")
	init_etcd()
	// Get store config.
	storeConf, err := GetStore()
	if err != nil {
		return err
	}
	if err := store.Init([]byte(storeConf)); err != nil {
		return err
	}
	// Get web config.
	if webConf, err := GetWeb(); err != nil {
		return err
	} else {
		if err = json.Unmarshal([]byte(webConf), &Conf); err != nil {
			return err
		}
	}

	return etcd_hb()
}
