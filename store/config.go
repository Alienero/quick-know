package store

var Config *config

type config struct {
	UserAddr string
	UserName string
	Clients  string
	Ctrls    string

	MsgAddr     string
	MsgName     string
	OfflineName string

	SubName  string
	SubsName string

	OfflineMsgs int
}
