package subscriber

import (
	"chatapp/config"
	"chatapp/user/globals"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"strings"
	"sync"
)

func StartMessageReciever(zmqContext *zmq.Context, selfName string, wg *sync.WaitGroup) {
	defer wg.Done()

	subscriber := createSubscriberSocket(zmqContext, selfName)

	for {
		if globals.GlobalShutdown {
			subscriber.Close()
			fmt.Println("Shutting down Subscriber")
			return
		}
		message, err := subscriber.Recv(0)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}
		handleRecievedMessage(message, selfName)
	}

}

func createSubscriberSocket(zmqContext *zmq.Context, selfName string) *zmq.Socket {
	subscriber, err := zmqContext.NewSocket(zmq.SUB)
	if err != nil {
		panic(err)
	}

	err = subscriber.Connect("tcp://" + config.ServerIP + ":5002")
	if err != nil {
		panic(err)
	}

	subscriber.SetSubscribe("@" + selfName + "^")
	subscriber.SetSubscribe("@all^")

	return subscriber
}

func handleRecievedMessage(message string, selfName string) {
	messageParts := strings.SplitN(message, " ", 3)
	if len(messageParts) < 3 {
		return
	}
	_, senderName, extractedMessage := messageParts[0], messageParts[1], messageParts[2]
	if senderName != selfName {
		fmt.Printf("\nFrom %s: %s\n", senderName, extractedMessage)
	}
}
