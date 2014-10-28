// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restful_tool

import (
	"testing"
)

func TestAddPrivateMsg(t *testing.T) {
	Url = "http://127.0.0.1:9901"
	ID = "1234"
	Psw = "10086"
	if err := AddPrivateMsg("apq5y6w9stc4", 0, []byte("Private msg")); err != nil {
		t.Error(err)
	} else {
		println("Pass")
	}
}
