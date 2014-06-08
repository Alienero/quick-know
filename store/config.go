package store

var Config *config

type config struct {
	UserAddr string
	UserName string

	MsgAddr     string
	MsgName     string
	OfflineName string

	OfflineMsgs int
}
