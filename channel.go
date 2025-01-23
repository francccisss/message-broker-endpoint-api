package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Channel struct {
	StreamID string
	BoundTo  string
	conn     net.Conn
}

type ChannelHandler interface {
	AssertQueue(MessageType string, Name string, Type string, Durable bool) (string, error)
	DeliverMessage(message string)
	Consume(route string, model string, durable bool)
	CloseChannel()
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
	Body        []byte
}

/*
  - Name is used to route endpoint messages to the specified queue,
    which then pushes messages to consumers that are subscribe to it,
  - Type could be 'P2P' (point to point) or 'PUBSUB' (pub-sub) model,
  - Durable persists data in the queue
*/
func (ch Channel) AssertQueue(Name string, Type string, Durable bool) (string, error) {
	q := Queue{
		Name:        Name,
		MessageType: "Queue",
		Type:        Type,
		Durable:     Durable,
	}
	b, err := json.Marshal(q)
	if err != nil {
		fmt.Printf("ERROR: Unable to Marshal Queue Message")
		return "", err
	}

	_, err = ch.conn.Write(b)
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

/*
Do i need QueueType??
*/
func (ch Channel) DeliverMessage(Route string, Message []byte, QueueType string) error {
	log.Println("Delivering Message")
	emsg := EPMessage{
		MessageType: "EPMessage",
		Route:       Route,
		Type:        QueueType,
		Body:        Message,
	}

	b, err := json.Marshal(emsg)
	if err != nil {
		log.Println("ERROR: Unable to Marshal EPMessage")
		return err
	}

	_, err = ch.conn.Write(b)

	if errors.Is(err, io.EOF) {
		fmt.Printf("ERROR: connection was closed")
		return err
	}
	if err != nil {
		fmt.Printf("ERROR: Unable to write to message broker")
		return err
	}
	return nil
}

func (ch Channel) Consume() {
}

func (ch Channel) CloseChannel() {
}
