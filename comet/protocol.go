package comet

import (
	"encoding/json"

	"github.com/Alienero/spp"
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
type loginResponse struct {
	ID   string
	Addr string

	State bool
	Info  string
}

func getLoginResponse(id, addr string, state bool, info string) ([]byte, error) {
	resp := &loginResponse{id, addr, state, info}
	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return body, nil
}
