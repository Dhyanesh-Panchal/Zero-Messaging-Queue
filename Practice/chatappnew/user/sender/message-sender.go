package sender

import (
	"chatapp/config"
	"chatapp/user/globals"
	"chatapp/user/utils"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"strings"
	"sync"
)

func StartSender(zmqContext *zmq.Context, selfName string, wg *sync.WaitGroup) {
	defer wg.Done()

	sender := createSenderSocket(zmqContext)

	for {
		if globals.GlobalShutdown {
			fmt.Println("Shutting down Sender")
			sender.Close()
			return
		}
		fmt.Print("Enter message: ")
		msg := utils.GetTerminalInput()

		if strings.HasPrefix(msg, "@") {
			_, err := sender.Send(selfName+" "+msg, 0)

			if err != nil {
				fmt.Println("Error sending the message", err)
			}
		} else {

			fmt.Println("Please mention the recipient with '@' at the beginning.")

		}
	}

}

func createSenderSocket(zmqContext *zmq.Context) *zmq.Socket {

	sender, err := zmqContext.NewSocket(zmq.PUSH)

	if err != nil {
		panic(err)
	}

	err = sender.Connect("tcp://" + config.ServerIP + ":5001")

	if err != nil {
		panic(err)
	}

	return sender
}
