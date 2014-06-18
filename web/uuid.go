// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock = new(sync.Mutex)

func get_uuid() string {
	h := md5.New()
	lock.Lock()
	io.WriteString(h, strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	lock.Unlock()
	return strings.Replace(fmt.Sprintf("% x", h.Sum(nil)), " ", "", -1)
}
