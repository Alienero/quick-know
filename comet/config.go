// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"strings"

	"github.com/Alienero/quick-know/comet/config"
	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/utils/json"
)

var Conf = config.Config{}

var (
	etcd_addr string
)

func init() {
	// flag.StringVar(&path, "path", "", "-path=comet.conf")
	flag.StringVar(&Conf.Listener.RPC_addr, "rpc", "", "-rpc=127.0.0.1:8899")
	flag.StringVar(&Conf.Listener.Listen_addr, "tcp_listen", "", "-tcp_listen=127.0.0.1:9001")
	flag.StringVar(&Conf.Listener.WebSocket_addr, "web_listen", "", "-web_listen=127.0.0.1:9002")
	// flag.BoolVar(&Conf.Listenner.Tls, "tls", false, "-tls=true")
	flag.StringVar(&etcd_addr, "etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")
}

func InitConf() error {
	// Get config from etcd.
	// Init etcd client.
	Conf.Etcd_addr = strings.Split(etcd_addr, ",")
	init_etcd()
	// Get the etcd config.
	if err := json.Getter(getEtcdConf, &Conf.Etcd); err != nil {
		return err
	}
	// Get listener's conf.
	if err := json.Getter(getListenConf, &Conf.Listener); err != nil {
		return err
	}
	// Get the redis config.
	// if err := json.Getter(getRedisConf, &Conf.Redis); err != nil {
	// 	return err
	// }
	// Get the Restrictiont config.
	if err := json.Getter(getRestrictiontConf, &Conf.Restriction); err != nil {
		return err
	}
	// Get the store config.
	strStore, err := getStoreConf()
	if err != nil {
		return err
	}
	if err := store.Init([]byte(strStore)); err != nil {
		return err
	}
	// Start etcd hb.
	return etcd_hb()
}
