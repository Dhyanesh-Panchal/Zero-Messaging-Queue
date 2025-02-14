/*
a stress user which will Generate multiple users generating messages at High pace, to stress test the system.
*/

package main

import (
	"fmt"
	"os"
	"strconv"

	zmt "github.com/pebbe/zmq4"
)

const SERVER_IP string = "10.20.40.165"

// var reader = bufio.NewReader(os.Stdin)

func recieve_msg(context *zmt.Context, self_name string, reciever_sync chan bool) {
	// Subscriber socket connected to central server's 5002 port
	reciever, _ := context.NewSocket(zmt.SUB)
	defer reciever.Close()
	reciever.Connect(fmt.Sprintf("tcp://%s:5002", SERVER_IP))

	// Subscribe to self and broadcst
	reciever.SetSubscribe("@" + self_name + "^")
	reciever.SetSubscribe("@all^")

	recieve_counter := 0
	for {
		// message, _ := reciever.Recv(0)
		reciever.Recv(0)
		recieve_counter++
		// Parse the message
		// message_split := strings.Split(message, " ")
		// _, sender_name, extracted_message := message_split[0], message_split[1], strings.Join(message_split[2:], " ")
		fmt.Printf("\t Reciever %s \n ", self_name)

		if recieve_counter%1 == 0 {
			fmt.Printf("\n %s Recieved %d messages", self_name, recieve_counter)
		}
	}
	reciever_sync <- true
}

// User method
func User(user_id int, total_users int, message_count int, sync chan bool) {
	context, _ := zmt.NewContext()
	defer context.Term()

	// Message pushing socket connected to server at port 5001
	server_snd, _ := context.NewSocket(zmt.PUSH)
	defer server_snd.Close()
	err := server_snd.Connect(fmt.Sprintf("tcp://%s:5001", SERVER_IP))

	if err != nil {
		fmt.Printf("Error Connecting to server for user %d", user_id)
	}
	self_name := strconv.Itoa(user_id)
	// fmt.Printf("Name: %s ", self_name)

	reciever_sync := make(chan bool)
	// start a go routine for listening to published message
	go recieve_msg(context, self_name, reciever_sync)

	// start Sending periodic messages.
	for i := 0; i < message_count; i++ {
		for target_user := 0; target_user < total_users; target_user++ {
			if target_user == user_id {
				continue
			}
			fmt.Printf("\t %d \n ", user_id)
			msg := fmt.Sprintf("THIS IS MESSAGE %d FROM %s", i+1, self_name)

			// // Increament target user
			// target_user = (target_user + 1) % total_users
			// if target_user == user_id {
			// 	target_user = (target_user + 1) % total_users
			// }

			// fmt.Printf(self_name)
			// Push the message to the server
			server_snd.Send(fmt.Sprintf("%s @%s %s", self_name, strconv.Itoa(target_user), msg), 0)
		}
	}
	<-reciever_sync
	sync <- true

}

func main() {
	// get no of users & no of messages/s from argument.
	n_users, _ := strconv.Atoi(os.Args[1])
	n_msgs, _ := strconv.Atoi(os.Args[2])

	sync := make(chan bool)
	//Instanciate Users routines
	for i := 0; i < n_users; i++ {
		go User(i, n_users, n_msgs, sync)
	}
	<-sync
}
