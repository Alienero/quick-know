// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"testing"
)

func TestRedis(t *testing.T) {
	if err := redis_login("1234"); err != nil {
		t.Error(err)
	}
	if err := redis_logout("1234"); err != nil {
		t.Error(err)
	}
}
