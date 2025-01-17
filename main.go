package main

import (
	"github.com/francccisss/message-broker-endpoint-api/client"
	"log"
)

func main() {

	conn, err := client.Connect("localhost:8080")
	if err != nil {
		log.Panicf(err.Error())
	}
	_, err = conn.CreateChannel()

}
