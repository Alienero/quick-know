package comet

import (
	"errors"
	"fmt"
	"time"

	"github.com/Alienero/quick-know/store"
	"github.com/Alienero/spp"

	// "github.com/golang/glog"
)

type client struct {
	queue *PackQueue
	id    string // Owner+id

	offlines chan *store.Msg
	onlines  chan *store.Msg

	onlineCache map[string]*store.Msg
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
		pack    *spp.Pack

		readChan = c.queue.ReadPackInLoop()
	)
loop:
	for {
		select {

		case msg = <-c.offlines:
			// Push the offline msg
			err = c.pushMsg(msg)
			if err != nil {
				break loop
			}
		case msg = <-c.onlines:
			// Push the online msg
			// Add the msg into cache
			if len(c.onlineCache) > Conf.MaxCacheMsg && Conf.MaxCacheMsg != 0 {
				err = fmt.Errorf("Online msg is out of range:%v", len(c.onlineCache))
				break loop
			}
			c.onlineCache[msg.Msg_id] = msg
			err = c.pushMsg(msg)
			if err != nil {
				break
			}
		case pAndErr = <-readChan:
			// If connetion has a error, should break
			// if it return a timeout error, illustrate
			// hava not recive a heart beat pack at an
			// given time.
			if pAndErr.err != nil {
				err = pAndErr.err
				break loop
			}
			// Choose the requst type
			switch pAndErr.pack.Body[1] {
			case OFFLINE:
				// Del the offline msg in the store
				store.DelOfflineMsg(string(packAndErr.pack.Body[1:]), c.id)
			case ONLINE:
				// Del the online msg in the msg cache
				delete(c.onlineCache, pAndErr.pack.Body[1:])
			case HEART_BEAT:
				// Reply the heart beat
				pack, err = c.setPack(HEART_BEAT, []byte("OK"))
				if err != nil {
					break loop
				}
				err = c.queue.WritePack(pack)
				if err != nil {
					break loop
				}
			default:
				err = fmt.Errorf("The type not define:%v", packAndErr.pack.Typ)
				break loop
			}
		}
	}

	return nil
}
func (c *client) pushMsg(msg *store.Msg) (err error) {
	var buf []byte
	buf, err = getMsg(msg)
	if err != nil {
		return
	}
	// Set a pack
	pack, err := c.setPack(PUSH_MSG, buf)
	if err != nil {
		return
	}
	// Write this pack
	err = c.queue.WritePack(pack)
	return
}
func (c *client) setPack(typ int, body []byte) (*spp.Pack, error) {
	return c.queue.rw.SetDefaultPack(typ, body)
}

func newClient(rw *spp.Conn, id string) *client {
	return &client{
		queue: NewPackQueue(rw),
		id:    id,
	}
}
