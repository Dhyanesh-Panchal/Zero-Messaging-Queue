package main

import (
	"chatapp/centralserver/globals"
	"chatapp/centralserver/publisher"
	receiver "chatapp/centralserver/reciever"
	"chatapp/centralserver/utils"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

func main() {
	userList := globals.Users{Container: make(map[string]bool)}
	sharedBuffer := utils.CreateMessageBuffer()
	var wg sync.WaitGroup

	wg.Add(1)
	go globals.HandleGlobalShutdown(&wg)

	zmqContext, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go receiver.StartReceiver(zmqContext, sharedBuffer, &userList, &wg)

	wg.Add(1)
	go publisher.StartPublisher(zmqContext, sharedBuffer, &wg)

	// Wait until Global shutdown
	for !globals.GlobalShutdown {
	}
	fmt.Println("termination called!")

	err = zmqContext.Term()
	if err != nil {
		panic(err)
	}

	fmt.Println("zmq terminated, waiting for goroutines to exit.")

	// wait for all Routines to return.
	wg.Wait()
}
