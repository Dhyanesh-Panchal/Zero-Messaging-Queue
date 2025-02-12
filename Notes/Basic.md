# ZeroMQ
ZMQ (ZeroMQ) is a high-performance asynchronous messaging library that provides efficient message queuing and communication between distributed applications. 

-  Unlike traditional message brokers (such as RabbitMQ or Kafka), ZMQ operates as a lightweight messaging layer that can be embedded directly into applications, eliminating the need for an external message broker.


It supports multiple messaging patterns, including:

Request-Reply
Publish-Subscribe
Push-Pull
Exclusive Pair
Multicast

## Working

ZMQ provides an abstraction over raw sockets, handling message queuing, connection management, and asynchronous communication internally. It uses socket-like objects but operates at a much higher level.


High-Level Workflow
Create a ZMQ context – A shared environment for sockets.
Create and bind/connect sockets – Choose appropriate socket types.
Send and receive messages – Messages are queued and processed asynchronously.
Handle errors and retries – ZMQ provides built-in reconnect and failover mechanisms.


Socket Types and Communication Patterns
ZMQ provides different socket types, which define the communication patterns between processes or nodes.

Socket Type	Description
ZMQ_REQ	Request socket (sends requests, expects replies).
ZMQ_REP	Reply socket (receives requests, sends replies).
ZMQ_PUB	Publisher socket (sends messages in Pub-Sub pattern).
ZMQ_SUB	Subscriber socket (receives messages in Pub-Sub pattern).
ZMQ_PUSH	Push socket (sends messages in a pipeline pattern).
ZMQ_PULL	Pull socket (receives messages in a pipeline pattern).
ZMQ_PAIR	Exclusive 1-to-1 communication.