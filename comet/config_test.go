package main

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	if err := InitConf(); err != nil {
		t.Error(err)
	}
	fmt.Println(Conf)
}
