package main

import (
	"chatapp/user/globals"
	"chatapp/user/sender"
	"chatapp/user/subscriber"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go globals.HandleGlobalShutdown(&wg)

	zmqContext, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}

	username := sender.GetUserName(zmqContext)

	wg.Add(1)
	go subscriber.StartMessageReciever(zmqContext, username, &wg)

	wg.Add(1)
	go sender.StartSender(zmqContext, username, &wg)

	// Wait until Global shutdown
	for !globals.GlobalShutdown {
	}
	fmt.Println("termination called!")

	// call context termination
	zmqContext.Term()

	// wait for all routines to return.
	wg.Wait()

}
