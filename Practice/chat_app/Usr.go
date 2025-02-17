package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

const SERVER_IP string = "10.20.40.165"

var reader = bufio.NewReader(os.Stdin)

func recieve_msg(context *zmq.Context, name_channel chan string) {
	name_verifier_conn, _ := context.NewSocket(zmq.REQ)
	name_verifier_conn.Connect(fmt.Sprintf("tcp://%s:5003", SERVER_IP))

	valid_name := false
	var name string
	for !valid_name {
		// Get the sender's name
		fmt.Printf("Enter your name: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSuffix(name, "\n")

		//verify with server
		name_verifier_conn.Send(name, 0)

		responce, _ := name_verifier_conn.Recv(0)

		if responce == "OK" {
			valid_name = true
		} else {
			fmt.Printf("This user already exist, please enter a unique name.")
		}
	}

	// share the name to main routine
	name_channel <- name

	// Subscriber socket connected to central server's 5002 port
	reciever, _ := context.NewSocket(zmq.SUB)
	defer reciever.Close()
	reciever.Connect(fmt.Sprintf("tcp://%s:5002", SERVER_IP))
	// Subscribe to self and broadcst
	reciever.SetSubscribe("@" + name + "^")
	reciever.SetSubscribe("@all^")
	for {
		message, _ := reciever.Recv(0)

		// Parse the message
		message_split := strings.Split(message, " ")
		_, sender_name, extracted_message := message_split[0], message_split[1], strings.Join(message_split[2:], " ")

		if sender_name != name {
			fmt.Printf("\nFrom %s: %s\n", sender_name, extracted_message)
		}
	}
}

// Just hava ma try karu chu

func main() {
	context, _ := zmq.NewContext()
	defer context.Term()

	// Message pushing socket connected to server at port 5001
	server_snd, _ := context.NewSocket(zmq.PUSH)
	defer server_snd.Close()
	server_snd.Connect(fmt.Sprintf("tcp://%s:5001", SERVER_IP))

	name_channel := make(chan string)
	// start a go routine for listening to published message
	go recieve_msg(context, name_channel)

	//wait for the slef declaration
	self_name := <-name_channel
	for {
		// fmt.Printf("Enter reciever name: ")
		// reciever_name, _ := reader.ReadString('\n')
		// reciever_name = strings.TrimSuffix(reciever_name, "\n")

		fmt.Printf("Enter message: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSuffix(msg, "\n")

		if strings.HasPrefix(msg, "@") {
			server_snd.Send(fmt.Sprintf("%s %s", self_name, msg), 0)
		} else {
			fmt.Println("Please mention the reciepent with '@' at the begining.")
		}

	}

}
