// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"runtime"
	"time"

	"github.com/Alienero/quick-know/mqtt"
	myrpc "github.com/Alienero/quick-know/rpc"
	"github.com/Alienero/quick-know/store"

	"github.com/golang/glog"
)

type ConnListener struct {
	listener *net.TCPListener
}

func newConnListener(l *net.TCPListener) *ConnListener {
	return &ConnListener{l}
}

func (l *ConnListener) Accept() (c net.Conn, err error) {
	var tcp *net.TCPConn
	tcp, err = l.listener.AcceptTCP()
	if err != nil {
		return
	}
	// Tcp Setting。
	if err = tcp.SetKeepAlive(true); err != nil {
		glog.Errorf("conn.SetKeepAlive() error(%v)", err)
		return
	}
	return tcp, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *ConnListener) Close() error { return l.listener.Close() }

// Addr returns the listener's network address.
func (l *ConnListener) Addr() net.Addr { return l.listener.Addr() }

func startListen(typ int, addr string) {
	err := func() error {
		var (
			tempDelay time.Duration // how long to sleep on accept failure
			l         net.Listener
			err       error
		)
		l, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		tl := l.(*net.TCPListener)
		l = newConnListener(tl)
		if Conf.Tls {
			tlsConf := new(tls.Config)
			tlsConf.Certificates = make([]tls.Certificate, 1)
			tlsConf.Certificates[0], err = tls.X509KeyPair(Conf.Cert, Conf.Key)
			if err != nil {
				glog.Fatal(err)
			}
			l = tls.NewListener(l, tlsConf)
		}

		for {
			rw, e := l.Accept()
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
					time.Sleep(tempDelay)
					continue
				}
				glog.Errorf("http: Accept error: %v; retrying in %v", e, tempDelay)
				return e
			}

			c := newConn(rw, typ)
			go c.serve()
		}
	}()
	glog.Fatal(err)
}

// Process connetion settings
type conn struct {
	// net's Connection
	rw net.Conn
	r  *bufio.Reader
	w  *bufio.Writer
	// The conn's listen type
	typ int
}

func newConn(rw net.Conn, typ int) *conn {
	return &conn{
		rw:  rw,
		r:   bufio.NewReader(rw),
		w:   bufio.NewWriter(rw),
		typ: typ,
	}
}

// Do login check and response or push
func (c *conn) serve() {
	var err error
	defer func() {
		if err := recover(); err != nil {
			buff := make([]byte, 4096)
			runtime.Stack(buff, false)
			glog.Errorf("conn.serve() panic(%v)\n info:%s", err, string(buff))
		}
		c.rw.Close()

	}()

	var l listener
	if l, err = login(c.r, c.w, c.rw, c.typ); err != nil {
		glog.Errorf("Login error :%v\n", err)
		return
	}
	err = mqtt.WritePack(mqtt.GetConnAckPack(0), c.w)
	if err != nil {
		return
	}
	l.listen_loop()
}

func login(r *bufio.Reader, w *bufio.Writer, conn net.Conn, typ int) (l listener, err error) {
	var pack *mqtt.Pack
	pack, err = mqtt.ReadPack(r)
	if err != nil {
		glog.Error("Read login pack error")
		return
	}
	if pack.GetType() != mqtt.CONNECT {
		err = fmt.Errorf("Recive login pack's type error:%v \n", pack.GetType())
		return
	}
	info, ok := (pack.GetVariable()).(*mqtt.Connect)
	if !ok {
		err = errors.New("It's not a mqtt connection package.")
		return
	}
	id := info.GetUserName()
	psw := info.GetPassword()

	switch typ {
	case CLIENT:
		if !store.Manager.Client_login(*id, *psw) {
			err = fmt.Errorf("Client Authentication is not passed id:%v,psw:%v", *id, *psw)
			break
		}
		// Has been already logon
		// TODO
		var (
			ok bool
			s  string
		)
	re:
		ok, s, err = redis.IsExist(*id)
		if err != nil {
			return
		}
		if ok {
			var client *rpc.Client
			client, err = rpc.DialHTTP("tcp", s)
			if err != nil {
				return
			}
			reply := new(myrpc.Reply)
			err = client.Call("Comet_RPC.Relogin", *id, reply)
			if err != nil {
				return
			}
			if reply.IsRe {
				goto re
			}
			if !reply.IsOk {
				err = errors.New("Has been relogining")
				return
			}
		}

		c := newClient(r, w, conn, *id, info.GetKeepAlive())
		// Redis Append.
		if err = redis.Login(*id, Conf.RPC_addr); err != nil {
			return
		}
		Users.Set(*id, c)
		l = c
	case CSERVER:
		// TODO
	default:
		fmt.Errorf("No such pack type :%v", typ)
	}
	return
}

// Listen the clients' or controller server's request
type listener interface {
	listen_loop() error
}
