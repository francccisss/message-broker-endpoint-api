package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	msgType "message-broker-endpoint-api/internal/types"
	"net"
)

type Channel struct {
	StreamID string
	BoundTo  string
	conn     net.Conn
	chanBuff chan msgType.EPMessage
}

type ChannelHandler interface {
	AssertQueue(MessageType string, Name string, Type string, Durable bool) (string, error)
	DeliverMessage(message string)
	Consume(route string, model string, durable bool)
	CloseChannel()
}

/*
  - Name is used to route endpoint messages to the specified queue,
    which then pushes messages to consumers that are subscribe to it,
  - Type could be 'P2P' (point to point) or 'PUBSUB' (pub-sub) model,
  - Durable persists data in the queue
*/
func (ch Channel) AssertQueue(Name string, Type string, Durable bool) (string, error) {
	q := msgType.Queue{
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
	emsg := msgType.EPMessage{
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

// - Route to consume messages from
// creating a route table for streams where if a message of route x is received by the mudem,
// it will try to pattern match the message's data type, parse it and read the route name in
// the message's header, if it exists in the table, the channels/streams connected on that route,
// will receive a message through their chanBuff channel buffer
func (ch Channel) Consume(route string) <-chan msgType.EPMessage {
	r, exists := RouteTable[route]
	if !exists {
		// This channel will consume any message on specified route
		RouteTable[route] = Route{
			Name:     route,
			channels: []*Channel{&ch},
		}
	}
	r.channels = append(r.channels, &ch)
	return ch.chanBuff
}

func (ch Channel) CloseChannel() {
}
