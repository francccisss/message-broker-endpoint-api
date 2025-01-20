package main

import (
	"log"
)

func main() {

	conn, err := Connect("localhost:8080")
	if err != nil {
		log.Panicf(err.Error())
	}
	_, err = conn.CreateChannel("Queue", "Somewhere", "P2P", true)
	if err != nil {
		log.Panicf(err.Error())
	}

	loop := make(chan struct{})
	<-loop
}
