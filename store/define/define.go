// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package define

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Msg struct {
	Msg_id int    // Msg ID
	Owner  string // Owner
	To_id  string
	Topic  string
	Body   []byte
	Typ    int // Online or Oflline msg

	Dup byte // mqtt dup

	Expired int64
}

type Msg_id struct {
	Id string
	M  *Msg
}

type User struct {
	Id    string
	Psw   string
	Owner string // Owner
}

type Ctrl struct {
	Id   string
	Auth string
}

type Sub struct {
	Id  string
	Own string

	Max int // Max of the group members
	Typ int
}

type Sub_map struct {
	Sub_id  string
	User_id string
}

var lock = new(sync.Mutex)
var key string

func SetKey(addr string) {
	h := sha1.New()
	io.WriteString(h, strconv.FormatInt(time.Now().UTC().UnixNano(), 36)+addr)
	key = strings.Replace(fmt.Sprintf("% x", h.Sum(nil)), " ", "", -1)
}

func Get_uuid() string {
	lock.Lock()
	defer lock.Unlock()
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 36)
}
