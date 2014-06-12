package comet

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"

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
	// The conn's listen type
	typ int
}

func newConn(rw net.Conn, typ int) *conn {
	return &conn{
		rw:  rw,
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
	packRW := spp.NewConn(tcp)
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

func login(rw *spp.Conn, typ int) (l listener, err error) {
	var pack *spp.Pack
	pack, err = rw.ReadPack()
	if err != nil {
		return
	}
	if pack.Typ != LOGIN {
		err = fmt.Errorf("Recive login pack's type error:%v \n", pack.Typ)
		return
	}
	// Marshal Json
	var req *loginRequst
	req, err = getLoginRequst(pack.Body)
	if err != nil {
		return
	}
	if req.Typ != typ {
		return nil, fmt.Errorf("request type error:%v", req.Typ)
	}

	switch req.Typ {
	case CLIENT:
		if !store.Client_login(req.Id, req.Psw, req.Owner) {
			err = fmt.Errorf("Client Authentication is not passed id:%v,psw:%v,owner:%v", req.Id, req.Psw, req.Owner)
			break
		}
		// Has been already logon
		if tc := uesers.get(req.Id); tc != nil {
			tc.CloseChan <- 1
			<-tc.CloseChan
		}
		c := newClient(rw, req.Id)
		uesers.set(req.Id, c)
		l = c
	case CSERVER:
		if !store.Ctrl_login(req.Id, req.Psw) {
			err = fmt.Errorf("Client Authentication is not passed id:%v,psw:%v", req.Id, req.Psw)
			break
		}
		// TODO limit ctrl server users
		cs := newCServer(rw, req.Id)
		ctrls.set(req.Id, cs)
		l = cs
	default:
		fmt.Errorf("No such pack type :%v", pack.Typ)
	}
	return
}

// Listen the clients' or controller server's request
type listener interface {
	listen_loop() error
}
