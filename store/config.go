// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
