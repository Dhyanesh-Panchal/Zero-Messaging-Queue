package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

const CHUNK_SIZE = 1024 * 1024

func send_large(socket *zmq.Socket, data []byte) {
	msg_len := len(data)

	total_chunks := int(math.Ceil(float64(msg_len) / CHUNK_SIZE))
	fmt.Printf("\nTotal Chunks %d", total_chunks)
	// Create Chunk and send with Send more flag

}

func main() {
	zctx, _ := zmq.NewContext()
	zctx.SetIoThreads(3)
	zctx.SetMaxMsgsz(1024 * 1024)

	s, _ := zctx.NewSocket(zmq.PUSH)
	// s2, _ := zctx.NewSocket(zmq.PUSH)
	// s.SetSndhwm(10)
	// s.SetSndbuf(10)
	defer s.Close()
	// defer s2.Close()
	s.Connect("tcp://10.20.40.165:5555")

	time.Sleep(3 * time.Second)
	n_msgs, _ := strconv.Atoi(os.Args[1])
	big_msg := make([]byte, 2*1024*1024)

	for i := 0; i < n_msgs; i++ {
		fmt.Printf("\nSending #%d", i)
		message := fmt.Sprintf("#%d, Message %s", i, big_msg)

		// send_large(s, []byte(message))
		_, err := s.Send(message, 0)
		if err != nil {
			fmt.Print(err)
		}
		time.Sleep(time.Millisecond * 10)
	}
}
