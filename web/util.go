package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func readAdnGet(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
