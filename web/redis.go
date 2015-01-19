// Copyright © 2015 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	myredis "github.com/Alienero/quick-know/redis"
)

var redis *myredis.Redis

func init_redis(jconf string) (err error) {
	redis, err = myredis.NewRedis(jconf)
	return
}
