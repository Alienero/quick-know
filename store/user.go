package store

import (
	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

type User struct {
	id    string
	psw   string
	owner string // Owner
}
type Ctrl struct {
	id  string
	psw string
}

func Client_login(id, psw, owner string) bool {
	sei := sei_user.New()
	c := sei.DB(Config.UserName).C(Config.Clients)
	var u *User
	err := c.Find(bson.M{"id": id, "psw": psw, "owner": owner}).One(u)
	if err != nil {
		glog.Errorf("find user error:%v", err)
		return false
	}
	if u == nil {
		return false
	}
	return true
}
func Ctrl_login(id, psw string) bool {
	sei := sei_user.New()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	var u *Ctrl
	err := c.Find(bson.M{"id": id, "psw": psw}).One(u)
	if err != nil {
		glog.Errorf("find ctrl error:%v", err)
		return false
	}
	if u == nil {
		return false
	}
	return true
}
