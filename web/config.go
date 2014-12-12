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

	"github.com/Alienero/quick-know/web/config"
)

var Conf = new(config.Config)

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
