package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	zmt "github.com/pebbe/zmq4"
)

const SERVER_IP string = "10.20.40.165"

var reader = bufio.NewReader(os.Stdin)

func recieve_msg(context *zmt.Context, name_channel chan string) {
	// Subscriber socket connected to central server's 5002 port
	reciever, _ := context.NewSocket(zmt.SUB)
	defer reciever.Close()
	reciever.Connect(fmt.Sprintf("tcp://%s:5002", SERVER_IP))

	// Get the sender's name
	fmt.Printf("Enter your name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSuffix(name, "\n")

	// share the name to main routine
	name_channel <- name

	// Subscribe to self and broadcst
	reciever.SetSubscribe(name)
	reciever.SetSubscribe("broadcast")
	for {
		message, _ := reciever.Recv(0)
		message = strings.TrimSuffix(message, "\n")

		// Parse the message
		message_split := strings.Split(message, " ")
		_, sender_name, extracted_message := message_split[0], message_split[1], strings.Join(message_split[2:], " ")

		if sender_name != name {
			fmt.Printf("\n\nFrom %s: %s\n", sender_name, extracted_message)
		}
	}
}

// Just hava ma try karu chu

func main() {

	context, _ := zmt.NewContext()
	defer context.Term()

	// Message pushing socket connected to server at port 5001
	server_snd, _ := context.NewSocket(zmt.PUSH)
	defer server_snd.Close()
	server_snd.Connect(fmt.Sprintf("tcp://%s:5001", SERVER_IP))

	name_channel := make(chan string)
	// start a go routine for listening to published message
	go recieve_msg(context, name_channel)

	//wait for the slef declaration
	self_name := <-name_channel
	for {
		fmt.Printf("Enter reciever name: ")
		reciever_name, _ := reader.ReadString('\n')
		reciever_name = strings.TrimSuffix(reciever_name, "\n")

		fmt.Printf("Enter message: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\n")

		server_snd.Send(fmt.Sprintf("%s %s %s", reciever_name, self_name, msg), 0)

	}

}
