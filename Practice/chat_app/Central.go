package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
)

var User_list = make(map[string]bool)

func Handle_new_usr(context *zmq.Context) {
	soc, _ := context.NewSocket(zmq.REP)
	defer soc.Close()

	soc.Bind("tcp://*:5003")
	for {
		user_name, _ := soc.Recv(0)

		// Verify that name isnt duplicate
		if User_list[user_name] == false {
			// user_name unique, update the list
			User_list[user_name] = true
			// Replay back with OK
			soc.Send("OK", 0)
		} else {
			// User already exist
			soc.Send("USER_ALREADY_EXIST", 0)
		}

	}

}

func main() {

	var wg sync.WaitGroup

	context, _ := zmq.NewContext()
	defer context.Term()

	// Reciever socket
	reciever, _ := context.NewSocket(zmq.PULL)
	reciever.SetLinger(time.Second)

	defer reciever.Close()
	// bind at port 5001
	reciever.Bind("tcp://*:5001")

	// Publisher socket
	publisher, _ := context.NewSocket(zmq.PUB)
	publisher.SetLinger(time.Second)

	defer publisher.Close()
	// bind at port 5002
	publisher.Bind("tcp://*:5002")

	// Instanciate User Name Verification routine
	wg.Add(1)
	go func(context *zmq.Context) {
		defer wg.Done()
		Handle_new_usr(context)
	}(context)

	// singular reciever-publisher

	rcv_message_count := 0
	snd_message_count := 0
	snapshot_time := time.Now()
	for {
		msg, err_rcv := reciever.Recv(0)
		if err_rcv == nil {
			rcv_message_count++
		}
		// parse the msg and extract the topic. msg formt: "topic msg"
		msg_split := strings.Split(msg, " ")
		sender_name, topic, msg_actual := msg_split[0], msg_split[1], strings.Join(msg_split[2:], " ")

		fmt.Printf("\n%s --> %s : %s", sender_name, topic, msg_actual)

		// Publish to that topic
		_, err_snd := publisher.Send(fmt.Sprintf("%s^ %s %s", topic, sender_name, msg_actual), 0)
		if err_snd == nil {
			snd_message_count++
		}

		const snapshot_length = 10

		if rcv_message_count%snapshot_length == 0 || snd_message_count%snapshot_length == 0 {
			elapsed_time := time.Since(snapshot_time)
			fmt.Printf("\nRecieved: %d, Sent: %d messages in %d ms", rcv_message_count, snd_message_count, elapsed_time.Milliseconds())
			snapshot_time = time.Now()
		}

	}

}
