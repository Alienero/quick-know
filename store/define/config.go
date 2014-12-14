// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package define

var Config = new(DBConfig)

type DBConfig struct {
	DBType string

	UserAddr string // DB addr.
	UserName string // DB name.
	Clients  string
	Ctrls    string
	Salt     string // Sha512's salt.

	MsgAddr     string // DB addr.
	MsgName     string // DB name.
	OfflineName string // Collection name.
	SubName     string // Collection name.
	SubsName    string // Collection name.

	OfflineMsgs int // The max of the offline msgs.
}
