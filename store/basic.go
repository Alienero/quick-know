package store

import (
	"labix.org/v2/mgo"
)

var sei_user *mgo.Session
var sei_msg *mgo.Session

func Init() (err error) {
	sei_user, err = mgo.Dial(Config.UserAddr)
	sei_user.EnsureSafe(safe * Safe)
	sei_user.SetMode(mgo.Monotonic, true)
	sei_user.Refresh()

	sei_msg, err = mgo.Dial(Config.MsgAddr)
	sei_msg.EnsureSafe(safe * Safe)
	sei_msg.SetMode(mgo.Monotonic, true)
	sei_msg.Refresh()
	return
}

func getOfflineMsg(id string) (string, []byte) {
}
func delOfflineMsg(msg_id string)
