package main

import (
	"log"
	msgType "message-broker-endpoint-api/internal/types"
	"message-broker-endpoint-api/internal/utils"
	"net"
)

type Route struct {
	Name     string
	channels []*Channel
}

var RouteTable = map[string]Route{}

// Multiplexer/Demultiplexer takes in the socket connection
// Mudem will handle the demultiplexing of messages from incoming tcp connection
func Mudem(c net.Conn) {
	for {
		buf := make([]byte, 10)
		_, err := c.Read(buf)
		if err != nil {
			log.Println("NOTIF: Return some error")
			return
		}

		// Message Parsing
		endpointMsg, err := utils.MessageParser(buf)
		if err != nil {
			log.Println("ERROR: Something when wrong with the message parser")
			log.Println(err.Error())
			// just log the error dont process it any further
			continue
		}
		// Pattern matching,
		switch msg := endpointMsg.(type) {
		case msgType.EPMessage:
			// Do a RouteTable Lookup
			MessageDispatcher(msg)
			// I still dont know what to do with different message types
		case msgType.Queue:
			log.Println(msg.MessageType)
		case msgType.ErrorMessage:
			log.Println(msg.MessageType)
			log.Println(string(msg.Body))
		default:
			log.Println("ERROR: Message not of any known type")
		}
	}
}

func MessageDispatcher(msg msgType.EPMessage) {
	log.Println(msg.MessageType)
	route, exists := RouteTable[msg.Route]
	if !exists {
		log.Println("NOTIF: Route does not exist")
		log.Println("NOTIF: Do nothing")
		return
	}
	for _, channel := range route.channels {
		channel.chanBuff <- msg
	}
}
