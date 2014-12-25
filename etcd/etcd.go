// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etcd

import (
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type Conn struct {
	Client *etcd.Client
}

func Connet(machines []string) *Conn {
	return &Conn{
		Client: etcd.NewClient(machines),
	}
}

func (c *Conn) Set(key, value string, ttl time.Duration) error {
	_, err := etcd_client.Set(key, value, ttl)
	return err
}

func (c *Conn) Get() {}

func (c *Conn) IntervalUpdate(key, value string, interval time.Duration) (stop chan bool) {
	timer := time.NewTicker(interval / 2)
	stop = make(chan bool)
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

func (c *Conn) Watch() {}

func (c *Conn) CompareAndSwap() {}
