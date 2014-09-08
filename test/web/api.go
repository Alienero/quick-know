package web

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ID  string
	Psw string
	c   = new(http.Client)
)

type m map[interface{}]interface{}

func AddUser(psw string) (string, error) {
	req, err := http.NewRequest("POST", "http:127.0.0.1:9901", bytes.NewReader([]byte(`{"psw":"`+psw+`"}`)))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(ID, Psw)
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	// Get the User ID
	id := m{"id": ""}
	err = getObject(id, resp.Body)
	return id["id"].(string), err
}

func getObject(v interface{}, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
