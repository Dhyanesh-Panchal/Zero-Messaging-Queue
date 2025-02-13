# Sockets in ZMQ

Sockets in ZeroMQ follow a lifecycle similar to BSD sockets, consisting of four main stages:

1. **Creating and Destroying Sockets**
   - Sockets are created and closed, forming a structured lifecycle.
   - In Go, this is done using:
   
   ```go
   socket, _ := ctx.NewSocket(zmq.REP)
   defer socket.Close()
   ```

2. **Configuring Sockets**
   - Options can be set and checked as needed.
   - Use `SetSockOpt()` and `GetSockOpt()` for configuration.

3. **Plugging Sockets into the Network Topology**
   - Establish connections using `socket.Bind()` and `socket.Connect()`.
   - A node calling `Bind()` acts as a server, and a node calling `Connect()` acts as a client.
   
   ```go
   socket.Bind("tcp://*:5555")
   socket.Connect("tcp://localhost:5555")
   ```

   - ZeroMQ connections are different from classic TCP connections:
     - They support multiple transports (`inproc`, `ipc`, `tcp`, `pgm`, `epgm`).
     - One socket may have many outgoing and incoming connections.
     - No `zmq_accept()` method is needed; binding a socket automatically starts accepting connections.
     - Connections occur in the background, with automatic reconnection.
     - A server node can bind to multiple endpoints using a single socket:
       
       ```go
       socket.Bind("tcp://*:5555")
       socket.Bind("tcp://*:9999")
       socket.Bind("inproc://somename")
       ```

4. **Sending and Receiving Messages**
   - Messages are sent and received using `Send()` and `Recv()` methods.
   - ZeroMQ sockets operate with discrete messages rather than byte streams.
   
   ```go
   socket.Send("Hello", 0)
   msg, _ := socket.Recv(0)
   ```
   
   - Additional characteristics:
     - ZeroMQ messages are **length-specified binary data**.
     - Sockets handle I/O in a background thread, with messages queued for asynchronous sending.
     - Built-in one-to-N routing behavior based on socket type.

## Unicast Transports
- **tcp**: Commonly used for network transport, supporting dynamic endpoint connections.
- **ipc**: Used for inter-process communication, typically with `.ipc` endpoint names.
- **inproc**: Fastest transport for thread-to-thread communication within a process.
---
# Summary on All Types of sockets : [Refer](https://zeromq.org/socket-api/)
