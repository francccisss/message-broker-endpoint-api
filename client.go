package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

// This is a wrapper, don't expose this tcp connection
// Returns the current connection with the message broker on the specified port

const MAX_STREAMS = 3

var STREAM_POOL = make(map[string]*Channel)

type connection struct {
	conn net.Conn
}

func Connect(address string) (connection, error) {
	conn := make(chan net.Conn)
	stopped := make(chan int)
	go connectTCP(address, conn, stopped)

	// waiting for succesful tcp connection
	for {
		select {
		case c := <-conn:
			return connection{conn: c}, nil
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
