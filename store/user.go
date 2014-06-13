package store

import (
	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id    string
	Psw   string
	Owner string // Owner
}
type Ctrl struct {
	Id   string
	Auth string
}

func Client_login(id, psw, owner string) bool {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	var u *User
	err := c.Find(bson.M{"ID": id, "Psw": psw, "Owner": owner}).One(u)
	if err != nil {
		glog.Errorf("find user error:%v", err)
		return false
	}
	if u == nil {
		return false
	}
	return true
}
func Ctrl_login(auth string) (bool, string) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Ctrls)
	var u *Ctrl
	err := c.Find(bson.M{"auth": auth}).One(u)
	if err != nil {
		glog.Errorf("find ctrl error:%v", err)
		return false, ""
	}
	if u == nil {
		return false, ""
	}
	return true, u.id
}
func Ctrl_login_alive(id, psw string) bool { return false }

// Add or del user
func AddUser(u *User) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	err := c.Insert(u)
	if err != nil {
		glog.Errorf("Insert a new user error:%v", err)
	}
}
func DelUser(id string, own string) {
	sei := sei_user.New()
	defer sei.Refresh()
	c := sei.DB(Config.UserName).C(Config.Clients)
	err := c.Remove(bson.M{"ID": id, "Owner": own})
	if err != nil {
		glog.Errorf("Del a user error:%v,ID:%v,Own:%v", err, id, own)
	}
}
