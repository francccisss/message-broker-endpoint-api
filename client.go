package main

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"log"
	"math"
	"net"
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
	go HandleIncomingMessages(conn)
	log.Printf("NOTIF: Successfully Connected to message broker on %s", address)
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

func HandleIncomingMessages(c net.Conn) {
	bodyBuf := make([]byte, READ_SIZE)
	headerBuf := make([]byte, HEADER_SIZE)
	for {
		var msgBuf bytes.Buffer
		_, err := c.Read(headerBuf)
		if err != nil {
			log.Println("ERROR: Unable to decode header prefix length")
			return
		}
		expectedMsgLength := int(binary.LittleEndian.Uint32(headerBuf[:HEADER_SIZE]))
		log.Printf("Prefix Length Receieved: %d\n", expectedMsgLength)

		for {
			_, err := c.Read(bodyBuf)
			if err != nil {
				log.Printf("ERROR: Unable to read the incoming message body ")
				break
			}
			remainingBytes := int(math.Min(float64(expectedMsgLength-msgBuf.Len()), float64(READ_SIZE)))
			// Writes the from the minimum value of remainingBytes into the buffer up to
			// 1024 that is to be read into the bodyBuf
			_, err = msgBuf.Write(bodyBuf[:remainingBytes])
			if err != nil {
				log.Printf("ERROR: Unable to append bytes to the message buffer ")
				break
			}

			log.Printf("Current Total in msgBuf: %+v\n", msgBuf.Len())
			if msgBuf.Len() == expectedMsgLength {
				log.Printf("NOTIF: Receieved all values: %d\n", msgBuf.Bytes())

				log.Printf("BODYBUF BEFORE:\n %+v\n", bodyBuf)

				// Currently head buff is occupied
				// so replace it with approrriate size with the excess from bodyBuf
				// to the headerBuff and leave the rest within the bodyBuf

				// Since TCP is a stream oriented protocol, each new requeust travels in a single
				// connection so to handle excess bytes within the stream, we need to extract
				// and place these excess bytes in to the header and the body buffers
				if len(bodyBuf[remainingBytes:]) < HEADER_SIZE {
					copy(headerBuf, bodyBuf)
					bodyBuf = bodyBuf[:0]
				} else {
					log.Printf("EXTRACTED HEADER LENGTH :%d\n", len(bodyBuf[remainingBytes:remainingBytes+HEADER_SIZE]))

					copy(headerBuf, bodyBuf[remainingBytes:remainingBytes+HEADER_SIZE])
					copy(bodyBuf, bodyBuf[remainingBytes+HEADER_SIZE:])
				}
			}
		}
	}
}
