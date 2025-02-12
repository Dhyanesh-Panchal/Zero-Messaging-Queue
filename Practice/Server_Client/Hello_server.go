package main

import (
	"bufio"
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()
	reader := bufio.NewReader(os.Stdin)
	s, _ := zctx.NewSocket(zmq.REP)
	defer s.Close()
	s.Bind("tcp://10.20.40.165:5555")

	for {
		// Wait for next request from client
		msg, _ := s.RecvBytes(0)
		fmt.Printf("Client: %s\n", string(msg))

		var reply string
		fmt.Print("Server:")
		reply, _ = reader.ReadString('\n')

		// Send reply back to client
		s.Send(reply, 0)
	}
}
