// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// Import config define
	comet "github.com/Alienero/quick-know/comet/config"
	define "github.com/Alienero/quick-know/config"
	redis "github.com/Alienero/quick-know/redis"
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
	OnCometTls bool
	OnWebTls   bool

	CometCert string
	CometKey  string
	WebCert   string
	WebKey    string

	Comet struct {
		comet.Etcd
		// comet.Redis
		comet.Restriction
		comet.Listener
	}
	Web struct {
		web.Balancer
		web.Etcd
		comet.Listener
	}
	Store store.DBConfig
	Redis redis.RedisConf
}

func main() {
	// Read config.
	if err := readFileInto(*path); err != nil {
		logger.Panic(err)
	}
	// Check tls.
	if Conf.OnCometTls {
		fileToStruct(Conf.CometCert, &Conf.Comet.Listener.Cert)
		fileToStruct(Conf.CometKey, &Conf.Comet.Listener.Key)
	}
	if Conf.OnWebTls {
		fileToStruct(Conf.WebCert, &Conf.Web.Listener.Cert)
		fileToStruct(Conf.WebKey, &Conf.Web.Listener.Key)
	}
	// Share config.
	logger.Println("Set Comet's config...")
	logger.Println("->Do comet.Etcd")
	if err := setNode(define.Etcd_comet_etcd, &Conf.Comet.Etcd); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("->Do comet.Listener")
	if err := setNode(define.Etcd_comet_listen, &Conf.Comet.Listener); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	// logger.Println("->Do comet.Redis")
	// if err := setNode(define.Etcd_comet_redis, &Conf.Comet.Redis); err != nil {
	// 	logger.Fatal(err)
	// }
	// logger.Println("Done.")
	logger.Println("->Do comet.Restriction")
	if err := setNode(define.Etcd_comet_rest, &Conf.Comet.Restriction); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")

	logger.Println("Set Web's config...")
	logger.Println("->Do web.Listener")
	if err := setNode(define.Etcd_web_listen, &Conf.Web.Listener); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("->Do web.Balancer")
	if err := setNode(define.Etcd_web_balancer, &Conf.Web.Balancer); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("->Do web.Etcd")
	if err := setNode(define.Etcd_web_etcd, &Conf.Web.Etcd); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("Set Store's config...")
	logger.Println("->Do store")
	if err := setNode(define.Etcd_store, &Conf.Store); err != nil {
		logger.Fatal(err)
	}
	logger.Println("Done.")
	logger.Println("Set Redis's config...")
	logger.Println("->Do redis")
	if err := setNode(define.Etcd_redis, &Conf.Redis); err != nil {
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
	switch v.(type) {
	case string:
		// String direct insert etcd.
		s := v.(string)
		_, err := etcdClient.Set(node, s, 0)
		return err
	case []byte:
		// Base64 encode.
		src := v.([]byte)
		_, err := etcdClient.Set(node, base64.StdEncoding.EncodeToString(src), 0)
		return err
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = etcdClient.Set(node, string(data), 0)
		return err
	}
}

// If read the file has an error,it will throws a panic.
func fileToStruct(path string, ptr *[]byte) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Panic(err)
	}
	*ptr = data
}
