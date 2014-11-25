// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package define

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

var Config = new(config)

type config struct {
	DBType string

	UserAddr string // DB addr
	UserName string // Collection name
	Clients  string
	Ctrls    string
	Salt     string

	MsgAddr     string
	MsgName     string
	OfflineName string

	SubName  string
	SubsName string

	OfflineMsgs int // The max of the offline msgs
}

func InitConfig() error {
	buf := new(bytes.Buffer)

	f, err := os.Open("store.conf")
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
	return json.Unmarshal(buf.Bytes(), Config)
}
