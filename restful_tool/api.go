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

func (u *Urls) Add(s ...string) {
	*u = append(*u, s...)
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
	us.Add("?to_id=", to_id)
	if expired > 0 {
		us.Add("&expired=", strconv.FormatInt(expired, 10))
	}
	req, err := http.NewRequest("PUT", us.String(), bytes.NewReader(msg))
	if err != nil {
		return err
	}
	req.SetBasicAuth(ID, Psw)
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
	us := NewUrls("/push/add_user")
	defer us.ReFresh()
	us.Add("?psw=", psw)
	req, err := http.NewRequest("PUT", us.String(), nil)
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
		Id string `json:"id"`
	}
	id := new(ID)
	defer resp.Body.Close()
	err = getObject(id, resp.Body)
	return id.Id, err
}

func DelUser(id string) error {
	us := NewUrls("/push/del_user")
	defer us.ReFresh()
	us.Add("?id=", id)
	req, err := http.NewRequest("DELETE", us.String(), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(ID, Psw)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("del user fail.")
	}
	return nil
}

func AddSub(max int, typ int) (string, error) {
	us := NewUrls("/push/add_sub")
	defer us.ReFresh()
	if max > 0 {
		us.Add("&max=", strconv.FormatInt(max, 10))
	}
	if typ > 0 {
		us.Add("&typ=", strconv.FormatInt(typ, 10))
	}
	req, err := http.NewRequest("PUT", us.String(), nil)
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
		Sub_id string `json:"sub_id"`
	}
	id := new(ID)
	defer resp.Body.Close()
	err = getObject(id, resp.Body)
	return id.Sub_id, err
}

// TODO:

func getObject(v interface{}, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
