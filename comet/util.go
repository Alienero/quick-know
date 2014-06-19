// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"fmt"

	"github.com/Alienero/spp"

	"github.com/golang/glog"
)

// Tcp write queue
type PackQueue struct {
	// The last error in the tcp connection
	writeError error
	// Notice read the error
	errorChan chan error

	writeChan chan *spp.Pack
	readChan  chan *packAndErr
	// Pack connection
	rw *spp.Conn
}
type packAndErr struct {
	pack *spp.Pack
	err  error
}

// Init a pack queue
func NewPackQueue(rw *spp.Conn) *PackQueue {
	return &PackQueue{
		rw:        rw,
		writeChan: make(chan *spp.Pack, Conf.WirteLoopChanNum),
		readChan:  make(chan *packAndErr, 1),
		errorChan: make(chan error, 1),
	}
}

// Start a pack write queue
// It should run in a new grountine
func (queue *PackQueue) writeLoop() {
	// defer recover()
	var err error
loop:
	for {
		select {
		case pack := <-queue.writeChan:
			if pack == nil {
				break loop
			}
			err = queue.rw.WritePack(pack)
			if err != nil {
				// Tell listen error
				queue.writeError = err
				break loop
			}
		}
	}
	// Notice the read
	if err != nil {
		queue.errorChan <- err
	}
}

// Write a pack , and get the last error
func (queue *PackQueue) WritePack(pack *spp.Pack) error {
	if queue.writeError != nil {
		return queue.writeError
	}
	queue.writeChan <- pack
	return nil
}

// Read a pack and retuen the write queue error
func (queue *PackQueue) ReadPack() (pack *spp.Pack, err error) {
	go func() {
		p := new(packAndErr)
		p.pack, p.err = queue.rw.ReadPack()
		queue.readChan <- p
	}()
	select {
	case err = <-queue.errorChan:
		// Hava an error
		// pass
	case pAndErr := <-queue.readChan:
		pack = pAndErr.pack
		err = pAndErr.err
	}
	return
}

// Get a read pack queue
// Only call once
func (queue *PackQueue) ReadPackInLoop(fin <-chan byte) <-chan *packAndErr {
	ch := make(chan *packAndErr, Conf.ReadPackLoop)
	go func() {
		// defer recover()
		p := new(packAndErr)
	loop:
		for {
			p.pack, p.err = queue.rw.ReadPack()
			select {
			case ch <- p:
				// if p.err != nil {
				// 	break loop
				// }
				// Without anything to do
			case <-fin:
				glog.Info("Recive fin (read loop chan)")
				break loop
			}
			p = new(packAndErr)
		}
		close(ch)
	}()
	return ch
}

// Close the all of queue's channels
func (queue *PackQueue) Close() error {
	close(queue.writeChan)
	close(queue.readChan)
	close(queue.errorChan)
	return nil
}

// Buffer
type buffer struct {
	index int
	data  []byte
}

func newBuffer(data []byte) *buffer {
	return &buffer{
		data:  data,
		index: 0,
	}
}
func (b *buffer) readString(length int) (s string, err error) {
	if (length + b.index) > len(b.data) {
		err = fmt.Errorf("Out of range error:%v", length)
		return
	}
	s = string(b.data[b.index:(length + b.index)])
	b.index += length
	return
}
func (b *buffer) readByte() (c byte, err error) {
	if (1 + b.index) > len(b.data) {
		err = fmt.Errorf("Out of range error")
		return
	}
	c = b.data[b.index]
	b.index++
	return
}
