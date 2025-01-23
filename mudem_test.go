package main

import (
	"log"
	msgType "message-broker-endpoint-api/internal/types"
	"message-broker-endpoint-api/internal/utils"
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

		// Message Parsing
		endpointMsg, err := utils.MessageParser(buf)
		if err != nil {
			log.Println(err.Error())
		}
		// Pattern matching,
		switch msg := endpointMsg.(type) {
		case msgType.EPMessage:
			// Do a RouteTable Lookup
			MessageDispatcher(msg)
		case msgType.Queue:
			log.Println(msg.MessageType)
		case msgType.ErrorMessage:
			log.Println(msg.MessageType)

		}
	}
}

func MessageDispatcher(msg msgType.EPMessage) {
	log.Println(msg.MessageType)
	route, exists := RouteTable[msg.Route]
	if !exists {
		log.Println("Route does not exist")
		log.Println("Do nothing")
		return
	}
	for _, channel := range route.channels {
		channel.chanBuff <- msg
	}
}
