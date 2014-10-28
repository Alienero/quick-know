// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restful_tool

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	Url string
	ID  string
	Psw string
	c   = new(http.Client)

	pool = &sync.Pool{
		New: func() interface{} {
			return make(Urls, 0, 4)
		},
	}
)

type Urls []string

func NewUrls(path string) Urls {
	s := pool.Get().(Urls)
	s = append(s, Url)
	return append(s, path)
}

func (u *Urls) Add(s string) {
	*u = append(*u, s)
}

func (u Urls) String() string {
	return strings.Join(u, "")
}

func (u Urls) ReFresh() {
	pool.Put(u)
}

func AddPrivateMsg(to_id string, expired int64, msg []byte) error {
	us := NewUrls("/push/private")
	defer us.ReFresh()
	us.Add("?to_id=" + to_id)
	if expired > 0 {
		us.Add("&expired=" + strconv.FormatInt(expired, 10))
	}
	req, err := http.NewRequest("PUT", us.String(), bytes.NewReader(msg))
	if err != nil {
		return err
	}
	req.SetBasicAuth(ID, Psw)
	println(us.String())
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Add private msg fail.")
	}
	return nil
}

func AddUser(psw string) (string, error) {
	req, err := http.NewRequest("POST", "http://127.0.0.1:9901/push/add_user", bytes.NewReader([]byte(`{"psw":"`+psw+`"}`)))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(ID, Psw)
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	// Get the User ID
	type ID struct {
		Id string
	}
	id := new(ID)
	defer resp.Body.Close()
	err = getObject(id, resp.Body)
	return id.Id, err
}

func PushMsg2All(msg []byte, expired int64) error {
	type mt struct {
		Body    []byte
		Expired int64
	}
	m := &mt{msg, expired}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// Post to server
	req, err := http.NewRequest("POST", "http://127.0.0.1:9901/push/all", bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(res) != `{"status":"success"}` {
		err = errors.New("push fail")
	}
	return err
}

func getObject(v interface{}, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
