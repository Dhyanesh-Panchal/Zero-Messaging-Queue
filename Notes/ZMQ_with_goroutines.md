# ZMQ and Go Routines

## `zmq.socket` is **NOT** Thread Safe.

Hence Passing the sockets through seperate go routine will cause anamolies
Sharing a single socket across goroutines without proper synchronization can lead to race conditions and undefined behavior.


### **BAD CODE**
```go

func myfn(socket *zmq.sokcet){
    // Using socket
}

func main(){
    context, _ := zmt.NewContext()
	socket, _ := context.NewSocket(zmt.PUB) // ANY type of socket

    go myfn(socket) //BAD IDEA
}
```

### Instead Pass the `zmq.Context` accross the Routines, Instanciate and Utilze `sockets` under same routine, & Synchronize execution using `Channels`.

### **GOOD CODE**
```go

func myfn(context *zmq.Context) // Pass channel variables if sync needed
{
    socket,_ := context.NewSocket(zmt.PUB)
    defer socket.Close()
    //utilize and end, exclusivily for this routine.

}

func main() {

	context, _ := zmt.NewContext()
	defer context.Term()

    go myfn(context)
}
```
