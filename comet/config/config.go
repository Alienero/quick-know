// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

type Config struct {
	Listener
	Restriction
	// Redis
	Etcd
}

type Listener struct {
	Listen_addr    string `json:"-"` // Client listener addr
	WebSocket_addr string `json:"-"`
	RPC_addr       string `json:"-"`
	Tls            bool   `json:"-"`
	Cert           []byte
	Key            []byte
}

type Restriction struct {
	WirteLoopChanNum int // Should > 1
	ReadPackLoop     int
	MaxCacheMsg      int
	ReadTimeout      int // Heart beat check (seconds)
	WriteTimeout     int
}

// type Redis struct {
// 	// Redis conf
// 	Network    string
// 	Address    string
// 	MaxIde     int
// 	IdeTimeout int // Second.
// }

type Etcd struct {
	// Etcd conf.
	Etcd_addr     []string `json:"-"`
	Etcd_interval uint64
	Etcd_dir      string
}
