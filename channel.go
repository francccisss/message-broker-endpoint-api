package msbqclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/francccisss/msbq-client-api/internal/types"
	"github.com/francccisss/msbq-client-api/internal/utils"
)

type ClientChannel struct {
	StreamID string
	BoundTo  string
	conn     net.Conn
	chanBuff chan types.EPMessage
}

type Channel interface {
	AssertQueue(Name string, Type string, Durable bool) (string, error)
	DeliverMessage(Route string, Message []byte, QueueType string) error
	Consume(route string) <-chan types.EPMessage
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
func (ch ClientChannel) AssertQueue(Name string, Type string, Durable bool) (string, error) {
	fmt.Println("NOTIF: Asserting a Queue")

	// TODO Something wrong with this, server is not able to parse
	// message from queue assertion for some reason
	q := types.Queue{
		Name:        Name,
		MessageType: "Queue",
		Type:        Type,
		Durable:     Durable,
		StreamID:    ch.StreamID,
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
func (ch ClientChannel) DeliverMessage(Route string, Message []byte, QueueType string) error {
	defer fmt.Println("NOTIF: Message delivered!")
	emsg := types.EPMessage{
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
}

// Calling consume returns a read-only channel buffer, this channel's channel buffer
// is used by the Mudem to push messages into it by doing a STREAM_POOL look up, since
// each stream in the table is a pointer to a specific channel, the channel will then
// push the incoming messages into the stream which then pushes the message into the
// the channel's channel buffer
func (ch *ClientChannel) Consume(route string) <-chan types.EPMessage {
	fmt.Printf("NOTIF: Consuming from \"%s\" \n", route)
	// implement some functionality to listen to a specific route
	// a channel can have multiple consumers, which can listen to different routes
	// register the channel to listen for specific routes
	return ch.chanBuff
}

func (ch ClientChannel) CloseChannel() {
}
