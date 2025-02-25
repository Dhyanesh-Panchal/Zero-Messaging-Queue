package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

const DefaultServerIP = "10.20.40.165"

var reader = bufio.NewReader(os.Stdin)

func receiveMessage(context *zmq.Context, nameChannel chan string) {
	nameVerifier, _ := context.NewSocket(zmq.REQ)
	defer nameVerifier.Close()
	nameVerifier.Connect(fmt.Sprintf("tcp://%s:5003", DefaultServerIP))

	var username string
	for {
		fmt.Print("Enter your name: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		nameVerifier.Send(username, 0)
		response, _ := nameVerifier.Recv(0)

		if response == "OK" {
			break
		}
		fmt.Println("This username is already taken. Please enter a unique name.")
	}

	nameChannel <- username

	subscriber, _ := context.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect(fmt.Sprintf("tcp://%s:5002", DefaultServerIP))
	subscriber.SetSubscribe("@" + username + "^")
	subscriber.SetSubscribe("@all^")

	for {
		message, _ := subscriber.Recv(0)
		messageParts := strings.SplitN(message, " ", 3)
		if len(messageParts) < 3 {
			continue
		}
		_, senderName, extractedMessage := messageParts[0], messageParts[1], messageParts[2]
		if senderName != username {
			fmt.Printf("\nFrom %s: %s\n", senderName, extractedMessage)
		}
	}
}

func main() {
	serverIP := DefaultServerIP
	if len(os.Args) > 1 {
		inputServerIP := os.Args[1]
		if matched, _ := regexp.MatchString(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`, inputServerIP); matched {
			serverIP = inputServerIP
		} else {
			fmt.Printf("\nInvalid IP address, using default: %s\n", DefaultServerIP)
		}
	}

	context, _ := zmq.NewContext()
	defer context.Term()

	sender, _ := context.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect(fmt.Sprintf("tcp://%s:5001", serverIP))

	nameChannel := make(chan string)
	go receiveMessage(context, nameChannel)

	selfName := <-nameChannel
	for {
		fmt.Print("Enter message: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if strings.HasPrefix(msg, "@") {
			sender.Send(fmt.Sprintf("%s %s", selfName, msg), 0)
		} else {
			fmt.Println("Please mention the recipient with '@' at the beginning.")
		}
	}
}
