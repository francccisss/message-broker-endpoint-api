package types

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
type ErrorMessage struct {
	MessageType string
	Body        []byte
}
