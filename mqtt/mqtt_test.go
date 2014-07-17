package mqtt

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func TestConnet(t *testing.T) {
	// Set the listener
	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		t.Fatal(err)
	}
	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	pack, err := ReadPack(r)
	if err != nil {
		t.Fatal(err)
	}
	// Check the pack
	c := pack.variable.(*connect)
	fmt.Println(*c.uname)
	fmt.Println(*c.upassword)
	fmt.Println(*c.id)

	// Return the connection ack
	pack = new(Pack)
	pack.msg_type = CONNACK

	ack := new(connack)
	ack.return_code = 0
	pack.variable = ack

	if err := WritePack(pack, w); err != nil {
		t.Error(err)
		return
	}

	pack = new(Pack)
	pack.qos_level = 1
	pack.dup_flag = 0
	pack.msg_type = PUBLISH
	pub := new(publish)
	pub.mid = 1
	s := "jcode/a"
	pub.topic_name = &s
	pub.msg = []byte("Hello push server")
	pack.variable = pub
	if err := WritePack(pack, w); err != nil {
		t.Error(err)
		return
	}

	if _, err := ReadPack(r); err != nil {
		t.Error(err)
	}
}
