package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

func Reciever(soc zmq.Socket) {
	for {
		msg, _ := soc.Recv(0)
		fmt.Printf("Received reply [ %s ]\n", msg)
	}
}

func Send(soc zmq.Socket, msg string) {
	soc.Send(msg, 0)
}

func main() {
	zctx, _ := zmq.NewContext()
	// Socket to talk to server
	fmt.Printf("Connecting to the server...\n")
	s, _ := zctx.NewSocket(zmq.REQ)
	s.Connect("tcp://10.20.40.165:5555")

	// Do 10 requests, waiting each time for a response
	for {
		var mssg string
		fmt.Printf("Client:")
		fmt.Scanln(&mssg)
		s.Send(mssg, 0)

		msg, _ := s.Recv(0)
		fmt.Printf("Server: %s \n", msg)
	}
}
