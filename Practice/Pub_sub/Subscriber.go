package main

import (
	"bufio"
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq.NewContext()
	socket, _ := context.NewSocket(zmq.SUB)
	reader := bufio.NewReader(os.Stdin)

	defer context.Term()
	defer socket.Close()

	socket.Connect("tcp://10.20.40.165:1234")

	fmt.Printf("Enter Listning Topic: ")
	topic, _ := reader.ReadString('\n')

	socket.SetSubscribe(topic)
	socket.SetSubscribe("Broadcast")
	for {
		data, _ := socket.Recv(0)
		fmt.Printf("Recieved: %s", string(data))
	}

}
