package utils

import (
	"encoding/json"
	"fmt"
	"log"
	protocol "message-broker-endpoint-api/internal/types"
)

// Takes in the incoming message from the message broker server
func MessageParser(message []byte) (interface{}, error) {
	// Given an empty interface where it can store any value and be represented as any type,
	// we need to assert that its of some known type by matching the "MessageType" of the incoming message,
	// once the "MessageType" of the message is known, we can then Unmarashal the message into the specified
	// known type that matched the "MessageType"
	var tmp map[string]interface{}
	err := json.Unmarshal(message, &tmp)
	if err != nil {
		log.Println("ERROR: Unable to parse incoming message")
		return nil, err
	}
	switch tmp["MessageType"] {
	case "EPMessage":
		var epMsg protocol.EPMessage
		err := json.Unmarshal(message, &epMsg)
		if err != nil {
			return nil, err
		}
		return epMsg, nil
	case "Queue":
		var qMsg protocol.Queue
		err := json.Unmarshal(message, &qMsg)
		if err != nil {
			return nil, err
		}
		return qMsg, nil
	case "ErrorMessage":
		var errMsg protocol.ErrorMessage
		err := json.Unmarshal(message, &errMsg)
		if err != nil {
			return nil, err
		}
		return errMsg, nil
	default:
		return nil, fmt.Errorf("ERROR: Unable to assert undefined type")
	}

}
