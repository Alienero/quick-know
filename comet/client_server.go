package comet

import (
	// "fmt"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"

	// "github.com/golang/glog"
)

type client struct {
	queue *PackQueue
	id    string

	offlines chan *store.Msg
}

func (c *client) listen_loop() error {
	// Push the offline msg
	// TODO List :
	// Get the offline msg
	// Del the offline msg
	// Push the msg
	// Heart beat reply
	return nil
}

// func InitAllCHs() {
// 	// Heart beat reply
// 	addCh(HEART_BEAT, func(v interface{}, pack *spp.Pack) error {
// 		c := v.(*client)
// 		var err error
// 		body, err := getbeat_heartResp(true)
// 		if err != nil {
// 			return err
// 		}
// 		resp_pack, _ := c.Rw.SetDefaultPack(HEART_BEAT, body)
// 		err = c.writePack(resp_pack)
// 		return err
// 	})
// }

func newClient(rw *spp.Conn, id string) *client {
	return &client{
		queue: NewPackQueue(rw),
		id:    id,
	}
}

// func (c *client) listen_loop() (err error) { return c.l.listen_loop(clients_handles) }
