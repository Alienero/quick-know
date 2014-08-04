// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/Alienero/quick-know/store"
	// "github.com/Alienero/spp"
	"github.com/Alienero/quick-know/mqtt"

	"github.com/golang/glog"
)

func startListen(typ int, addr string) error {
	var tempDelay time.Duration // how long to sleep on accept failure
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
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
	tcp := c.rw.(*net.TCPConn)
	if err = tcp.SetKeepAlive(true); err != nil {
		glog.Errorf("conn.SetKeepAlive() error(%v)\n", err)
		return
	}
	if err = tcp.SetReadDeadline(time.Second * Conf.ReadTimeout); err != nil {
		glog.Errorf("conn.SetReadDeadline() error(%v)\n", err)
		return
	}
	if err = tcp.SetWriteDeadline(time.Second * Conf.WriteTimeout); err != nil {
		glog.Errorf("conn.SetWriteDeadline() error(%v)\n", err)
		return
	}

	// packRW := spp.NewConn(tcp)
	// packRW.SetWriteDeadline(time.Duration(Conf.WriteTimeout) * time.Second)
	// packRW.SetReadDeadline(time.Duration(Conf.ReadTimeout) * time.Second)
	var l listener
	if l, err = login(packRW, c.typ); err != nil {
		glog.Errorf("Login error :%v\n", err)
		return
	}
	body, err := getLoginResponse("1", "127.0.0.1", true, "")
	if err != nil {
		return
	}
	pack, _ := packRW.SetDefaultPack(LOGIN, body)
	err = packRW.WritePack(pack)
	if err != nil {
		return
	}
	l.listen_loop()
}

func login(r *bufio.Reader, typ int) (l listener, err error) {
	var pack *mqtt.Pack
	pack, err = mqtt.ReadPack(r)
	if err != nil {
		glog.Error("Read login pack error")
		return
	}
	if pack.GetType() != mqtt.CONNECT {
		err = fmt.Errorf("Recive login pack's type error:%v \n", pack.Typ)
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
		if !store.Client_login(*id, *psw) {
			err = fmt.Errorf("Client Authentication is not passed id:%v,psw:%v", *id, *psw)
			break
		}
		// Has been already logon
		if tc := Users.Get(*id); tc != nil {
			tc.lock.Lock()
			if !tc.isLetClose {
				tc.lock.Unlock()
				select {
				case tc.CloseChan <- 1:
					tc.lock.Lock()
					tc.isLetClose = true
					tc.lock.Unlock()
					<-tc.CloseChan
				case <-time.After(3 * time.Second):
					if tc := Users.Get(req.Id); tc != nil {
						return nil, errors.New("Close the logon user timeout")
					}
				}
			} else {
				return nil, errors.New("Has been relogining")
			}

		}
		c := newClient(rw, req.Id)
		Users.Set(req.Id, c)
		l = c
	case CSERVER:
		// TODO : Base64
		if !store.Ctrl_login_alive() {
			err = fmt.Errorf("Client Authentication is not passed id:%v,psw:%v", *id, *psw)
			break
		}
		// TODO limit ctrl server users
		cs := newCServer(rw, req.Id)
		ctrls.Set(*id, cs)
		l = cs
	default:
		fmt.Errorf("No such pack type :%v", typ)
	}
	return
}

// Listen the clients' or controller server's request
type listener interface {
	listen_loop() error
}
