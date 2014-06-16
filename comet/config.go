// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"encoding/json"
	"io/ioutil"
)

var Conf = &config{}

type config struct {
	Listen_addr string // Client listener addr

	WirteLoopChanNum int // Should > 1

	ReadPackLoop int

	MaxCacheMsg int

	ReadTimeout  int // Heart beat check (seconds)
	WriteTimeout int
}

func InitConf() error {
	data, err := ioutil.ReadFile("comet.conf")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, Conf)
	if err != nil {
		return err
	}
	return nil
}
