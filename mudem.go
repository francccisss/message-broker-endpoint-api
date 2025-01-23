package main

import (
// "log"
// "net"
)

type Route struct {
	Name     string
	channels []*Channel
}

var RouteTable = map[string]Route{}

// Multiplexer/Demultiplexer takes in the socket connection
// Mudem will handle the demultiplexing of messages from incoming tcp connection
// func Mudem(c net.Conn) {
// 	for {
// 		buf := make([]byte, 10)
// 		_, err := c.Read(buf)
// 		if err != nil {
// 			log.Println("Return some error")
// 			return
// 		}
// 	}
//
// }
