package main

import "testing"

func TestMudem(t *testing.T) {

	conn, err := Connect("localhost:8080")
	if err != nil {
		t.Fatal(err.Error())
	}

	ch, err := conn.CreateChannel()
	if err != nil {
		t.Fatal(err.Error())
	}
	ch.Consume()
}
