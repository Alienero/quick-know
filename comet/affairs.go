package comet

import (
	"fmt"

	"github.com/Alienero/spp"

	// "github.com/golang/glog"
)

var handles = make(map[int]handler)

type handler interface {
	serve(c *client, pack *spp.Pack) error
}
type handle func(c *client, pack *spp.Pack) error

func (h handle) serve(c *client, pack *spp.Pack) error {
	return h(c, pack)
}

func addHandle(typ int, f handle) {
	handles[typ] = f
}

type listener interface {
	listen_loop() error
}
type client struct {
	// Write loop chan
	writeChan chan *spp.Pack
	// write error
	writeError error

	rw *spp.Conn
}

func InitAllHandles() {
	addHandle(HEART_BEAT, func(c *client, pack *spp.Pack) error {
		var err error
		body, err := getbeat_heartResp(true)
		if err != nil {
			return err
		}
		resp_pack, _ := c.rw.SetDefaultPack(HEART_BEAT, body)
		err = c.writePack(resp_pack)
		return err
	})
}

func newClient(rw *spp.Conn) *client {
	return &client{
		writeChan: make(chan *spp.Pack, Conf.WirteLoopChanNum),
		rw:        rw,
	}
}
func (c *client) listen_loop() (err error) {
	defer func() {
		// Close the res
		close(c.writeChan)
	}()
	var pack *spp.Pack
	for {
		// Listen
		pack, err = c.rw.ReadPack()
		if err != nil {
			// glog.Errorf("clientLoop read pack error:%v\n", err)
			break
		}
		f := handles[pack.Typ]
		if f == nil {
			err = fmt.Errorf("No such pack type:%v", pack.Typ)
			break
		}
		// Call function f
		err = f.serve(c, pack)
		if err != nil {
			// glog.Errorf("clientLoop() f.serve() error:%v\n", err)
			break
		}
	}
	return
}

// Server write queue
func (c *client) writePack(pack *spp.Pack) error {
	if c.writeError != nil {
		return c.writeError
	}
	c.writeChan <- pack
	return nil
}
func (c *client) writeLoop() {
loop:
	for {
		select {
		case pack := <-c.writeChan:
			if pack == nil {
				break loop
			}
			err := c.rw.WritePack(pack)
			if err != nil {
				// Tell listen error
				c.writeError = err
				break loop
			}
		}
	}
}

// Control server
type ControlServer struct {
	conn *spp.Conn
}

func newCServer(rw *spp.Conn) *ControlServer {
	return &ControlServer{
		conn: rw,
	}
}
func (c *ControlServer) listen_loop() error { return nil }
