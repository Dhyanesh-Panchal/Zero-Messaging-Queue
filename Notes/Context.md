# Context in ZMQ

ZeroMQ applications always start by creating a context, and then using that for creating sockets.

In Go, it’s the `zmq.NewContext()` call. You should create and use exactly one context in your process.

The context is the container for all sockets in a single process and acts as the transport for *inproc* sockets, which are the fastest way to connect threads in one process. If at runtime a process has two contexts, these are like separate ZeroMQ instances, which is OK to have, but not advised.

Internally, it manages I/O threads and socket state.

```go
ctx, _ := zmq.NewContext()
defer ctx.Term()
```
