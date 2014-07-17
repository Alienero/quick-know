// Copyright © 2014 Alienero. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mqtt

import (
	"bufio"
	"encoding/binary"
	"errors"
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
	protocol         *string
	version          byte
	keep_alive_timer int
	return_code      byte
	topic_name       *string

	user_name     bool
	password      bool
	will_retain   bool
	will_qos      int
	will_flag     bool
	clean_session bool
	rserved       bool

	// Playload
	id         *string
	will_topic *string
	will_msg   *string
	uname      *string
	upassword  *string
}

type connack struct {
	reserved    byte
	return_code byte
}

type publish struct {
	topic_name *string
	mid        int
	msg        []byte
}

type puback struct {
	mid int
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
	flag.will_qos = int(b >> 3)
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
		n         int
		temp_byte byte
		count_len = 1
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
	temp_byte, err = r.ReadByte()
	if err != nil {
		return
	}

	// Read the high
	multiplier := 1
	for {
		count_len++
		pack.length += (int(temp_byte&127) * multiplier)
		if temp_byte>>7 != 0 && count_len < 4 {
			temp_byte, err = r.ReadByte()
			if err != nil {
				return
			}
			multiplier *= 128
		} else {
			break
		}
	}
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
		conn.keep_alive_timer, err = readInt(r, 2)
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
	case PUBACK:
		if pack.length == 2 {
			ack := new(puback)
			ack.mid, err = readInt(r, 2)
			if err != nil {
				break
			}
			pack.variable = ack
		} else {
			err = fmt.Errorf("Pack(%v) length(%v) != 2", pack.msg_type, pack.length)
		}
	case PINGREQ:
		// Pass
		// Nothing to do
	}

	return
}

func readString(r *bufio.Reader) (s *string, nn int, err error) {
	temp_string := ""
	s = &temp_string
	nn, err = readInt(r, 2)
	if err != nil {
		return
	}
	if nn > 0 {
		buf := make([]byte, nn)
		_, err = io.ReadFull(r, buf)
		if err == nil {
			*s = string(buf)
		}
	} else {
		*s = ""
	}
	return
}
func readInt(r *bufio.Reader, length int) (int, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf[:length])
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(buf[:length])), nil
}

func WritePack(pack *Pack, w *bufio.Writer) (err error) {
	// Write the fixed header
	var fixed byte
	// Byte 1
	fixed = pack.msg_type << 4
	fixed |= pack.dup_flag << 3
	fixed |= pack.qos_level << 1
	if err = w.WriteByte(fixed); err != nil {
		return
	}
	// Byte2
	switch pack.msg_type {
	case CONNACK:
		ack := pack.variable.(*connack)
		err = writeFull(w, getRemainingLength(2))
		if err != nil {
			return
		}
		// Write the variable
		if err = writeFull(w, []byte{ack.reserved, ack.return_code}); err != nil {
			return
		}
	case PUBLISH:
		// Publish the msg to the client
		pub := pack.variable.(*publish)
		if err = writeFull(w, getRemainingLength(4+len([]byte(*pub.topic_name)))); err != nil {
			return
		}
		if err = writeString(w, pub.topic_name); err != nil {
			return
		}
		if err = writeInt(w, pub.mid, 2); err != nil {
			return
		}
		if err = writeFull(w, pub.msg); err != nil {
			return
		}
	}
	err = w.Flush()
	return
}

func getRemainingLength(length int) []byte {
	b := make([]byte, 4)
	count := 0
	for {
		digit := length % 128
		length = length / 128
		if length > 0 {
			digit |= 128
			b[count] = byte(digit)
		} else {
			b[count] = byte(digit)
			break
		}
		count++
	}
	return b[:count+1]
}

func writeString(w *bufio.Writer, s *string) error {
	// Write the length of the string
	if s == nil {
		return errors.New("nil pointer")
	}
	data := []byte(*s)
	// Write the string length
	err := writeInt(w, len(data), 2)
	if err != nil {
		return err
	}
	return writeFull(w, data)
}
func writeInt(w *bufio.Writer, i, size int) error {
	b := make([]byte, size)
	binary.BigEndian.PutUint16(b, uint16(i))
	return writeFull(w, b)
}

// wirteFull write the data into the Writer's buffer
func writeFull(w *bufio.Writer, b []byte) (err error) {
	hasRead, n := 0, 0
	for n == len(b) {
		n, err = w.Write(b[hasRead:])
		if err != nil {
			break
		}
		hasRead += n
	}
	return err
}
