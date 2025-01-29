package types

type Queue struct {
	MessageSize uint32
	MessageType string
	Name        string
	Type        string
	Durable     bool
	StreamID    string
}

type EPMessage struct {
	MessageSize uint32
	MessageType string
	Route       string
	Type        string
	Body        []byte
	StreamID    string
}

// TODO Figure out how to handle errors from different requests
type ErrorMessage struct {
	MessageSize uint32
	MessageType string
	Body        []byte
	StreamID    string
}
