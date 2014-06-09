package comet

import (
	"fmt"
	"net"
	"runtime"
	"time"

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
	//TODO: DB Check
	switch req.Typ {
	case CLIENT:
		l = newClient(rw, req.Id)
	case CSERVER:
		l = newCServer(rw, req.Id)
	default:
		fmt.Errorf("No such pack type :%v", pack.Typ)
	}
	return
}

// Listen the clients' or controller server's request
type listener interface {
	listen_loop() error
}

// Tcp write queue
type PackQueue struct {
	// The last error in the tcp connection
	writeError error
	// Notice read the error
	errorChan chan error

	writeChan chan *spp.Pack
	readChan  chan *packAndErr
	// Pack connection
	rw *spp.Conn
}
type packAndErr struct {
	pack *spp.Pack
	err  error
}

func NewPackQueue(rw *spp.Conn) *PackQueue {
	return &PackQueue{
		rw:        rw,
		writeChan: make(chan *spp.Pack, Conf.WirteLoopChanNum),
		readChan:  make(chan *packAndErr, 1),
		errorChan: make(chan error, 1),
	}
}
func (queue *PackQueue) writeLoop() {
	// defer recover()
	var err error
loop:
	for {
		select {
		case pack := <-queue.writeChan:
			if pack == nil {
				break loop
			}
			err = queue.rw.WritePack(pack)
			if err != nil {
				// Tell listen error
				queue.writeError = err
				break loop
			}
		}
	}
	// Notice the read
	if err != nil {
		queue.errorChan <- err
	}
}

// Server write queue
func (queue *PackQueue) WritePack(pack *spp.Pack) error {
	if queue.writeError != nil {
		return queue.writeError
	}
	queue.writeChan <- pack
	return nil
}
func (queue *PackQueue) ReadPack() (pack *spp.Pack, err error) {
	go func() {
		defer recover()
		p := new(packAndErr)
		p.pack, p.err = queue.rw.ReadPack()
		queue.readChan <- p
	}()
	select {
	case err = <-queue.errorChan:
		// Hava an error
		// pass
	case pAndErr := <-queue.readChan:
		pack = pAndErr.pack
		err = pAndErr.err
	}
	return
}

// Only call once
func (queue *PackQueue) ReadPackInLoop() <-chan *packAndErr {
	ch := make(chan *packAndErr, Conf.ReadPackLoop)
	go func() {
		defer recover()
		p := new(packAndErr)
		for {
			p.pack, p.err = queue.rw.ReadPack()
			ch <- p
			if p.err != nil {
				break
			}
			p = new(packAndErr)
		}
	}()
	return ch
}
func (queue *PackQueue) Close() error {
	close(queue.writeChan)
	close(queue.readChan)
	close(queue.errorChan)
	return nil
}
