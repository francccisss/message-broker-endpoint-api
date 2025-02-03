package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	msgType "message-broker-endpoint-api/internal/types"
	"message-broker-endpoint-api/internal/utils"
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
Channel asserts that a message queue exists on the server, which then connects to that
specified message queue based on the route that the queue is bound to
  - Name is used to route endpoint messages to the specified queue,
    which then pushes messages to consumers that are subscribe to it,
  - Type could be 'P2P' (point to point) or 'PUBSUB' (pub-sub) model,
  - Durable persists data in the queue
*/
func (ch Channel) AssertQueue(Name string, Type string, Durable bool) (string, error) {
	fmt.Println("NOTIF: Asserting a Queue")

	// TODO Something wrong with this, server is not able to parse
	// message from queue assertion for some reason
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

	appQBuff, err := utils.AppendPrefixLength(b)

	if err != nil {
		fmt.Printf("ERROR: Unable to append queue message prefix length")
		return "", err
	}
	_, err = ch.conn.Write(appQBuff)
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
	fmt.Println("NOTIF: Delivering Message...")
	defer fmt.Println("NOTIF: Message delivered!")
	emsg := msgType.EPMessage{
		MessageType: "EPMessage",
		Route:       Route,
		Type:        QueueType,
		Body:        Message,
		StreamID:    ch.StreamID,
	}

	body, err := json.Marshal(emsg)
	if err != nil {
		fmt.Println("ERROR: Unable to Marshal EPMessage")
		return err
	}

	appMessBuff, err := utils.AppendPrefixLength(body)
	if err != nil {
		fmt.Printf("ERROR: Unable to append message prefix length")
		return err
	}
	_, err = ch.conn.Write(appMessBuff)
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

type Consumer struct {
	MessageType string
	Route       string
	StreamID    string
}

// Calling consume returns a read-only channel buffer, this channel's channel buffer
// is used by the Mudem to push messages into it by doing a STREAM_POOL look up, since
// each stream in the table is a pointer to a specific channel, the channel will then
// push the incoming messages into the stream which then pushes the message into the
// the channel's channel buffer
func (ch Channel) Consume(route string) <-chan msgType.EPMessage {
	fmt.Printf("NOTIF: Consuming from %s\n", route)

	b, err := json.Marshal(Consumer{
		MessageType: "Consumer",
		Route:       route,
	})
	if err != nil {
		fmt.Println("ERROR: Unable Marshal Consume Message")
	}
	appConsBuff, err := utils.AppendPrefixLength(b)
	if err != nil {
		fmt.Printf("ERROR: Unable to append consumer message prefix length")
	}
	_, err = ch.conn.Write(appConsBuff)
	if err != nil {
		fmt.Println("ERROR: Unable to create a consumer")
	}

	return ch.chanBuff
}

func (ch Channel) CloseChannel() {
}
