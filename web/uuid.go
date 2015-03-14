// Copyright © 2015 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/nu7hatch/gouuid"
)

func getID() (string, error) {
	uu, err := uuid.NewV4()
	return uu.String(), err
}
