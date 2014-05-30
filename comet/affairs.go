package comet

import (
	"github.com/Alienero/spp"

	"github.com/golang/glog"
)

var handles = make(map[int]handler)

type handler interface {
	serve(c *client, pack *spp.Pack) error
}
type handle func(c *client, pack *spp.Pack) error

func (h handle) serve(c *client, pack *spp.Pack) error {
	return h(c, pack)
}

func addHandle(typ int, f handle) handler {
	// TODO
}

func Init() {
	// handle[HEART_BEAT] = func(pack *spp.Pack) error {}
}

type client struct {
	rw *spp.Conn
}

func newClient(rw *spp.Conn) *client {
	return &client{
		rw: rw,
	}
}
func (c *client) clientLoop() {
	for {
		// Listen
		pack, err := c.rw.ReadPack()
	}
}
