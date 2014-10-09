package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Alienero/quick-know/store/define"
)

func TestPrivateMsg(t *testing.T) {
	u := &define.Msg{Body: []byte("这是离线消息1"), To_id: "apq5y6w9stc4"}
	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:9901/push/private", bytes.NewReader(data))
	if err != nil {
		t.Error(err)
		return
	}
	c := new(http.Client)
	str := base64.StdEncoding.EncodeToString([]byte("1234:10086"))
	// req.SetBasicAuth("username", "password")

	req.Header.Add("Authorization", " Basic "+str)
	resp, err := c.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	println(string(s))
}
