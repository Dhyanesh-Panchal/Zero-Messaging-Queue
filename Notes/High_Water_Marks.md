# High Water Marks (HWM) in ZeroMQ

## Problem:
When messages are sent rapidly between processes, memory can quickly become a bottleneck. A short processing delay in one process can cause a backlog of messages, potentially overwhelming the system. 

For example, if process A is sending messages at high frequency to process B, and B becomes temporarily busy, messages will accumulate in various places:
- In B’s network buffers
- On the Ethernet wire
- In A’s network buffers
- In A’s memory, if there’s no flow control

## Solution: High-Water Mark (HWM)
- ZeroMQ introduces **HWM (High-Water Mark)** to define the capacity of its internal pipes.
- Each connection into or out of a socket has its own pipe, with separate **send** and/or **receive** HWMs, depending on the socket type.
- Different socket types manage their buffers differently:
  - **PUB, PUSH**: Only have **send** buffers.
  - **SUB, PULL, REQ, REP**: Only have **receive** buffers.
  - **DEALER, ROUTER, PAIR**: Have both **send** and **receive** buffers.

### Behavior When HWM is Reached
- If a socket reaches its HWM, it will **either block or drop messages**, depending on its type:
  - **PUB and ROUTER** sockets will **drop** messages.
  - Other socket types will **block** the sender until space is available.
- For **inproc** transport, both sender and receiver share the same buffer, so the actual HWM is the sum of the values set by both ends.

### Example in Go
To set the High-Water Mark in ZeroMQ:

```go
socket, _ := ctx.NewSocket(zmq.PUB)
socket.SetSndhwm(1000) // Set send high-water mark
socket.SetRcvhwm(1000) // Set receive high-water mark
```
By carefully tuning the HWM values, we can prevent excessive memory usage while maintaining efficient message flow.
