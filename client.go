package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/google/uuid"
)

// This is a wrapper, don't expose this tcp connection
// Returns the current connection with the message broker on the specified port

const (
	MAX_STREAMS       = 3
	DEFAULT_READ_SIZE = 50
)

var STREAM_POOL = make(map[string]*Channel)

type connection struct {
	conn net.Conn
}

func Connect(address string) (connection, error) {
	conn, err := net.Dial("tcp", address) // need to change this
	if err != nil {
		fmt.Println(err.Error())
		return connection{}, err
	}
	go HandleIncomingRequests(conn)
	fmt.Printf("NOTIF: Successfully Connected to message broker on %s\n", address)
	return connection{conn}, nil
}

// # Creates a stream and channel
//   - A Channel is an abstracted logical connection between two application endpoints
//   - A Stream is a representation of data flowing from a tcp socket connection
//
// With multiplexing and demultiplexing the protocol can create multiple streams
// dedicated to specific channels, we can look at them as lightweight tcp connections.
//
// eg: CreateChannel() creates a new channel where it listens to the specific stream
// for incoming data
//
// # How are different streams handled?
//
// Each in the Hashmap holds a pointer to a specific channel on creation where the mudem
// can push messages into it using the channel's channel buffer (the "channel buffer' for IPC not
// to be confused with messaging systems' concept of channels).
//
// The Mudem() is responsible for handling different messages coming from the message broker
// and parses and then uses the STREAM_POOL look up table to push new messages into the specific
// stream that is specified on the StreamID field that is included in every message type
func (c connection) CreateChannel() (Channel, error) {
	newStreamID := uuid.NewString()
	ch := Channel{
		StreamID: newStreamID,
		conn:     c.conn,
	}
	_, exists := STREAM_POOL[newStreamID]
	if !exists {
		STREAM_POOL[newStreamID] = &ch
	}
	return ch, nil
}

// TODO Update this using using the server's implementation of handling
// of incoming requests
func HandleIncomingRequests(c net.Conn) {
	connBufReader := bufio.NewReader(c)
	defer c.Close()

	// fixed sized header length to extract from message stream
	// Outer loop will always take 4 iterations
	headerBuf := make([]byte, HEADER_SIZE)
	readSize := DEFAULT_READ_SIZE
	for {
		var msgBuf bytes.Buffer
		_, err := c.Read(headerBuf)
		if err != nil {
			fmt.Println("ERROR: Unable to decode header prefix length")
			return
		}

		expectedMsgLength := int(binary.LittleEndian.Uint32(headerBuf[:HEADER_SIZE]))
		fmt.Printf("Prefix Length Receieved: %d\n", expectedMsgLength)
		for {
			bodyBuf := make([]byte, readSize)
			_, err := connBufReader.Read(bodyBuf)
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

			// Updates the readsize for the next stream of bytes to be captured in bulk
			// This formula returns the minimum int between the two, if there is space
			// to fit the stream of bytes in the bodyBuf then return current readSize which is the DEFAULT_READ_SIZE
			// else if current readSize is greater than the remaining bytes left from the expected message
			// return n bytes up to the length of the remaining bytes of the current message.
			readSize = int(math.Min(float64(expectedMsgLength-msgBuf.Len()), float64(readSize)))

			// finishes the current stream request
			if msgBuf.Len() == expectedMsgLength {
				readSize = DEFAULT_READ_SIZE
				fmt.Println("NOTIF: Message sequence complete.")
				fmt.Println("NOTIF: Do something with the new request.")
				break
			}
		}
	}
}
