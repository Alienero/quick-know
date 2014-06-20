// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	// "bytes"
	"encoding/json"

	"github.com/Alienero/quick-know/store"
)

// Socket protocol
const (
	CLIENT  = 0
	CSERVER = 1

	PUSH_INFO    = 21
	PUSH_MSG     = 22
	PUSH_OFFLINE = 23

	// Client requst type
	OFFLINE = 41
	ONLINE  = 42

	HEART_BEAT = 31

	LOGIN  = 101
	LONGON = 102
)

type loginRequst struct {
	Id  string
	Psw string
	// 1 is Control server , 0 is client
	Typ   int
	Owner string

	// Subscribe string
}

func getLoginRequst(data []byte) (l *loginRequst, err error) {
	l = new(loginRequst)
	buf := newBuffer(data)
	idLen, err := buf.readByte()
	if err != nil {
		return
	}
	l.Id, err = buf.readString(int(idLen))
	if err != nil {
		return
	}

	pswLen, err := buf.readByte()
	if err != nil {
		return
	}
	l.Psw, err = buf.readString(int(pswLen))
	if err != nil {
		return
	}

	owLen, err := buf.readByte()
	if err != nil {
		return
	}
	l.Owner, err = buf.readString(int(owLen))
	return
}

func getLoginResponse(id, addr string, status bool, info string) ([]byte, error) {
	type loginResponse struct {
		ID   string
		Addr string

		Status bool
		Info   string
	}
	resp := &loginResponse{id, addr, status, info}
	return marshalJson(resp)
}

type beat_heart struct{}
type beat_heartResp struct {
	Status bool
}

func getbeat_heartResp(status bool) ([]byte, error) {
	resp := beat_heartResp{status}
	return marshalJson(resp)
}

// Get the msg's data
type user_msg struct {
	Body   []byte
	Msg_id string
	Typ    int
}

func getMsg(msg *store.Msg) ([]byte, error) {
	return marshalJson(getUserMsg(msg))
}
func getUserMsg(msg *store.Msg) *user_msg {
	return &user_msg{
		Body:   msg.Body,
		Msg_id: msg.Msg_id,
		Typ:    msg.Typ,
	}
}

type offineMsg struct {
	Ms []*user_msg
}

func getOffineMsg(msg *offineMsg) ([]byte, error) {
	return marshalJson(msg)
}
func marshalJson(v interface{}) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func unMarshalJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
