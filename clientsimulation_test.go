package main

import (
	"log"
	"testing"
)

const (
	PRV_COUNT  = 2
	CONS_COUNT = 1
)

func TestClientSimulation(t *testing.T) {

}

func providers() {
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

func consumers() {
	conn, err := Connect("localhost:8080")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("Successfully created a new Channel")
	_, err = ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panic(err.Error())
	}

	msg := ch.Consume("route")
	m := <-msg
	log.Printf("Received Message from route: %s", string(m.Body))

}
