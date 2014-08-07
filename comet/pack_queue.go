// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comet

import (
	"bufio"
	"fmt"
	"net"
	"time"

	// "github.com/Alienero/spp"
	"github.com/Alienero/quick-know/mqtt"
)

// Tcp write queue
type PackQueue struct {
	// The last error in the tcp connection
	writeError error
	// Notice read the error
	errorChan chan error

	writeChan chan *pakcAdnType
	readChan  chan *packAndErr
	// Pack connection
	r *bufio.Reader
	w *bufio.Writer

	conn net.Conn
}
type packAndErr struct {
	pack *mqtt.Pack
	err  error
}

// 1 is delay, 0 is no delay, 2 is just flush.
type pakcAdnType struct {
	pack *mqtt.Pack
	typ  byte
}

// Init a pack queue
func NewPackQueue(r *bufio.Reader, w *bufio.Writer, conn net.Conn) *PackQueue {
	return &PackQueue{
		r:         r,
		w:         w,
		conn:      conn,
		writeChan: make(chan *pakcAdnType, Conf.WirteLoopChanNum),
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
		case pt := <-queue.writeChan:
			if pt == nil {
				break loop
			}
			if Conf.WriteTimeout > 0 {
				queue.conn.SetWriteDeadline(time.Now().Add(time.Duration(Conf.WriteTimeout)))
			}
			switch pt.typ {
			case 0:
				err = mqtt.WritePack(pt.pack, queue.w)
			case 1:
				err = mqtt.DelayWritePack(pt.pack, queue.w)
			case 2:
				err = queue.w.Flush()
			}

			if err != nil {
				// Tell listener the error
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
func (queue *PackQueue) WritePack(pack *mqtt.Pack) error {
	if queue.writeError != nil {
		return queue.writeError
	}
	queue.writeChan <- &pakcAdnType{pack: pack}
	return nil
}

func (queue *PackQueue) WriteDelayPack(pack *mqtt.Pack) error {
	if queue.writeError != nil {
		return queue.writeError
	}
	queue.writeChan <- &pakcAdnType{
		pack: pack,
		typ:  1,
	}
	return nil
}

func (queue *PackQueue) Flush() error {
	if queue.writeError != nil {
		return queue.writeError
	}
	queue.writeChan <- &pakcAdnType{typ: 2}
	return nil
}

// Read a pack and retuen the write queue error
func (queue *PackQueue) ReadPack() (pack *mqtt.Pack, err error) {
	go func() {
		p := new(packAndErr)
		if Conf.ReadTimeout > 0 {
			queue.conn.SetReadDeadline(time.Now().Add(time.Duration(Conf.ReadTimeout)))
		}
		p.pack, p.err = mqtt.ReadPack(queue.r)
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
			if Conf.ReadTimeout > 0 {
				queue.conn.SetReadDeadline(time.Now().Add(time.Duration(Conf.ReadTimeout)))
			}
			p.pack, p.err = mqtt.ReadPack(queue.r)
			select {
			case ch <- p:
				// Without anything to do
			case <-fin:
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
