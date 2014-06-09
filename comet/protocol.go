package comet

import (
	"encoding/json"

	"github.com/Alienero/quick-know/store"
)

// Socket protocol
const (
	CLIENT  = 0
	CSERVER = 1

	PUSH_INFO = 21
	PUSH_MSG  = 22

	// Client requst type
	OFFLINE    = 11
	ONLINE     = 12
	HEART_BEAT = 31

	LOGIN  = 101
	LONGON = 102
)

type loginRequst struct {
	Id  string
	Psw string
	// 1 is Control server , 0 is client
	Typ int

	// Subscribe string
}

func getLoginRequst(data []byte) (l *loginRequst, err error) {
	l = new(loginRequst)
	err = unMarshalJson(data, l)
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

func getMsg(msg *store.Msg) ([]byte, error) {
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
