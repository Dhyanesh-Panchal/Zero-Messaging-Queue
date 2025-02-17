/*
a stress user which will Generate multiple users generating messages at High pace, to stress test the system.
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	zmt "github.com/pebbe/zmq4"
)

const SERVER_IP string = "10.20.40.165"

var reciever_ready int = 0

// var reader = bufio.NewReader(os.Stdin)

func recieve_msg(context *zmt.Context, self_name string) {
	// Subscriber socket connected to central server's 5002 port
	reciever, _ := context.NewSocket(zmt.SUB)
	reciever.SetLinger(time.Second)
	reciever.SetRcvhwm(5)
	defer reciever.Close()
	reciever.Connect(fmt.Sprintf("tcp://%s:5002", SERVER_IP))

	// Subscribe to self and broadcst
	reciever.SetSubscribe("@" + self_name + "^")
	reciever.SetSubscribe("@all^")

	recieve_counter := 0

	// Reciever ready increament the global state.
	reciever_ready++
	for {

		_, err_rcv := reciever.Recv(0)
		// fmt.Printf("\n Recieved: %s", message)
		if err_rcv == nil {
			recieve_counter++
		} else {
			fmt.Println("ERROR")
		}
		// Parse the message
		// message_split := strings.Split(message, " ")
		// _, sender_name, extracted_message := message_split[0], message_split[1], strings.Join(message_split[2:], " ")

		if recieve_counter%50 == 0 {
			fmt.Printf("\n %s Recieved %d messages", self_name, recieve_counter)
		}
		// time.Sleep(100 * time.Millisecond)
	}
}

// User method
func User(user_id int, total_users int, message_count int) {
	context, _ := zmt.NewContext()
	defer context.Term()

	// Message pushing socket connected to server at port 5001
	server_snd, _ := context.NewSocket(zmt.PUSH)
	server_snd.SetLinger(time.Second) // Linger for 1 sec
	defer server_snd.Close()
	err := server_snd.Connect(fmt.Sprintf("tcp://%s:5001", SERVER_IP))

	if err != nil {
		fmt.Printf("Error Connecting to server for user %d", user_id)
	}
	// fmt.Printf("\nName: %d ", user_id)
	self_name := strconv.Itoa(user_id)

	// start a go routine for listening to published message
	var wg_reciver sync.WaitGroup
	wg_reciver.Add(1)

	// Use channel to synchronise with the reciever
	go func() {
		defer wg_reciver.Done()
		// fmt.Printf("\nName: %s ", self_name)
		recieve_msg(context, self_name)
	}()
	// // Wait untill all recievers ready
	for reciever_ready < total_users {
		// Wait.
	}
	// start Sending periodic messages.
	fmt.Printf("\nSend started for %d, reciver ready count is: %d", user_id, reciever_ready)
	for i := 0; i < message_count; i++ {
		for target_user := 0; target_user < total_users; target_user++ {
			// time.Sleep(time.Millisecond)
			if target_user == user_id {
				continue
			}
			// fmt.Printf("\t %d \n ", user_id)
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
	wg_reciver.Wait()

}

func main() {
	// get no of users & no of messages/s from argument.
	n_users, _ := strconv.Atoi(os.Args[1])
	n_msgs, _ := strconv.Atoi(os.Args[2])

	var wg sync.WaitGroup

	// A reciever ready shared variable to aware the Users of ready recievers.
	//Instanciate Users routines
	for i := 0; i < n_users; i++ {
		wg.Add(1) // Add one more routine to the waitgroup.

		go func(user_id int) {
			defer wg.Done()
			User(user_id, n_users, n_msgs)
		}(i)
	}

	wg.Wait()
}
