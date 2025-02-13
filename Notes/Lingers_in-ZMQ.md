## Linger in ZeroMQ

Linger determines how long a socket will wait to finish sending or receiving queued messages after calling `socket.Close()` or terminating the context.

### Linger Settings
Linger is configured using `socket.SetLinger(ms)`, where `ms` is the timeout in milliseconds:

- `0`: Immediate shutdown (discard pending messages). This is the default behavior.
- `-1`: Block indefinitely until all messages are processed (potential risk of deadlocks).
- `> 0`: Wait for `ms` milliseconds before closing, providing a safer timeout mechanism.

Example:

```go
socket.SetLinger(1000) // Wait 1 second before closing
socket.SetLinger(0)    // Discard messages immediately
```