package store

import (
	// "sync"

	"labix.org/v2/mgo"
)

var sei_user *mgo.Session
var sei_msg *mgo.Session

func Init() (err error) {
	connect(sei_user, Config.UserAddr)
	connect(sei_msg, Config.MsgAddr)
}

func connect(sei *mgo.Session, addr string) {
	if sei != nil {
		sei.Close()
	}
	sei, err = mgo.Dial(addr)
	sei.EnsureSafe(&mgo.Safe{})
	sei.SetMode(mgo.Monotonic, true)
	sei.Refresh()
}
