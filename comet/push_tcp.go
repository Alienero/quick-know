package comet

import (
	"net"
	"runtime/debug"

	"github.com/Alienero/spp"

	"github.com/golang/glog"
)

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
			glog.Errorf("conn.serve() panic(%v)\n info:%s", err, string(debug.Stack()))
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
	pack, err := packRW.ReadPack()
	if err != nil {
		glog.Errorf("Recive login pack error:%v \n", err)
	}
	if !c.login(pack) {
		return
	}
	body, err := getLoginResponse("1", "127.0.0.1", true, "")
	if err != nil {
		return
	}
	pack, _ = packRW.SetDefaultPack(LOGIN, body)
	err = packRW.WritePack(pack)
	if err != nil {
		return
	}
	newClient(packRW).clientLoop()

}
func (c *conn) login(pack *spp.Pack) bool {
	if pack.Typ != LOGIN {
		glog.Errorf("Recive login pack's type error:%v \n", pack.Typ)
		return false
	}
	return true
}
