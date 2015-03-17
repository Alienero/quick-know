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
	"time"

	"github.com/nu7hatch/gouuid"
)

const (
	// Client requst type
	OFFLINE = 41
	ONLINE  = 42
)

type Msg struct {
	Id     string
	Msg_id int    // Msg ID
	Owner  string // Owner
	To_id  string
	Topic  string
	Body   []byte
	Typ    int // Online or Oflline msg

	Dup byte // mqtt dup

	Expired int64

	// sub msgs.
	IsSub bool
}

// type Msg_id struct {
// 	// Id    string
// 	IsSub bool
// 	M     *Msg
// }

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

type SubMsgs struct {
	Id    string
	Count int64
	Body  []byte
}

var key string

func SetKey(addr string) {
	h := sha1.New()
	io.WriteString(h, strconv.FormatInt(time.Now().UTC().UnixNano(), 36)+addr)
	key = strings.Replace(fmt.Sprintf("% x", h.Sum(nil)), " ", "", -1)
}

func Get_uuid() string {
	uu, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uu.String()
}
