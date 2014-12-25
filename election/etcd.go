// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package election

import (
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var (
	etcd_client *etcd.Client
)

func setEtcdClient(nodes []string) {
	if len(nodes) == 0 {
		panic("none machines.")
	}
	etcd_client = etcd.NewClient(nodes)
}

func etcd_hb(interval time.Duration) (err error) {
	timer := time.NewTicker(interval / 2)
	go func() {

		for {
			select {
			case <-timer.C:
				// Heart beat and update the etcd node's infomation.
				etcd
			}
		}
	}()
}
