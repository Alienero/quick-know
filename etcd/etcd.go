// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package etcd

import (
	"strings"
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

func InitClient(machines []string, dir, key, value string, ttl time.Duration) (conn *Conn, err error) {
	conn = Connet(machines)
	if strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	return conn, conn.Set(dir+key, value, ttl)
}

func (c *Conn) Set(key, value string, ttl time.Duration) error {
	_, err := c.Client.Set(key, value, uint64(ttl))
	return err
}

func (c *Conn) Get(key string) (string, error) {
	resp, err := c.Client.Get(key, false, false)
	return resp.Node.Value, err
}

func (c *Conn) GetAll(key string) (etcd.Nodes, error) {
	resp, err := c.Client.Get(key, false, true)
	return resp.Node.Nodes, err
}

func (c *Conn) IntervalUpdate(key, value string, interval time.Duration) (stop chan bool) {
	timer := time.NewTicker(interval / 2)
	stop = make(chan bool)
	go func() {
		for {
			select {
			case <-timer.C:
				// Heart beat and update the etcd node's infomation.
				if _, err := c.Client.Update(key, value, uint64(interval)); err != nil {
					glog.Fatalf("Comet system will be closed ,err:%v\n", err)
				}
			case <-stop:
				timer.Stop()
			}
		}
	}()
	return stop
}

func (c *Conn) watch(key string, recursive bool, receiver chan *etcd.Response, stop chan bool) error {
	_, err := c.Client.Watch(key, 0, recursive, receiver, stop)
	return err
}

func (c *Conn) WatchByChan(key string, receiver chan *etcd.Response) (stop chan bool, err error) {
	stop = make(chan bool)
	err = c.watch(key, false, receiver, stop)
	return
}

func (c *Conn) WatchAll(key string, receiver chan *etcd.Response) (stop chan bool, err error) {
	stop = make(chan bool)
	err = c.watch(key, true, receiver, stop)
	return
}

func (c *Conn) Watch(key string) (receiver chan *etcd.Response, stop chan bool, err error) {
	receiver = make(chan *etcd.Response, 10)
	stop = make(chan bool)
	err = c.watch(key, false, receiver, stop)
	return
}

func (c *Conn) WatchByBuff(key string, buffer int) (receiver chan *etcd.Response, stop chan bool, err error) {
	receiver = make(chan *etcd.Response, buffer)
	stop = make(chan bool)
	err = c.watch(key, false, receiver, stop)
	return
}

func (c *Conn) CompareAndSwap(key, old, nw string, ttl time.Duration) error {
	_, err := c.Client.CompareAndSwap(key, nw, uint64(ttl), old, 0)
	return err
}
