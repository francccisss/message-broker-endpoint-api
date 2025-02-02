package types

type Queue struct {
	MessageType string
	Name        string
	Type        string
	Durable     bool
	StreamID    string
}

type EPMessage struct {
	MessageType string
	Route       string
	Type        string
	Body        []byte
	StreamID    string
}

// TODO Figure out how to handle errors from different requests
type ErrorMessage struct {
	MessageType string
	Body        []byte
	StreamID    string
}
