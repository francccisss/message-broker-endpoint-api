package main

import (
	"fmt"
	"log"
	"testing"
)

const (
	CONS_COUNT = 1
	PRV_COUNT  = 5
)

func TestSpawnProviders(t *testing.T) {

	l := make(chan struct{})
	for tag := range PRV_COUNT {
		go providers(tag, "route")
	}
	fmt.Printf("TEST_NOTIF: Spawned %d Providers\n", PRV_COUNT)
	<-l
}

func TestSpawnConsumers(t *testing.T) {
	l := make(chan struct{})
	for tag := range CONS_COUNT {
		go consumers(tag)
	}
	fmt.Printf("TEST_NOTIF: Spawned %d Consumers\n", CONS_COUNT)
	<-l
}

func providers(tag int, route string) {
	defer fmt.Printf("TEST_NOTIF: Provider #%d exited\n", tag)
	conn, err := Connect("localhost:5671")
	if err != nil {
		log.Println(err.Error())
		return
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("TEST_NOTIF: Successfully created a new Channel")
	fmt.Println("TEST_NOTIF: Delivering message...")

	err = ch.DeliverMessage(route, []byte(route), "P2P")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func consumers(tag int) {
	defer fmt.Printf("TEST_NOTIF: Consumer #%d exited\n", tag)
	conn, err := Connect("localhost:5671")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("TEST_NOTIF: Successfully created a new Channel")
	_, err = ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panic(err.Error())
	}

	msg := ch.Consume("route")

	fmt.Printf("TEST_NOTIF: Consumer #%d Waiting for message to consume\n", tag)
	m := <-msg
	fmt.Printf("TEST_NOTIF: Received Message from route: %s\n", string(m.Body))

}
