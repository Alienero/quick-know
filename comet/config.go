// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
)

var Conf = config{}

type config struct {
	Listenner struct {
		Listen_addr    string `json:"-"` // Client listener addr
		WebSocket_addr string `json:"-"`
		RPC_addr       string `json:"-"`
		Tls            bool   `json:"-"`
	}

	Restriction struct {
		WirteLoopChanNum int // Should > 1
		ReadPackLoop     int
		MaxCacheMsg      int
		ReadTimeout      int // Heart beat check (seconds)
		WriteTimeout     int
	}

	Redis struct {
		// Redis conf
		Network    string
		Address    string
		MaxIde     int
		IdeTimeout int // Second.
	}

	Etcd struct {
		// Etcd conf.
		Etcd_addr     []string `json:"-"`
		Etcd_interval uint64
		Etcd_dir      string
	}
}

func confFromFile(path string) error {
	buf := new(bytes.Buffer)

	f, err := os.Open("comet.conf")
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
	return json.Unmarshal(buf.Bytes(), Conf)
}

var (
	path      string
	etcd_addr string
)

func init() {
	flag.StringVar(&path, "path", "", "-path=comet.conf")
	flag.StringVar(&Conf.Listenner.RPC_addr, "rpc", "", "-rpc=127.0.0.1:8899")
	flag.StringVar(&Conf.Listenner.Listen_addr, "tcp_listen", "", "-tcp_listen=127.0.0.1:9001")
	flag.StringVar(&Conf.Listenner.WebSocket_addr, "web_listen", "", "-web_listen=127.0.0.1:9002")
	flag.BoolVar(&Conf.Listenner.Tls, "tls", false, "-tls=true")
	flag.StringVar(&etcd_addr, "etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")
}

func InitConf() error {
	// Get config from etcd.
}
