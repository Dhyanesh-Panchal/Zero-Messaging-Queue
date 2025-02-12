package main

import (
	"bufio"
	"fmt"
	"os"

	zmt "github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmt.NewContext()
	socket, _ := context.NewSocket(zmt.PUB) // Publisher socket

	reader := bufio.NewReader(os.Stdin)

	defer context.Term()
	defer socket.Close()

	socket.Bind("tcp://*:1234")

	// Continous Broadcast
	for {
		fmt.Printf("Enter Publish Topic: ")
		topic, _ := reader.ReadString('\n')
		fmt.Printf("Enter Publish Message: ")
		msg, _ := reader.ReadString('\n')

		socket.Send(fmt.Sprintf("%s %s", topic, msg), 0)
	}
}
