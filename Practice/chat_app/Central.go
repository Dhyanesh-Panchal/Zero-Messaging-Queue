package main

import (
	"fmt"
	"strings"
	"time"

	zmt "github.com/pebbe/zmq4"
)

// func recieve_msg(reciever *zmt.Socket)

// Just hava ma try karu chu

func main() {
	context, _ := zmt.NewContext()
	defer context.Term()

	// Reciever socket
	reciever, _ := context.NewSocket(zmt.PULL)
	defer reciever.Close()
	// bind at port 5001
	reciever.Bind("tcp://*:5001")

	// Publisher socket
	publisher, _ := context.NewSocket(zmt.PUB)
	defer publisher.Close()
	// bind at port 5002
	publisher.Bind("tcp://*:5002")

	// singular reciever-publisher

	message_count := 0
	snapshot_time := time.Now()
	for {
		msg, _ := reciever.Recv(0)
		message_count++
		// parse the msg and extract the topic. msg formt: "topic msg"
		msg_split := strings.Split(msg, " ")
		sender_name, topic, msg_actual := msg_split[0], msg_split[1], strings.Join(msg_split[2:], " ")

		// fmt.Printf("\n%s --> %s : %s", sender_name, topic, msg_actual)

		// Publish to that topic
		publisher.Send(fmt.Sprintf("%s^ %s %s", topic, sender_name, msg_actual), 0)

		if message_count%1000 == 0 {
			elapsed_time := time.Since(snapshot_time)
			fmt.Printf("\nSent %d messages in %d ms", message_count, elapsed_time.Milliseconds())
			snapshot_time = time.Now()
		}

	}

}
