package web

import (
	"strconv"
	"sync"
	"time"
)

var lock = new(sync.Mutex)

func get_uuid() string {
	lock.Lock()
	defer lock.Unlock()
	return strconv.FormatInt(time.Now().Unix(), 10)
}
