# Context in ZMQ

ZeroMQ applications always start by creating a context, and then using that for creating sockets.

In Go, it’s the `zmq.NewContext()` call. You should create and use exactly one context in your process.

The context is the container for all sockets in a single process and acts as the transport for *inproc* sockets, which are the fastest way to connect threads in one process. If at runtime a process has two contexts, these are like separate ZeroMQ instances, which is OK to have, but not advised.

Internally, it manages I/O threads and socket state.

```go
ctx, _ := zmq.NewContext()
defer ctx.Term()
```

### Context is responsible for:

- Managing the underlying I/O threads that handle message passing.
- Allocating and managing sockets.
- Handling global configurations such as thread pool size.


## Configurations possible in a Context

### 1. Number of I/O Threads
Controls the number of background threads for processing messages.
```go
context.SetIoThreads(N)
```

### 2. MaxSockets per Context
Limits the number of sockets that can be opened.
```go
context.SetMaxSockets(N)
```

### 3. Message Limits and Timeouts
You can configure:

#### Enable/Disable IPv6
```go
context.SetIpv6(true) // Enable IPv6
```

#### Thread Scheduling Policy
```go
context.SetThreadSchedPolicy(1) // Custom scheduling policy
```

#### Thread Priority
```go
context.SetThreadPriority(10) // Set priority to 10
```

#### Max Message Size Restriction
```go
context.SetMaxMessageSize(10240) // Limit message size to 10KB
```

