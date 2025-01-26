package types

type Queue struct {
	MessageSize uint32
	MessageType string
	Name        string
	Type        string
	Durable     bool
}

type EPMessage struct {
	MessageSize uint32
	MessageType string
	Route       string
	Type        string
	Body        []byte
}
type ErrorMessage struct {
	MessageSize uint32
	MessageType string
	Body        []byte
}
