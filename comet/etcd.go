// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"time"

	define "github.com/Alienero/quick-know/config"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

var etcd_client *etcd.Client

func init_etcd() {
	etcd_client = etcd.NewClient(Conf.Etcd_addr)
}

func etcd_hb() error {
	flush_time := time.Duration(float64(Conf.Etcd_interval) / 1.5)
	// Connect the etcd.
	_, err := etcd_client.Set(Conf.Etcd_dir+"/"+Conf.RPC_addr, "0", Conf.Etcd_interval)
	if err != nil {
		return err
	}
	c_time := time.NewTicker(flush_time * time.Second)
	go func() {
		for {
			select {
			case <-c_time.C:
				// Flush the etcd node time.
				if _, err = etcd_client.Update(Conf.Etcd_dir+"/"+Conf.RPC_addr, strconv.Itoa(Users.Len()), Conf.Etcd_interval); err != nil {
					glog.Fatalf("Comet system will be closed ,err:%v\n", err)
				}
			}
		}
	}()
	return nil
}

func getStoreConf() (string, error) {
	resp, err := etcd_client.Get(define.Etcd_store, false, false)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func getRedisConf() (string, error) {
	resp, err := etcd_client.Get(define.Etcd_comet_redis, false, false)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func getRestrictiontConf() (string, error) {
	resp, err := etcd_client.Get(define.Etcd_comet_rest, false, false)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func getEtcdConf() (string, error) {
	resp, err := etcd_client.Get(define.Etcd_comet_etcd, false, false)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func getListenConf() (string, error) {
	resp, err := etcd_client.Get(define.Etcd_comet_listen, false, false)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// func setStoreConf(conf interface{}) error {
// 	data, err := json.Marshal(conf)
// 	if err != nil {
// 		_, err = etcd_client.Set("quick-know/store_conf", string(data), 0)
// 	}
// 	return err
// }

// func setRedisConf(conf interface{}) error {
// 	data, err := json.Marshal(conf)
// 	if err != nil {
// 		_, err = etcd_client.Set("quick-know/redis_conf", string(data), 0)
// 	}
// 	return err
// }

// func setRestrictiontConf(conf interface{}) error {
// 	data, err := json.Marshal(conf)
// 	if err != nil {
// 		_, err = etcd_client.Set("quick-know/comet_conf", string(data), 0)
// 	}
// 	return err
// }

// func setEtcdConf(conf interface{}) error {
// 	data, err := json.Marshal(conf)
// 	if err != nil {
// 		_, err = etcd_client.Set("quick-know/etcd_conf", string(data), 0)
// 	}
// 	return err
// }
