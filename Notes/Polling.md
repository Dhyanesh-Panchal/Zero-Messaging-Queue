# Polling

Polling is the mechanism used to read from multiple endpoints at same time.


Refer [Docs](https://pkg.go.dev/github.com/pebbe/zmq4#NewPoller)

```go

poller := zmq.NewPoller()
poller.Add(socket0, zmq.POLLIN)
poller.Add(socket1, zmq.POLLIN)
//  Process messages from both sockets
for {
    sockets, _ := poller.Poll(-1)
    for _, socket := range sockets {
        switch s := socket.Socket; s {
        case socket0:
            msg, _ := s.Recv(0)
            //  Process msg
        case socket1:
            msg, _ := s.Recv(0)
            //  Process msg
        }
    }
}

```