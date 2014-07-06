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
	pack, err := ReadPack(r)
	if err != nil {
		t.Fatal(err)
	}
	// Check the pack
	c := pack.variable.(*connect)
	fmt.Println(*c.uname)
	fmt.Println(*c.upassword)
}
