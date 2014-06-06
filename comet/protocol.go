package comet

import (
	"encoding/json"

	// "github.com/Alienero/spp"
)

// Socket protocol
const (
	PUSH_INFO  = 21
	HEART_BEAT = 31

	LOGIN  = 101
	LONGON = 102
)

type loginRequst struct {
	UserName string
	Psw      string

	Subscribe string
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

func marshalJson(v interface{}) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return body, nil
}
