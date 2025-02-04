package main

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

// This is a wrapper, don't expose this tcp connection
// Returns the current connection with the message broker on the specified port

const (
	MAX_STREAMS       = 3
	DEFAULT_READ_SIZE = 50
)

var STREAM_POOL = make(map[string]*ClientChannel)

type connection struct {
	conn net.Conn
}

func Connect(address string) (connection, error) {
	conn, err := net.Dial("tcp", address) // need to change this
	if err != nil {
		fmt.Println(err.Error())
		return connection{}, err
	}
	go Mudem(conn)
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
	ch := ClientChannel{
		StreamID: newStreamID,
		conn:     c.conn,
	}
	_, exists := STREAM_POOL[newStreamID]
	if !exists {
		STREAM_POOL[newStreamID] = &ch
	}
	return ch, nil
}
