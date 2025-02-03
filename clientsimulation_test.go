package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

const (
	CONS_COUNT = 0
	PRV_COUNT  = 1
)

func TestClientSimulation(t *testing.T) {

	var wg sync.WaitGroup

	for tag := range PRV_COUNT {
		wg.Add(1)
		providers(&wg, tag, "route")
	}

	wg.Wait()
	for consTag := range CONS_COUNT {
		go consumers(consTag)
	}
	loop := make(chan struct{})
	<-loop
	fmt.Println("Simulation Ended")
}

func providers(wg *sync.WaitGroup, tag int, route string) {
	defer fmt.Printf("NOTIF: Provider #%d exited\n", tag)
	defer wg.Done()
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
	fmt.Println("NOTIF: Successfully created a new Channel")
	_, err = ch.AssertQueue(route, "P2P", false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = ch.DeliverMessage(route, []byte(route), "P2P")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func consumers(tag int) {
	defer fmt.Printf("NOTIF: Consumer #%d exited\n", tag)
	conn, err := Connect("localhost:5671")
	if err != nil {
		log.Panic(err.Error())
	}
	ch, err := conn.CreateChannel()
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("NOTIF: Successfully created a new Channel")
	_, err = ch.AssertQueue("route", "P2P", false)
	if err != nil {
		log.Panic(err.Error())
	}

	msg := ch.Consume("route")

	fmt.Printf("NOTIF: Consumer #%d Waiting for message to consume\n", tag)
	m := <-msg
	fmt.Printf("NOTIF: Received Message from route: %s\n", string(m.Body))

}
