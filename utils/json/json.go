// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"encoding/json"
	"fmt"
)

func Getter(call func() (string, error), v interface{}) error {
	str, err := call()
	if err != nil {
		return err
	}
	fmt.Println("Etcd-->", str)
	return json.Unmarshal([]byte(str), v)
}
