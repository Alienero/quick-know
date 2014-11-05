// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"testing"
)

func TestRedis(t *testing.T) {
	func() {
		// Redis init.
		Conf.Network = "tcp"
		Conf.IdeTimeout = 240
		Conf.MaxIde = 3
		Conf.Address = "127.0.0.1:6379"
		Conf.Listen_addr = "127.0.0.1:9900"
		Conf.RPC_addr = "127.0.0.1:8899"
	}()
	if err := redis_login("1234"); err != nil {
		t.Error(err)
	}
	if b, s, err := redis_isExist("1234"); err != nil {
		t.Error(err)
	} else {
		if !b {
			t.Error(errors.New("Not exist."))
		}
		println("ok:", s)
	}
	if err := redis_logout("1234"); err != nil {
		t.Error(err)
	}
	if b, _, err := redis_isExist("1234"); err != nil {
		t.Error(err)
	} else {
		if b {
			t.Error(errors.New("exist."))
		}
	}
}
