package store

import (
	"labix.org/v2/mgo/bson"
)

const (
	OFFLINE = 11
	ONLINE  = 12
)

type Msg struct {
	Msg_id string // Msg ID
	Body   []byte
	Typ    int

	Owner string // Owner
}

func GetOfflineMsg(mID string, ch chan<- *Msg) {
	defer recover()
	// Find in the db
	sei := sei_msg.New()
	defer sei.Refresh()
	c := sei.DB(Config.MsgName).C(Config.OfflineName)
	iter := c.Find(bson.M{"Msg_id": mID}).Limit(Config.OfflineMsgs).Iter()
	defer iter.Close()
	msg := new(Msg)
	for iter.Next(msg) {
		ch <- msg
		msg = new(Msg)
	}
}
func DelOfflineMsg(msg_id string, id string) {}
