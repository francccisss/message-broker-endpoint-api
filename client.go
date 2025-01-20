package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

// This is a wrapper, don't expose this tcp connection
// Returns the current connection with the message broker on the specified port

const MAX_STREAMS = 3

var STREAM_POOL = make(map[string]*Channel)

type connection struct {
	conn *net.Conn
}

func Connect(address string) (connection, error) {
	conn := make(chan net.Conn)
	stopped := make(chan int)
	go connectTCP(address, conn, stopped)

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

/*
  - Name is used to route endpoint messages to the specified queue,
    which then pushes messages to consumers that are subscribe to it,
  - Type could be `P2P` (point to point) or `PSB` (pub-sub) model,
  - Durable persists data in the queue
*/
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

type Channel struct {
	StreamID string
	BoundTo  string
	conn     *net.Conn
}

type ChannelHandler interface {
	AssertQueue(MessageType string, Name string, Type string, Durable bool) (string, error)
	DeliverMessage(message string)
	Consume(route string, model string, durable bool)
	CloseChannel()
}

func (ch Channel) AssertQueue(Name string, MessageType string, Type string, Durable bool) (string, error) {
	q := Queue{
		Name,
		MessageType,
		Type,
		Durable,
	}
	b, err := json.Marshal(q)
	if err != nil {
		fmt.Printf("ERROR: Unable to marshal queue message")
		return "", err
	}

	c := *ch.conn
	_, err = c.Write(b)
	if errors.Is(err, io.EOF) {
		fmt.Printf("ERROR: connection was closed")
		return "", err
	}
	if err != nil {
		fmt.Printf("ERROR: Unable to write to message broker")
		return "", err
	}

	return q.Name, nil
}

func (ch Channel) DeliverMessage() {
}

func (ch Channel) Consume() {
}

func (ch Channel) CloseChannel() {
}

type Queue struct {
	MessageType string
	Name        string
	Type        string
	Durable     bool
}

type EPMessage struct {
	MessageType string
	Route       string
	Type        string
	Body        any
}

func CreateEPMessage(MessageType string, Route string, Type string, Body any) EPMessage {
	return EPMessage{
		MessageType,
		Route,
		Type,
		Body,
	}
}
