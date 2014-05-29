package comet

import (
	"encoding/json"
	"io/ioutil"
)

var Conf = &config{}

type config struct {
	Listen_addr string
}

func InitConf() error {
	data, err := ioutil.ReadFile("comet.conf")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, Conf)
	if err != nil {
		return err
	}
}
