// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"

	// Import config define
	comet "github.com/Alienero/quick-know/comet/config"
	define "github.com/Alienero/quick-know/config"
	store "github.com/Alienero/quick-know/store/define"
	web "github.com/Alienero/quick-know/web/config"

	"github.com/coreos/go-etcd/etcd"
)

var (
	path     = flag.String("path", "", "-path=qk.conf")
	etcd_tmp = flag.String("etcd", "", "-etcd=http://127.0.0.1:4001,http://127.0.0.1:4002,http://127.0.0.1:4003")

	logger = log.New(os.Stdout, "qk_conf", log.Ltime|log.Lshortfile|log.LstdFlags)

	Conf = config{}

	etcdClient = new(etcd.Client)
)

func init() {
	flag.Parse()
	// Init etcd.
	etcdClient = etcd.NewClient(strings.Split(*etcd_tmp, ","))
}

type config struct {
	Comet struct {
		comet.Etcd
		comet.Redis
		comet.Restriction
	}
	Web struct {
		web.Balancer
		web.Etcd
	}
	Store store.DBConfig
}

func main() {
	// Read config.
	if err := readFileInto(*path); err != nil {
		logger.Panic(err)
	}
	// Share config.
	logger.Println("Set Comet's config...")
	logger.Println("->Do comet.Etcd")
	if err := setNode(define.Etcd_comet_etcd, &Conf.Comet.Etcd); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("->Do comet.Redis")
	if err := setNode(define.Etcd_comet_redis, &Conf.Comet.Redis); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("->Do comet.Restriction")
	if err := setNode(define.Etcd_comet_rest, &Conf.Comet.Restriction); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("Do web.Balancer")
	if err := setNode(define.Etcd_web_balancer, &Conf.Web.Balancer); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("Do web.Etcd")
	if err := setNode(define.Etcd_web_etcd, &Conf.Web.Etcd); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("Do store.")
	if err := setNode(define.Etcd_store, &Conf.Store); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")

	logger.Println("Shared!")
}

func readFileInto(path string) error {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	return json.Unmarshal(data, &Conf)
}

func setNode(node string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = etcdClient.Set(node, string(data), 0)
	return err
}
