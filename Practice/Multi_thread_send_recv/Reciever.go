package main

import (
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()
	zctx.SetIoThreads(3)
	// Context level Max size is not enforced.
	max_msg_size, _ := zctx.GetMaxMsgsz()
	fmt.Printf("\nMax Message size %d\n", max_msg_size)
	s, _ := zctx.NewSocket(zmq.PULL)

	// Socket Level Maxsize is Enforced
	s.SetMaxmsgsize(1 * 1024 * 1024)
	// s.SetRcvhwm(10)
	// s.SetRcvbuf(10)

	defer zctx.Term()
	defer s.Close()
	s.Bind("tcp://*:5555")

	for {
		// Open a file to append the Data.
		file, _ := os.OpenFile("./Responce.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// fmt.Printf("started recieving\n")
		reply, err := s.Recv(0)

		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("%s\n", reply[:11])
			fmt.Printf("len: %d\n", len(reply))
		}
		_, file_write_error := file.WriteString(reply)
		if err != nil {
			fmt.Print(file_write_error)
		}
		file.Close()
	}
}
