package main

import (
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
	conn, err := net.Dial("tcp", address) // need to change this
	if err != nil {
		log.Println(err.Error())
		return connection{}, err
	}
	go Mudem(conn)
	log.Printf("NOTIF: Successfully Connected to message broker on %s", address)
	return connection{conn}, nil
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
