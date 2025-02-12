package msbqclient

import (
	"fmt"
	"github.com/francccisss/msbq-client-api/internal/types"
	"github.com/google/uuid"
	"net"
)

const (
	MAX_STREAMS       = 3
	DEFAULT_READ_SIZE = 50
)

type clientConnection struct {
	conn net.Conn
	// StreamPool contains an array of channels
	// wher each channel corresponds to specific routes
	// that they listen to, the "stream" takes in
	// data from server maps the route that corresponds to
	// a channel

	streamPool map[string]*ClientChannel
}
type Connection interface {
	CreateChannel() (Channel, error)
	Close()
}

func Connect(address string) (Connection, error) {
	conn, err := net.Dial("tcp", address) // need to change this
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	newConnection := &clientConnection{
		conn:       conn,
		streamPool: map[string]*ClientChannel{},
	}

	go Mudem(newConnection)
	fmt.Printf("NOTIF: Successfully Connected to message broker on %s\n", address)
	return newConnection, nil
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
func (c *clientConnection) CreateChannel() (Channel, error) {
	newStreamID := uuid.NewString()
	ch, exists := c.streamPool[newStreamID]
	if !exists {
		ch = &ClientChannel{
			StreamID: newStreamID,
			conn:     c.conn,
			chanBuff: make(chan types.EPMessage),
		}
		c.streamPool[newStreamID] = ch
	}
	fmt.Printf("NOTIF: New Channel created bound to stream %s\n", newStreamID)
	fmt.Printf("NOTIF: ^^ %+v\n", ch)
	fmt.Printf("NOTIF: Stream Pool Length %d\n", len(c.streamPool))
	return ch, nil
}
func (c clientConnection) Close() {}
