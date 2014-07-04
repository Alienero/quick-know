// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mqtt

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	Rserved = iota
	CONNECT
	CONNACK

	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP

	SUBSCRIBE
	SUBACK

	UNSUBSCRIBE
	UNSUBACK

	PINGREQ
	PINGRESP

	DISCONNECT
)

type Pack struct {
	// Fixed header
	msg_type  byte
	dup_flag  byte
	qos_level byte
	retain    byte
	// Remaining Length
	length int

	// Variable header and playload
	variable interface{}
}

type connect struct {
	protocol         string
	version          byte
	keep_alive_timer int
	return_code      byte
	topic_name       string

	user_name     bool
	password      bool
	will_retain   bool
	will_qos      int
	will_flag     bool
	clean_session bool
	rserved       bool

	// Playload
	id         string
	will_topic string
	will_msg   string
	uname      string
	upassword  string
}

type publish struct {
	topic_name string
	mid        int
	msg        []byte
}

// Parse the connect flags
func parse_flags(b byte, flag *connect) {
	if b>>7 != 0 {
		flag.user_name = true
	}
	if b = b & 127; b>>6 != 0 {
		flag.password = true
	}
	if b = b & 63; b>>5 != 0 {
		flag.will_retain = true
	}
	b = b & 31
	flag.will_qos = b >> 3
	if b = b & 7; b>>2 != 0 {
		flag.will_flag = true
	}
	if b = b & 3; b>>1 != 0 {
		flag.clean_session = true
	}
	if b&1 != 0 {
		flag.rserved = true
	}
}

// Read and Write a mqtt pack
func ReadPack(r *bufio.Reader) (pack *Pack, err error) {
	// Read the fixed header
	var (
		fixed     byte
		count_len int
		n         int
		length    = make([]byte, 4)
	)
	fixed, err = r.ReadByte()
	if err != nil {
		return
	}
	// Parse the fixed header
	pack = new(Pack)
	pack.msg_type = fixed >> 4
	fixed = fixed & 15
	pack.dup_flag = fixed >> 3
	fixed = fixed & 7
	pack.qos_level = fixed >> 1
	pack.retain = fixed & 1
	// Get the length of the pack
	length[count_len], err = r.ReadByte()
	if err != nil {
		return
	}
	for length[count_len]>>7 != 0 && count_len < 4 {
		count_len++
		length[count_len], err = r.ReadByte()
		if err != nil {
			return
		}
	}
	temp, e := binary.Varint(length)
	if e < 1 {
		err = fmt.Errorf("Remaining Length error :%v", e)
		return
	}
	pack.length = int(temp)
	// Read the Variable header and the playload
	// Check the msg type
	switch pack.msg_type {
	case CONNECT:
		// Read the protocol name
		var flags byte
		var conn = new(connect)
		pack.variable = conn
		conn.protocol, n, err = readString(r)
		if err != nil {
			break
		}
		if n > (pack.length - 4) {
			err = fmt.Errorf("out of range:%v", pack.length-n)
			break
		}
		// Read the version
		conn.version, err = r.ReadByte()
		if err != nil {
			break
		}
		flags, err = r.ReadByte()
		if err != nil {
			break
		}
		// Read the keep alive timer
		pack.keep_alive_timer, err = readInt(r, 2)
		if err != nil {
			break
		}
		parse_flags(flags, conn)
		// Read the playload
		playload_len := pack.length - 2 - n - 4
		// Read the Client Identifier
		conn.id, n, err = readString(r)
		if err != nil {
			break
		}
		if n > 23 || n < 1 {
			err = fmt.Errorf("Identifier Rejected length is:%v", n)
			conn.return_code = 2
			break
		}
		playload_len -= n
		if n < 1 && (conn.will_flag || conn.password || n < 0) {
			err = fmt.Errorf("length error : %v", playload_len)
			break
		}
		if conn.will_flag {
			// Read the will topic and the will message
			conn.will_topic, n, err = readString(r)
			if err != nil {
				break
			}
			playload_len -= n
			if playload_len < 0 {
				err = fmt.Errorf("length error : %v", playload_len)
				break
			}
			conn.will_msg, n, err = readString(r)
			if err != nil {
				break
			}
			playload_len -= n
		}
		if conn.user_name && playload_len > 0 {
			conn.uname, n, err = readString(r)
			if err != nil {
				break
			}
			playload_len -= n
			if playload_len < 0 {
				err = fmt.Errorf("length error : %v", playload_len)
				break
			}
		}
		if conn.password && playload_len > 0 {
			conn.upassword, n, err = readString(r)
			if err != nil {
				break
			}
			playload_len -= n
			if playload_len < 0 {
				err = fmt.Errorf("length error : %v", playload_len)
				break
			}
		}
	case PUBLISH:
		pub := new(publish)
		pack.variable = pub
		// Read the topic
		pub.topic_name, n, err = readString(r)
		if err != nil {
			break
		}
		vlen := pack.length - n
		if n < 1 || vlen < 2 {
			err = fmt.Errorf("length error :%v", vlen)
			break
		}
		// Read the msg id
		pub.mid, err = readInt(r, 2)
		if err != nil {
			break
		}
		vlen -= 2
		// Read the playload
		pub.msg = make([]byte, vlen)
		_, err = io.ReadFull(r, pub.msg)
	case PINGREQ:
		// Pass
		// Nothing to do
	}

	return
}

func readString(r *bufio.Reader) (s string, nn int, err error) {
	length := make([]byte, 2)
	length[0], err = r.ReadByte()
	if err != nil {
		return
	}
	length[1], err = r.ReadByte()
	if err != nil {
		return
	}
	i, n := binary.Varint(length)
	if n < 1 {
		err = fmt.Errorf("Get the length error:%v", n)
	} else {
		buf := make([]byte, i)
		_, err = io.ReadFull(r, buf)
		if err == nil {
			s = string(buf)
			nn = int64(i)
		}
	}
	return
}
func readInt(r *bufio.Reader, length int) (int, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf[:length])
	if err != nil {
		return 0, err
	}
	i, n := binary.Varint(buf[:length])
	if n < 1 {
		return 0, fmt.Errorf("varint error:%v", n)
	}
	return int(i), nil
}

func WritePack(pack *Pack, w *bufio.Writer) error {}
