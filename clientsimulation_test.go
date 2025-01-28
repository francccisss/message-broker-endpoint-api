package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

const (
	PRV_COUNT  = 3
	CONS_COUNT = 1
)

func TestClientSimulation(t *testing.T) {

	var wg sync.WaitGroup

	for provTag := range PRV_COUNT {
		wg.Add(1)
		go providers(&wg, provTag)
	}

	wg.Wait()
	for consTag := range CONS_COUNT {
		go consumers(consTag)
	}
	loop := make(chan struct{})
	<-loop
	fmt.Println("Simulation Ended")
}

func providers(wg *sync.WaitGroup, tag int) {
	defer fmt.Printf("Provider #%d exited\n", tag)
	conn, err := Connect("localhost:5671")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("NOTIF: Successfully created a new Channel")
	_, err = ch.AssertQueue("route", "P3P", false)
	if err != nil {
		log.Panic(err.Error())
	}
	err = ch.DeliverMessage("route", []byte("Hello"), "P2P")
	if err != nil {
		log.Println(err.Error())
	}
	wg.Done()
}

func consumers(tag int) {
	defer fmt.Printf("Consumer #%d exited\n", tag)
	conn, err := Connect("localhost:5671")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("NOTIF: Successfully created a new Channel")
	_, err = ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panic(err.Error())
	}

	msg := ch.Consume("route")

	log.Println("Waiting for message to consume")
	m := <-msg
	log.Printf("NOTIF: Received Message from route: %s\n", string(m.Body))

}
