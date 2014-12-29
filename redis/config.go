// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

type RedisConf struct {
	// Redis conf
	Network    string
	Address    string
	MaxIde     int
	IdeTimeout int // Second.
}
