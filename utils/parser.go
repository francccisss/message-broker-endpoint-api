package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
Takes in the incoming message from the message broker server and returns
the unmarshalled messages as an interface{} so that we can perform
a message dispatch or pattern matching based on the message's data type
*/
func MessageParser(message []byte) (interface{}, error) {
	/*
	 Given an empty interface where it can store any value and be represented as any type,
	 we need to assert that its of some known type by matching the "MessageType" of the incoming message,
	 once the "MessageType" of the message is known, we can then Unmarashal the message into the specified
	 known type that matched the "MessageType".
	*/
	var tmp map[string]interface{}
	err := json.Unmarshal(message, &tmp)
	if err != nil {
		log.Println("ERROR: Unable to parse incoming message")
		return nil, err
	}

	/*
	  Unmarshalling json data that can have any type from the server, the json package
	  unmarshalls the bytes into interface maps, and not of concrete type, so in order
	  for us to assert the type of the message returned by the json package, we could
	  assert and Unmarshal again to some known concrete type.
	*/

	switch tmp["MessageType"] {
	case "EPMessage":
		var epMsg EPMessage
		err := json.Unmarshal(message, &epMsg)
		if err != nil {
			return nil, err
		}
		return epMsg, nil
	case "Queue":
		var qMsg Queue
		err := json.Unmarshal(message, &qMsg)
		if err != nil {
			return nil, err
		}
		return qMsg, nil
	case "ErrorMessage":
		var errMsg ErrorMessage
		err := json.Unmarshal(message, &errMsg)
		if err != nil {
			return nil, err
		}
		return errMsg, nil
	default:
		return nil, fmt.Errorf("ERROR: Unable to assert undefined type")
	}

}
