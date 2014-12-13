// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

type Config struct {
	Tls         bool   `jsong:"-"`
	Listen_addr string `json:"-"`

	BalancerType  string // CoreBanlancing or Addr or domain
	CometRPC_addr string
	// CoreBanlancing conf.
	Cbl_addr string

	// Etcd conf.
	Etcd_addr     []string `json:"-"`
	Etcd_interval uint64
	Etcd_dir      string
}
