package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	zmt "github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmt.NewContext()
	socket, _ := context.NewSocket(zmt.PUB) // Publisher socket

	reader := bufio.NewReader(os.Stdin)

	defer context.Term()
	defer socket.Close()

	socket.Bind("tcp://*:1234")

	n_msgs, _ := strconv.Atoi(os.Args[1])

	// Continous Broadcast
	for i := 0; i < n_msgs; i++ {

		socket.Send(fmt.Sprintf("%s %s", "T", fmt.Sprintf("MESSAGE #1")), 0)
	}
}
