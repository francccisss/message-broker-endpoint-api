package main

import (
	"log"
)

func main() {

	conn, err := Connect("localhost:8080")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("Successfully created a new Channel")
	chName, err := ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println(chName)
	err = ch.DeliverMessage("route", []byte("Hello"), "P2P")
	if err != nil {
		log.Println(err.Error())
	}

	loop := make(chan struct{})
	<-loop
}
