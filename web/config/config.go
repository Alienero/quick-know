// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

type Config struct {
	Listen_addr string

	// Etcd conf.
	Etcd_addr     []string
	Etcd_interval uint64
	Etcd_dir      string
	From_etcd     bool

	Balancer     string // CoreBanlancing or Addr or domain
	Comet_domain string
	Comet_port   string

	Comet_addr string
	// CoreBanlancing conf.
	Cbl_addr string
}
