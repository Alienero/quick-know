// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provide a simple way to use etcd to
// election the leader of the cluster.
package election

type Cluster struct{}

func GetCluster() {}

type Node struct {
}

func NewNode() *Node {
	return &Node{}
}

var (
	node *Node
)
