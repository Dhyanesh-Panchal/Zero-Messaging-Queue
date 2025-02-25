package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
)

var userList = make(map[string]bool)

func handleNewUser(zmqContext *zmq.Context) {
	socket, _ := zmqContext.NewSocket(zmq.REP)
	defer socket.Close()

	socket.Bind("tcp://*:5003")
	for {
		userName, _ := socket.Recv(0)

		// Verify that name isn't duplicate
		if !userList[userName] {
			// userName is unique, update the list
			userList[userName] = true
			// Reply back with OK
			socket.Send("OK", 0)
			fmt.Printf("\n\tNEW USER JOINED: %s", userName)
		} else {
			// User already exists
			socket.Send("USER_ALREADY_EXIST", 0)
		}
	}
}

func main() {
	var wg sync.WaitGroup

	zmqContext, _ := zmq.NewContext()
	defer zmqContext.Term()

	// Receiver socket
	receiver, _ := zmqContext.NewSocket(zmq.PULL)
	receiver.SetLinger(time.Second)
	defer receiver.Close()

	// Bind at port 5001
	receiver.Bind("tcp://*:5001")

	// Publisher socket
	publisher, _ := zmqContext.NewSocket(zmq.PUB)
	publisher.SetLinger(time.Second)
	defer publisher.Close()

	// Bind at port 5002
	publisher.Bind("tcp://*:5002")

	// Instantiate User Name Verification routine
	wg.Add(1)
	go func(zmqContext *zmq.Context) {
		defer wg.Done()
		handleNewUser(zmqContext)
	}(zmqContext)

	// Singular receiver-publisher
	rcvMessageCount := 0
	sndMessageCount := 0
	snapshotTime := time.Now()

	for {
		msg, errRcv := receiver.Recv(0)
		if errRcv == nil {
			rcvMessageCount++
		}

		// Parse the msg and extract the topic. msg format: "topic msg"
		msgSplit := strings.Split(msg, " ")
		senderName, topic, msgContent := msgSplit[0], msgSplit[1], strings.Join(msgSplit[2:], " ")

		fmt.Printf("\n%s --> %s : %s", senderName, topic, msgContent)

		// Publish to that topic
		_, errSnd := publisher.Send(fmt.Sprintf("%s^ %s %s", topic, senderName, msgContent), 0)
		if errSnd == nil {
			sndMessageCount++
		}

		const snapshotLength = 10

		if rcvMessageCount%snapshotLength == 0 || sndMessageCount%snapshotLength == 0 {
			elapsedTime := time.Since(snapshotTime)
			fmt.Printf("\nReceived: %d, Sent: %d messages in %d ms", rcvMessageCount, sndMessageCount, elapsedTime.Milliseconds())
			snapshotTime = time.Now()
		}
	}
}
