# msbq-client-api

The client side go package used with [Message Broker](https://github.com/francccisss/msbq-server) for enabling application endpoints to communication through a message broker middleware.

Mind you this is just a small package, in a small project used for learning, it is not meant to be a fully fledged rabbitmq clone but just implementation of basic functionality.

### How to use

```go
    // Connect to the message broker using TCP
    conn, err := Connect("localhost:5671")

    // Creates a new Channel and a stream
    ch, err := conn.CreateChannel()

    // Asserts that queue exists and creates one if it doesnt then
    // registers the connection to specific message queue defined
    // defined by its route
    _, err = ch.AssertQueue("route", "P2P", false)


    // To consume and any mssages sent to the channel
    msg := ch.Consume("route")
    m := <-msg


    // For dispatching message to specific message queues
	err = ch.DeliverMessage(route, []byte(route), "P2P")
```
