package comet

import (
	// "fmt"
	"errors"
	"time"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"

	// "github.com/golang/glog"
)

type client struct {
	queue *PackQueue
	id    string

	offlines chan *store.Msg
	onlines  chan *store.Msg

	// hb       *time.Timer
	// notAlive bool
}

var notAlive = errors.New("Connection was dead")

func (c *client) listen_loop() (e error) {
	defer func() {
		// Close channels
	}()
	go c.queue.writeLoop()
	// Push the offline msg
	// TODO List :
	// Get the offline msg
	store.GetOfflineMsg(c.id, c.offlines)
	// Start push
	var (
		err     error
		msg     *store.Msg
		pAndErr *packAndErr

		readChan = c.queue.ReadPackInLoop()
	)
loop:
	for {
		select {
		// Heart beat
		// case c.hb.C:
		// 	// check whether hava send
		// 	if c.notAlive {
		// 		// Connetion is dead
		// 		e = notAlive
		// 		break loop
		// 	}
		// 	c.notAlive = true
		// 	body, _ := getbeat_heartResp(true)
		// 	resp_pack, _ := c.queue.rw.SetDefaultPack(HEART_BEAT, body)
		// 	err = c.queue.WritePack(pack)
		// 	if err != nil {
		// 		e = err
		// 		break loop
		// 	}
		// Offline msg
		case msg = <-c.offlines:
		case msg = <-c.onlines:
		case pAndErr = <-readChan:
		}
	}
	// Del the offline msg
	// Push the msg
	// Heart beat reply
	return nil
}
func (c *client) doResponse() {

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
