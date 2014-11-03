// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"github.com/garyburd/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxIdle:     Conf.MaxIde,
	IdleTimeout: Conf.IdeTimeout * time.Second,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial(Conf.Network, Conf.Address)
		if err != nil {
			return nil, err
		}
		// if _, err := c.Do("AUTH", password); err != nil {
		// 	c.Close()
		// 	return nil, err
		// }
		return c, err
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

func redis_login(id string) error {
	conn := redisPool.Get()
	defer conn.Close()
	conn.Send("SET", id, Conf.Listen_addr)
	conn.Send("LPUSH", Conf.Listen_addr, id)
	err := conn.Flush()
	if err != nil {
		return err
	}
	_, err = conn.Receive()
	return err
}

func redis_logout(id string) error {
	conn := redisPool.Get()
	defer conn.Close()
	conn.Send("DEL", id)
	conn.Send("LREM", Conf.Listen_addr, 1, id)
	err := conn.Flush()
	if err != nil {
		return err
	}
	_, err = conn.Receive()
	return err
}
