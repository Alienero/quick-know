// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The package get a comet from the CoreBalancing
// (https://github.com/CoreTalk/CoreBanlancing).
// And rpc a msg to comet server.

package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/rpc"

	myrpc "github.com/Alienero/quick-know/rpc"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/store/define"
)

var http_clinet = new(http.Client)

func get_comet() (string, error) {
	if Conf.BalancerType != "CoreBanlancing" {
		return Conf.CometRPC_addr, nil
	} else {
		resp, err := http_clinet.Get(Conf.Cbl_addr + "/get_server")
		if err != nil {
			return "", err
		}
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}

func write_msg(msg *define.Msg) error {
	//  Check the user online.
	b, addr, err := redis.IsExist(msg.To_id)
	if !b {
		if err != nil {
			return err
		}
		// User is offline.
		return writeOfflineMsg(msg)
	}
	return writeByAddr(msg, addr)
}

func writeByGetComet(msg *define.Msg) error {
	addr, err := get_comet()
	if err != nil {
		return err
	}
	c, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return err
	}
	reply := new(myrpc.Reply)
	if err = c.Call("Comet_RPC.WriteMsg", msg, reply); err != nil {
		return err
	}
	if !reply.IsOk {
		err = errors.New("RPC:wirte msg fail")
	}
	return err
}

func writeByAddr(msg *define.Msg, addr string) error {
	c, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return err
	}
	reply := new(myrpc.Reply)
	if err = c.Call("Comet_RPC.WriteMsg", msg, reply); err != nil {
		return err
	}
	if !reply.IsOk {
		err = errors.New("RPC:wirte msg fail")
	}
	return err
}

func writeOfflineMsg(msg *define.Msg) error {
	msg.Typ = define.OFFLINE
	return store.Manager.InsertOfflineMsg(msg)
}
