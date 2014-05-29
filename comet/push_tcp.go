package comet

import (
	"net"
	"runtime/debug"

	"github.com/golang/glog"
)

type conn struct {
	rw net.Conn
}

func newConn(rw net.Conn) *conn {
	return &conn{
		rw: rw,
	}
}
func (c *conn) serve() {
	var err error
	defer func() {
		if err = recover(); err != nil {
			glog.Errorf("conn.serve() panic(%v)\n info:%s", err, string(debug.Stack()))
		}
		c.rw.Close()

	}()
	tcp := c.rw.(*net.TCPConn)
	if err = tcp.SetKeepAlive(true); err != nil {
		glog.Errorf("conn.SetKeepAlive() error(%v)\n", err)
		return
	}
	// TODO: 连接认证，离线消息

}
