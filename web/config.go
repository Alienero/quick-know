// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

var Conf = new(config)

type config struct {
	Listen_addr string

	// Etcd conf.
	Etcd_addr     []string
	Etcd_interval uint64
	Etcd_dir      string
	From_etcd     bool

	Balancer     string // CoreBanlancing or Addr or domain
	Comet_domain string
	Comet_port   string

	Comet_addr string
	// CoreBanlancing conf.
	Cbl_addr string
}

func InitConf() error {
	buf := new(bytes.Buffer)
	f, err := os.Open("web.conf")
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
