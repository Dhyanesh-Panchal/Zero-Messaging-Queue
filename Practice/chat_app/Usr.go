package main

import (
	"bufio"
	"fmt"
	"os"

	zmt "github.com/pebbe/zmq4"
)

func recieve_msg(context *zmt.Context) {
	// Publisher socket
	reciever, _ := context.NewSocket(zmt.SUB)
	defer reciever.Close()
	// Connect at port 5002
	reciever.Connect("tcp://localhost:5002")
	fmt.Printf("Enter your name: ")
	name, _ := reader.ReadString('\n')

	// Subscribe to self and broadcst
	reciever.SetSubscribe(name)
	reciever.SetSubscribe("broadcast")
	for {
		message, _ := reciever.Recv(0)
		// time.Sleep(1 * time.Second)
		fmt.Printf("Recieved: %s", string(message))
	}
}

// Just hava ma try karu chu

func main() {
	reader := bufio.NewReader(os.Stdin)

	context, _ := zmt.NewContext()
	defer context.Term()

	// Reciever socket
	server_snd, _ := context.NewSocket(zmt.PUSH)
	defer server_snd.Close()
	// Connect at port 5001
	server_snd.Connect("tcp://localhost:5001")

	// start a go routine for listening to published message
	go recieve_msg(context)
	for {

		fmt.Printf("Enter reciever name: ")
		reciever_name, _ := reader.ReadString('\n')
		fmt.Printf("Enter message: ")
		msg, _ := reader.ReadString('\n')

		server_snd.Send(fmt.Sprintf("%s %s", reciever_name, msg), 0)

	}

}
