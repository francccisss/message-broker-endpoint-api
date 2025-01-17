package client

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

// This is a wrapper, don't expose this tcp connection
// Returns the current connection with the message broker on the specified port

const MAX_STREAMS = 3

var STREAM_POOL map[string]*Channel

type Channel struct {
	StreamID string
	BoundTo  string
}

type channelHandler interface {
	AsserQueue(route string, model string, durable bool)
	DeliverMessage(message string)
	Consume(route string, model string, durable bool)
	CloseChannel()
}

type connection struct {
	conn *net.Conn
}

// Creates a new channel, a channel is a logical connection between 2 applications,
// it does not know and shouldn't care about the existence of a message broker
func (c connection) CreateChannel() (Channel, error) {
	newStreamID := uuid.NewString()
	ch := Channel{
		StreamID: newStreamID,
	}
	_, exists := STREAM_POOL[newStreamID]
	if !exists {
		STREAM_POOL[newStreamID] = &ch
	}
	return ch, nil
}

func Connect(address string) (connection, error) {
	conn := make(chan net.Conn)
	stopped := make(chan int)
	go connectTCP("localhost:5671", conn, stopped)

	// waiting for succesful tcp connection
	for {
		select {
		case c := <-conn:
			return connection{conn: &c}, nil
		case <-stopped:
			return connection{}, fmt.Errorf("Unable to create a TCP connection with the message broker")
		}
	}
}

func connectTCP(address string, c chan net.Conn, stopped chan int) {
	conn, err := net.Dial("tcp", address) // need to change this
	if err != nil {
		stopped <- 0
		return
	}
	log.Printf("Successfully Connected to message broker on %s", address)
	c <- conn
}
