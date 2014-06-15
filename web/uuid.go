// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"strconv"
	"sync"
	"time"
)

var lock = new(sync.Mutex)

func get_uuid() string {
	lock.Lock()
	defer lock.Unlock()
	return strconv.FormatInt(time.Now().Unix(), 10)
}
