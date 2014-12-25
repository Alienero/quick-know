// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package election

import (
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

var (
	etcd_client *etcd.Client
)

func setEtcdClient(nodes []string) {
	if len(nodes) == 0 {
		panic("none machines.")
	}
	etcd_client = etcd.NewClient(nodes)
	// init the client.
	if _, err := etcd_client.Set(cluster.etcd_dir+cluster.etcd_nodes, "active", cluster.interval); err != nil {
		panic(err)
	}
}

func etcd_hb(dir, key string, interval time.Duration) chan byte {
	timer := time.NewTicker(interval / 2)
	stop := make(chan byte)
	go func() {
		for {
			select {
			case <-timer.C:
				// Heart beat and update the etcd node's infomation.
				if _, err = etcd_client.Update(dir, key, interval); err != nil {
					glog.Fatalf("Comet system will be closed ,err:%v\n", err)
				}
			case <-stop:
				timer.Stop()
			}
		}
	}()
	return stop
}

func cas() {
	// etcd_client.
}
