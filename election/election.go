// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provide a simple way to use etcd to
// election the leader of the cluster.
package election

import (
	"time"
)

type Cluster struct {
	// etcd information.
	etcd_dir    string
	etcd_nodes  string
	etcd_leader string
	interval    time.Duration
}

type Node struct {
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) GetCluster() {}

var (
	node    *Node
	cluster *Cluster
)
