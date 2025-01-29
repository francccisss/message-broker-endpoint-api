package main

import (
	"log"
	msgType "message-broker-endpoint-api/internal/types"
	"message-broker-endpoint-api/internal/utils"
)

const READ_SIZE = 1024
const HEADER_SIZE = 4

// Multiplexer/Demultiplexer takes in the socket connection
// Mudem will handle the demultiplexing of messages from incoming tcp connection
// using message dispatch based on the message's type
func Mudem(incomingMessage []byte) {

	msg, err := utils.MessageParser(incomingMessage)
	if err != nil {
		log.Printf("ERROR: Unable to parse message")
		log.Println(err.Error())
		return
	}

	// type assertion switch statement for different processing

	switch m := msg.(type) {
	case msgType.EPMessage:
		log.Println(m.MessageType)

		route, exists := RouteTable[m.Route]
		if !exists {
			log.Println("NOTIF: Route does not exist")
			log.Println("NOTIF: Do nothing")
			return
		}
		// Handling incoming messages to be dispatched to different channels
		// TODO Make STREAMS send data to specific channels
		// instead of basing the message delivery on the channels
		// connected route, each channel is bound to a specific stream,
		// the stream will contain the message the is demultiplexed by the mudem
		// and push it to the specified channel based on the streams, streamID

		// - Message received
		// - Message parsed
		// - Mudem reads message
		// - Mudem Looks at message's streamID
		// - Mudem looks up the stream pool
		// - The stream pool contains channels using the stream
		// - Each Stream contains a pointer to a Channel

		for _, channel := range route.channels {
			channel.chanBuff <- m
		}
	case msgType.ErrorMessage:
		log.Println(m.MessageType)
	case msgType.Queue:
		log.Println(m.MessageType)
	default:
		log.Println("ERROR: Unidentified type")
	}
}
