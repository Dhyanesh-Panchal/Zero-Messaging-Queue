package publisher

import (
	"chatapp/centralserver/globals"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"strings"
	"sync"
)

func processMessage(message string) string {
	msgSplit := strings.Split(message, " ")
	senderName, topic, msgContent := msgSplit[0], msgSplit[1], strings.Join(msgSplit[2:], " ")
	publishMessage := topic + "^ " + senderName + " " + msgContent
	return publishMessage
}

func StartPublisher(zmqContext *zmq.Context, sharedMsgBuffer chan string, wg *sync.WaitGroup) {

	publisher, _ := zmqContext.NewSocket(zmq.PUB)
	publisher.Bind("tcp://*:5002")

	for {

		if globals.GlobalShutdown {

			publisher.Close()

			wg.Done()
			fmt.Println("Returning for publisher")
			return
		}

		message := <-sharedMsgBuffer

		if message == "" {
			continue
		}

		publishMessage := processMessage(message)

		_, errSnd := publisher.Send(publishMessage, 0)
		if errSnd != nil {
			fmt.Println(errSnd)
		}
	}
}
