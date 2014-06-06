package comet

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/Alienero/spp"

	"github.com/golang/glog"
)

func StartListen() error {
	var tempDelay time.Duration // how long to sleep on accept failure
	l, err := net.Listen("tcp", Conf.Listen_addr)
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
		c := newConn(rw)
		go c.serve()
	}
}

// Process connetion settings
type conn struct {
	// net's Connection
	rw net.Conn
	// Small pack Connection
	// packRW spp.Conn
}

func newConn(rw net.Conn) *conn {
	return &conn{
		rw: rw,
	}
}
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
	// TODO: get the offline msg
	// Init the ssp
	packRW := spp.NewConn(tcp)
	// pack, err := packRW.ReadPack()
	// if err != nil {
	// 	glog.Errorf("Recive login pack error:%v \n", err)
	// }
	var l listener
	if l, err = c.login(packRW); err != nil {
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
func (c *conn) login(rw *spp.Conn) (l listener, err error) {
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
	//TODO: DB Check
	switch req.Typ {
	case 0:
		l = newClient(rw)
	case 1:
		l = newCServer(rw)
	default:
		fmt.Errorf("No such pack type :%v", pack.Typ)
	}
	return
}

// Listen the clients' or controller server's request
type listener interface {
	listen_loop() error
}

// For call the serve
type handle func(c interface{}, pack *spp.Pack) error

func (h handle) serve(c interface{}, pack *spp.Pack) error {
	return h(c, pack)
}

// The Listen implement
type listen struct {
	// Write loop chan
	WriteChan chan *spp.Pack
	// Write error
	writeError error

	// Pack connection
	Rw *spp.Conn
	// Handles
	Handles map[int]handle
}

// func (l *listen) initListen(rw *spp.Conn) {
// 	l.rw = rw
// writeChan:
// 	make(chan *spp.Pack, Conf.WirteLoopChanNum)
// }
func (l *listen) listen_loop() (err error) {
	defer func() {
		// Close the res
		close(l.WriteChan)
	}()
	var pack *spp.Pack
	for {
		// Listen
		pack, err = l.Rw.ReadPack()
		if err != nil {
			// glog.Errorf("clientLoop read pack error:%v\n", err)
			break
		}
		f := l.Handles[pack.Typ]
		if f == nil {
			err = fmt.Errorf("No such pack type:%v", pack.Typ)
			break
		}
		// Call function f
		err = f.serve(l, pack)
		if err != nil {
			// glog.Errorf("clientLoop() f.serve() error:%v\n", err)
			break
		}
	}
	return
}

// Server write queue
func (l *listen) writePack(pack *spp.Pack) error {
	if l.writeError != nil {
		return l.writeError
	}
	l.WriteChan <- pack
	return nil
}
func (l *listen) writeLoop() {
loop:
	for {
		select {
		case pack := <-l.WriteChan:
			if pack == nil {
				break loop
			}
			err := l.Rw.WritePack(pack)
			if err != nil {
				// Tell listen error
				l.writeError = err
				break loop
			}
		}
	}
}
