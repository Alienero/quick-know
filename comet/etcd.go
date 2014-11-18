// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

func Init_etcd() error {
	flush_time := time.Duration(float64(Conf.Etcd_interval) / 1.5)
	// Connect the etcd.
	client := etcd.NewClient(Conf.Etcd_addr)
	_, err := client.Set(Conf.Etcd_dir+"/"+Conf.RPC_addr, "0", Conf.Etcd_interval)
	if err != nil {
		return err
	}
	c_time := time.NewTicker(flush_time * time.Second)
	go func() {
		for {
			select {
			case <-c_time.C:
				// Flush the etcd node time.
				if _, err = client.Update(Conf.Etcd_dir+"/"+Conf.RPC_addr, strconv.Itoa(Users.Len()), Conf.Etcd_interval); err != nil {
					glog.Fatalf("Comet system will be closed ,err:%v\n", err)
				}
			}
		}
	}()
	return nil
}
