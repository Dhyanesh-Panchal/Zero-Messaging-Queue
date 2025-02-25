package sender

import (
	"chatapp/config"
	"chatapp/user/utils"
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

func GetUserName(zmqContext *zmq.Context) string {

	nameVerifier, _ := zmqContext.NewSocket(zmq.REQ)
	defer nameVerifier.Close()
	nameVerifier.Connect("tcp://" + config.ServerIP + ":5003")

	var username string
	for {
		fmt.Print("Enter your name: ")
		username = utils.GetTerminalInput()

		nameVerifier.Send(username, 0)
		response, _ := nameVerifier.Recv(0)

		if response == "OK" {
			break
		}
		fmt.Println("This username is already taken. Please enter a unique name.")
	}

	nameVerifier.Close()
	return username
}
