package main

import (
	"log"
	"net"
	"testing"
)

func TestMudem(t *testing.T) {

	conn, err := Connect("localhost:5671")
	if err != nil {
		t.Fatal(err.Error())
	}
	go Mudem(conn.conn)

	ch, err := conn.CreateChannel()
	if err != nil {
		t.Fatal(err.Error())
	}
	_ = ch.Consume("route")
	log.Println("Not Blocked")

	loop := make(chan struct{})
	<-loop
}

func Mudem(c net.Conn) {
	for {
		buf := make([]byte, 10)
		_, err := c.Read(buf)
		if err != nil {
			log.Println("Return some error")
			return
		}
		// Do pattern matching,
		//  - Do RouteTable look up after parsing
		// Error receiving,
		// Parse message
	}

}
