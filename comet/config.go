// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

var Conf = &config{}

type config struct {
	Listen_addr string // Client listener addr

	WirteLoopChanNum int // Should > 1

	ReadPackLoop int

	MaxCacheMsg int

	ReadTimeout  int // Heart beat check (seconds)
	WriteTimeout int

	// Redis conf
	Network    string
	Address    string
	MaxIde     int
	IdeTimeout int // Second.
}

func InitConf() error {
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
