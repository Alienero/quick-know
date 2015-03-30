// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provide a simple way to use etcd to
// election the leader of the cluster.
package election

import (
	"strings"
	"time"

	myetcd "github.com/Alienero/quick-know/etcd"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type Cluster struct {
	// etcd information.
	etcd_dir    string
	etcd_nodes  string
	etcd_leader string
	interval    time.Duration

	leader string
}

type Node struct {
	keeper *myetcd.Conn
	addr   string
}

// Listen and watch.
func (n *Node) Serve() {
	// Get the leader node. If leader not exist condition leader.
loop:
	stop, nw, run, lastVaule := n.watch()
	if nw {
		if err := n.keeper.Create(cluster.etcd_dir+cluster.etcd_leader, node.addr, cluster.interval); err != nil {
			goto loop
		}

		// Local host become leader.
		n.leaderKeeper()
		return
	}
	if run {
		// election.
		if stop != nil {
			stop <- true
		}
		if err := n.keeper.CompareAndSwap(cluster.etcd_dir+cluster.etcd_leader, lastVaule, node.addr, cluster.interval); err != nil {
			goto loop
		}
		n.leaderKeeper()
	}
}

func (n *Node) leaderKeeper() {
	n.keeper.IntervalUpdate(cluster.etcd_dir+cluster.etcd_leader, node.addr, cluster.interval)
}

// watch the leader (TODO: other nodes).
func (n *Node) watch() (stop chan bool, nw bool, goon bool, lastVaule string) {
	// Watch leader node.
	// Get the leader's addr.
	var (
		reciver chan *etcd.Response
		err     error
	)
	leader, err := n.keeper.Get(cluster.etcd_dir + cluster.etcd_leader)
	if err != nil {
		goon = true
		if e, ok := err.(*etcd.EtcdError); ok {
			if e.ErrorCode == 100 {
				nw = true
			}
		}
		return
	}
	if leader == n.addr {
		// Pass
		goon = false
		return
	}
	reciver, stop = n.keeper.Watch(cluster.etcd_dir + cluster.etcd_leader)
	// timer := time.NewTimer(d)
	for {
		select {
		case resp := <-reciver:
			if resp == nil {
				// Warning!
			}
			if resp.Action != "update" {
				goon = true
				return
			}
		}
	}
}

func (n *Node) intervalUpdate() chan bool {
	return n.keeper.IntervalUpdate(cluster.etcd_dir+cluster.etcd_nodes, node.addr, cluster.interval)
}

func (n *Node) GetCluster() {}

var (
	node    *Node // local host.
	cluster *Cluster
)

func Init(machines []string, etcd_dir, etcd_nodes, etcd_leader, addr string, interval time.Duration) {
	// Init var.
	cluster = new(Cluster)
	node = new(Node)
	node.addr = addr
	var err error
	node.keeper, err = myetcd.InitClient(machines, etcd_dir, etcd_nodes, node.addr, interval)
	if err != nil {
		panic(err)
	}
	cluster.etcd_dir = etcd_dir
	if !strings.HasSuffix(cluster.etcd_dir, "/") {
		cluster.etcd_dir += "/"
	}
	cluster.etcd_nodes = etcd_nodes
	cluster.etcd_leader = etcd_leader
	cluster.interval = interval

	node.intervalUpdate()
}
