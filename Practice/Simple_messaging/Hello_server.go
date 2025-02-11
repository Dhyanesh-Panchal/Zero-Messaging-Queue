package main

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.REP)
	s.Bind("tcp://10.20.40.165:5555")

	for {
		// Wait for next request from client
		msg, _ := s.Recv(0)
		log.Printf("Client: %s\n", msg)

		var reply string
		fmt.Print("Server:")
		fmt.Scanln(&reply)

		// Send reply back to client
		s.Send(reply, 0)
	}
}
