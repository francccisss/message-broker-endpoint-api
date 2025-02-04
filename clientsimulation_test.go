package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

const (
	CONS_COUNT = 4
	PRV_COUNT  = 4
)

func TestSpawnProviders(t *testing.T) {

	var wg sync.WaitGroup
	for tag := range PRV_COUNT {
		wg.Add(1)
		go providers(&wg, tag, "route")
	}
	wg.Wait()
	fmt.Printf("TEST_NOTIF: Spawned %d Providers\n", PRV_COUNT)
}

func TestSpawnConsumers(t *testing.T) {
	CONS_COUNT := 4
	var wg sync.WaitGroup
	for tag := range CONS_COUNT {
		wg.Add(1)
		go consumers(tag)
	}
	wg.Wait()
	fmt.Printf("TEST_NOTIF: Spawned %d Providers\n", CONS_COUNT)
}

func providers(wg *sync.WaitGroup, tag int, route string) {
	defer fmt.Printf("TEST_NOTIF: Provider #%d exited\n", tag)
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

	fmt.Println("TEST_NOTIF: Successfully created a new Channel")
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
