package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/Alienero/quick-know/comet"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"
)

// func TestAddUser(t *testing.T) {
// 	u := &store.User{Psw: "1024"}
// 	data, err := json.Marshal(u)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	req, err := http.NewRequest("POST", "http://127.0.0.1:9901/push/add_user", bytes.NewReader(data))
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	c := new(http.Client)
// 	str := base64.StdEncoding.EncodeToString([]byte("test001"))
// 	// req.SetBasicAuth("username", "password")

// 	req.Header.Add("Authorization", " Basic "+str)
// 	resp, err := c.Do(req)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	s, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		println(err.Error())
// 	}
// 	println(string(s))

// 	// Add private msg

// 	// req,err = http.NewRequest("POST", "http://127.0.0.1:9901", body)
// }

func addMsg(t *testing.T) {
	u := &store.Msg{Body: []byte("hello push server"), To_id: "29d2b76f47e4f2e36e732a53c74e2731"}
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
	str := base64.StdEncoding.EncodeToString([]byte("615582195:1"))
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
func loginAndGetMsg(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9900")
	if err != nil {
		t.Error(err)
		return
	}

	c := spp.NewConn((conn.(*net.TCPConn)))
	// Login
	id := "29d2b76f47e4f2e36e732a53c74e2731"
	psw := "1024"
	owner := "615582195"
	buff := new(bytes.Buffer)
	buff.WriteByte(byte(len(id)))
	buff.Write([]byte(id))

	buff.WriteByte(byte(len(psw)))
	buff.Write([]byte(psw))

	buff.WriteByte(byte(len(owner)))
	buff.Write([]byte(owner))

	pack, err := c.SetDefaultPack(comet.LOGIN, buff.Bytes())
	if err != nil {
		t.Error(err)
		return
	}
	c.WritePack(pack) // login
	// Recive the response pack
	_, err = c.ReadPack()
	if err != nil {
		t.Error(err)
		return
	}

	pack, err = c.ReadPack()
	if err != nil {
		t.Error(err)
		return
	}
	println(string(pack.Body))

	// Response the msg
	msg := new(store.Msg)
	if err := json.Unmarshal(pack.Body, msg); err != nil {
		t.Error(err)
		return
	}
	pack, err = c.SetDefaultPack(41, []byte(msg.Msg_id))
	if err != nil {
		t.Error(err)
		return
	}
	if err = c.WritePack(pack); err != nil {
		t.Error(err)
	}
}

func TestPrivateMsg(t *testing.T) {
	addMsg(t)
	loginAndGetMsg(t)
}
