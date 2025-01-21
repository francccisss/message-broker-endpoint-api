package main

import (
	"log"
)

func main() {

	conn, err := Connect("localhost:8080")
	if err != nil {
		log.Panicf(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panicf(err.Error())
	}
	log.Println("Successfully created a new Channel")
	chName, err := ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panicf(err.Error())
	}
	log.Println(chName)

	loop := make(chan struct{})
	<-loop
}
