package main

import (
	"github.com/francccisss/message-broker-endpoint-api/client"
	"log"
)

func main() {

	conn, err := client.Connect("localhost:5671")
	if err != nil {
		log.Panicf(err.Error())
	}
	ch, err := conn.CreateChannel()

}
