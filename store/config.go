// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"encoding/json"
	"io/ioutil"
)

var Config = new(config)

type config struct {
	UserAddr string // DB addr
	UserName string // Collection name
	Clients  string
	Ctrls    string

	MsgAddr     string
	MsgName     string
	OfflineName string

	SubName  string
	SubsName string

	OfflineMsgs int // The max of the offline msgs
}

func initConfig() error {
	data, err := ioutil.ReadFile("store.conf")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}
