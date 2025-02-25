package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type Users struct {
	container map[string]bool
	lock      sync.Mutex
}

// Globals
const snapshotLength = 100
const sharedMessageBufferSize = 100

func handleUserManagementMessage(message string, userList *Users) string {
	userName := message
	userList.lock.Lock()
	defer userList.lock.Unlock()
	// Verify that name isn't duplicate
	if !userList.container[userName] {
		// userName is unique, update the list
		userList.container[userName] = true
		// Reply
		fmt.Printf("\n\tNEW USER JOINED: %s", userName)
		return "OK"
	} else {
		// User already exists
		return "USER_ALREADY_EXIST"
	}
}

func senderRoutine(zmqContext *zmq.Context, sharedMsgBuffer chan string) {
	var sndMessageCount int

	// Publisher socket
	publisher, _ := zmqContext.NewSocket(zmq.PUB)
	publisher.SetLinger(time.Second)
	defer publisher.Close()

	// Bind at port 5002
	publisher.Bind("tcp://*:5002")

	// Listen for continuous messages
	for {
		message := <-sharedMsgBuffer
		msgSplit := strings.Split(message, " ")
		senderName, topic, msgContent := msgSplit[0], msgSplit[1], strings.Join(msgSplit[2:], " ")

		// Publish to that topic
		_, errSnd := publisher.Send(fmt.Sprintf("%s^ %s %s", topic, senderName, msgContent), 0)
		if errSnd == nil {
			sndMessageCount++
			if sndMessageCount%snapshotLength == 0 {
				fmt.Printf("\nReceive count: %d", sndMessageCount)
			}
		}
	}
}

func receiverRoutine(zmqContext *zmq.Context, sharedMsgBuffer chan string, userList *Users) {
	var rcvMessageCount int

	// User management Socket (REQ/REP)
	usrManagementSoc, _ := zmqContext.NewSocket(zmq.REP)
	defer usrManagementSoc.Close()
	usrManagementSoc.Bind("tcp://*:5003")

	// Message receiver socket (PULL)
	receiverSoc, _ := zmqContext.NewSocket(zmq.PULL)
	receiverSoc.SetLinger(time.Second)
	defer receiverSoc.Close()
	receiverSoc.Bind("tcp://*:5001")

	// Poller for listening on both sockets
	poller := zmq.NewPoller()
	poller.Add(usrManagementSoc, zmq.POLLIN)
	poller.Add(receiverSoc, zmq.POLLIN)

	for {
		sockets, _ := poller.Poll(-1)
		fmt.Printf("\n Received a Message")
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case usrManagementSoc:
				message, err := s.Recv(0)
				if err == nil {
					response := handleUserManagementMessage(message, userList)
					// Send the response
					s.Send(response, 0)
				}
			case receiverSoc:
				fmt.Printf("\n Received a Message from user")
				msg, errRcv := s.Recv(0)
				if errRcv == nil {
					// Push message to the buffer
					sharedMsgBuffer <- msg
					rcvMessageCount++
					if rcvMessageCount%snapshotLength == 0 {
						fmt.Printf("\nReceive count: %d", rcvMessageCount)
					}
				}
			}
		}
	}
}

func main() {
	userList := Users{container: make(map[string]bool)}
	var wg sync.WaitGroup

	sharedMsgBuffer := make(chan string, sharedMessageBufferSize)

	zmqContext, _ := zmq.NewContext()
	defer zmqContext.Term()

	// Start the receiver routine
	wg.Add(1)
	go func(zmqContext *zmq.Context, sharedMsgBuffer chan string, userList *Users) {
		defer wg.Done()
		receiverRoutine(zmqContext, sharedMsgBuffer, userList)
	}(zmqContext, sharedMsgBuffer, &userList)

	// Start the sender routine
	wg.Add(1)
	go func(zmqContext *zmq.Context, sharedMsgBuffer chan string) {
		defer wg.Done()
		senderRoutine(zmqContext, sharedMsgBuffer)
	}(zmqContext, sharedMsgBuffer)

	wg.Wait()
}
