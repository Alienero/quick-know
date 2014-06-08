package store

import (
	"labix.org/v2/mgo/bson"
)

type Msg struct {
	id   string
	body []byte
}

func GetOfflineMsg(id string, ch chan<- *Msg) {
	defer recover()
	// Find in the db
	sei := sei_msg.New()
	defer sei.Refresh()
	c := sei.DB(Config.MsgName).C(Config.OfflineName)
	iter := c.Find(bson.M{"id": id}).Limit(Config.OfflineMsgs).Iter()
	defer iter.Close()
	msg := new(Msg)
	for iter.Next(msg) {
		ch <- msg
		msg = new(Msg)
	}
}
func DelOfflineMsg(msg_id string) {}
