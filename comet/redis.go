// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

var redisPool = &redis.Pool{
	MaxIdle:     Conf.MaxIde,
	IdleTimeout: time.Duration(Conf.IdeTimeout) * time.Second,
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
	conn.Send("SET", id, Conf.RPC_addr)
	conn.Send("LPUSH", Conf.RPC_addr, id)
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
	conn.Send("LREM", Conf.RPC_addr, 1, id)
	err := conn.Flush()
	if err != nil {
		return err
	}
	_, err = conn.Receive()
	return err
}

func redis_isExist(id string) (bool, string, error) {
	conn := redisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", id)
	if err != nil {
		return false, "", err
	}
	if reply == nil {
		// not exist.
		return false, "", nil
	}
	if s, ok := reply.(string); ok && s != "" {
		// exist.
		return true, s, nil
	}
	return false, "", errors.New("redis get reply not type string.")
}
