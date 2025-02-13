package msbqclient

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"

	"github.com/francccisss/msbq-client-api/utils"
)

const READ_SIZE = 1024
const HEADER_SIZE = 4

// Multiplexer/Demultiplexer takes in the socket connection
// Mudem will handle the demultiplexing of messages from incoming tcp connection
// using message dispatch based on the message's type

func mudem(c *clientConnection) {
	msgChan := make(chan []byte)
	go HandleIncomingMessage(c.conn, msgChan)

	// Code does not run until msgChan has data in its channel buffer
	// from the HandleIncomingMessage()
	for {
		msg := <-msgChan
		go dispatchMessage(msg, &c.streamPool)
	}
}

/*
Handling incoming messages to be dispatched to different channels,
each channel is bound to a specific stream, the stream will contain
the a pointer to the channel, messages will be pushed into the channel's
channel buffer
*/
func dispatchMessage(incomingMessage []byte, streamPool *map[string]*ClientChannel) {
	fmt.Println("NOTIF: Dispatching message...")

	msg, err := utils.MessageParser(incomingMessage)
	if err != nil {
		fmt.Printf("ERROR: Unable to parse message")
		fmt.Println(err.Error())
		return
	}
	// Type assertion to marashal incoming json stream as
	// concrete type defined in the package's message types
	switch m := msg.(type) {
	case utils.EPMessage:
		chann, exists := (*streamPool)[m.StreamID]
		if !exists {
			fmt.Println("NOTIF: Stream does not exist")
			fmt.Println("NOTIF: Do nothing")
			return
		}
		var epMsg utils.EPMessage
		err := json.Unmarshal(incomingMessage, &epMsg)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		chann.chanBuff <- epMsg
	case utils.ErrorMessage:
		fmt.Println(m.MessageType)
	default:
		fmt.Println("ERROR: Unidentified type")
	}
}

// Using prefix length for handling tcp data stream
func HandleIncomingMessage(c net.Conn, msgChan chan []byte) {
	defer fmt.Println("NOTIF: Exiting request listener")

	defer c.Close()

	headerBuf := make([]byte, HEADER_SIZE)
	for {
		var msgBuf bytes.Buffer
		_, err := c.Read(headerBuf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("ERROR: Client has abrubtly terminated the connection")
				return
			}
		}

		expectedMsgLength := int(binary.LittleEndian.Uint32(headerBuf[:HEADER_SIZE]))

		// Initial read size to check if whether incoming request exceeds
		// current read size then read up to DEFAULT_READ_SIZE else only
		// up to remaning bytes to be read instead
		// This formula returns the minimum int between the two, if there is space
		// to fit the stream of bytes in the bodyBuf then return current readSize which is the DEFAULT_READ_SIZE
		// else if current readSize is greater than the remaining bytes left from the expected message
		// return n bytes up to the length of the remaining bytes of the current message.
		currentReadSize := int(math.Min(float64(expectedMsgLength-msgBuf.Len()), float64(DEFAULT_READ_SIZE)))
		for {

			// creates a buffer up to the calculated readSize
			bodyBuf := make([]byte, currentReadSize)

			_, err := c.Read(bodyBuf)
			if err != nil {
				fmt.Printf("ERROR: Unable to read the incoming message body ")
				break
			}

			// store bytes from stream up to the current readsize length into the
			// msgBuf (msgBuf is the current accumulated requested stream from client)
			_, err = msgBuf.Write(bodyBuf[:])
			if err != nil {
				fmt.Printf("ERROR: Unable to write incoming bytes to the message buffer ")
				break
			}

			// Updates the readSize for the next bytes to be read
			remainingBytesLen := expectedMsgLength - msgBuf.Len()
			if currentReadSize < remainingBytesLen {
				currentReadSize = remainingBytesLen
			}
			// finishes the current stream request
			if msgBuf.Len() == expectedMsgLength {
				msgChan <- msgBuf.Bytes()
				break
			}
		}
	}
}
