package receiver

import (
	"fmt"
	"sync"

	"chatapp/centralserver/globals"
	zmq "github.com/pebbe/zmq4"
)

// Handles user registration messages
func handleUserManagementMessage(message string, userList *globals.Users) string {
	userName := message
	userList.RWlock.Lock()
	defer userList.RWlock.Unlock()
	if !userList.Container[userName] {
		userList.Container[userName] = true
		fmt.Printf("\n\tNEW USER JOINED: %s", userName)
		return "OK"
	}
	return "USER_ALREADY_EXIST"
}

// Receiver Routine (Handles REP & PULL sockets)
func StartReceiver(zmqContext *zmq.Context, sharedMsgBuffer chan string, userList *globals.Users, wg *sync.WaitGroup) {
	// User management (REP)
	usrManagementSoc, _ := zmqContext.NewSocket(zmq.REP)
	usrManagementSoc.Bind("tcp://*:5003")

	// Message receiver (PULL)
	receiverSoc, _ := zmqContext.NewSocket(zmq.PULL)
	receiverSoc.Bind("tcp://*:5001")

	// Poller for both sockets
	poller := zmq.NewPoller()
	poller.Add(usrManagementSoc, zmq.POLLIN)
	poller.Add(receiverSoc, zmq.POLLIN)

	for {

		if globals.GlobalShutdown {
			receiverSoc.Close()
			usrManagementSoc.Close()

			close(sharedMsgBuffer)

			wg.Done()
			return
		}

		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case usrManagementSoc:
				message, err := s.Recv(0)
				if err != nil {
					break
				}
				response := handleUserManagementMessage(message, userList)
				s.Send(response, 0)

			case receiverSoc:
				msg, errRcv := s.Recv(0)
				if errRcv != nil {
					break
				} else {
					sharedMsgBuffer <- msg
				}
			}
		}
	}
}
