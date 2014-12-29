// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	pool *redis.Pool
}

func NewRedis(jconf string) (*Redis, error) {
	conf := struct {
		// Redis conf
		Network    string
		Address    string
		MaxIde     int
		IdeTimeout int // Second.
	}{}
	if err := json.Unmarshal([]byte(jconf), &conf); err != nil {
		return nil, err
	}
	client := &Redis{
		pool: &redis.Pool{
			MaxIdle:     conf.MaxIde,
			IdleTimeout: time.Duration(conf.IdeTimeout) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(conf.Network, conf.Address)
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
		},
	}
	return client, nil
}

func (r *Redis) Login(id, value string) error {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Send("SETNX", id, value)
	conn.Send("LPUSH", value, id)
	err := conn.Flush()
	if err != nil {
		return err
	}
	reply, err := conn.Receive()
	if err != nil {
		return err
	}
	if i, _ := redis.Int(reply, err); i != 1 {
		return errors.New("id not exist.")
	}
	_, err = conn.Receive()
	return err
}

func (r *Redis) Logout(id, value string) error {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Send("DEL", id)
	conn.Send("LREM", value, 1, id)
	err := conn.Flush()
	if err != nil {
		return err
	}
	_, err = conn.Receive()
	return err
}

func (r *Redis) IsExist(id string) (bool, string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", id)
	if err != nil {
		return false, "", err
	}
	s, _ := redis.String(reply, nil)
	if s != "" {
		// exist.
		return true, s, nil
	}
	return false, "", nil
}
