// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	// "encoding/json"
	"flag"
	"strings"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/quick-know/utils/json"
	"github.com/Alienero/quick-know/web/config"
)

var Conf = config.Config{}

var etcd_addr_temp string

func init() {
	flag.StringVar(&Conf.Listen_addr, "listen", "", "-listen=127.0.0.1:9002")
	// flag.BoolVar(&Conf.Tls, "tls", false, "-tls=true")
	flag.StringVar(&etcd_addr_temp, "etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")
}

func InitConf() error {
	// Read config frome etcd.
	// Init etcd client.
	Conf.Etcd_addr = strings.Split(etcd_addr_temp, ",")
	init_etcd()
	// Get store config.
	storeConf, err := getStore()
	if err != nil {
		return err
	}
	if err := store.Init([]byte(storeConf)); err != nil {
		return err
	}
	// Get web config.
	// Get listener's config.
	if err := json.Getter(getListener, &Conf.Listener); err != nil {
		return err
	}
	// Get balancer config.
	if err := json.Getter(getBalancer, &Conf.Balancer); err != nil {
		return err
	}
	// Get etcd config.
	if err := json.Getter(getEtcd, &Conf.Etcd); err != nil {
		return err
	}

	return etcd_hb()
}
