package comet

import (
	"net"
	"time"

	"github.com/golang/glog"
)

func StartListen() error {
	l, err := net.Listen("tcp", Conf.Listen_addr)
	if err != nil {
		return err
	}
	for {
		rw, err := l.Accept()
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
