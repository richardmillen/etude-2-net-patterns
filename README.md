# Étude 2 - Networking Patterns

A study of networking patterns over TCP & UDP, with implementations in Go & *(eventually)* C++.

## Examples

The examples start small & simple, increasing in size & complexity. They aren't intended to be useful applications, but explore one or more aspects of distributed computing.

Each one has an associated a readme which is kind of a mini spec, providing high-level information and low-level protocol details.

1. [Echo Server](docs/EchoServer.md)
1. [Hello World Server](docs/HelloWorldServer.md)
1. [State Server](docs/StateServer.md)
1. [Arithmetic Server](docs/ArithmeticServer.md)
1. [Timeouts](docs/Timeouts.md)

## Demos

The demo folder contains reasonably complete apps & simulations. Each one comes with its own readme much like the examples.

## Notes on Versions

### v0.1

*n.b. Having had a degree of success with my initial investigations, I didn't like the direction it was
heading. I needed to make substantive changes to the supporting packages, so to avoid the overhead of working through the technical debt, I opted to start over.*

### v0.2

For this version I've elected to go with a more state-centric design for a number of reasons:

+ self-documenting protocol; define states that accept defined inputs etc.
+ it should help promote a more declarative style
+ various strategies e.g. heartbeating, beacons, ping-pong and so forth can be implemented as concurrent states in a hierarchical state machine. this should help with composability.
+ simplifies how invalid messages are handled.

This approach does come with some inherent complexities / liabilities including:

+ ensuring each automaton remains isolated. for instance, a server will have at least one connection for each remote client and each of those connections may be in different state.
+ input events will be coming from two places; the network and the api consumer. should all input require an associated input event to be defined against the state? if a server receives a valid request from a client, should the servers response also be known to the current state whether it requires validation or not?
+ network traffic will often be validated, but under what circumstances should the input from the API consumer code be validated?
+ what happens when a state machine reaches it's defined *final state* ? close the connection? do nothing? return to the *initial state*? to what extent should this behaviour be configurable?
+ internal vs. external events? to what extent should they be differentiated? does `myState.HandleSomeEvent` represent an internal event while `myService.HandleSomeEvent` represents an external one? should they be defined different on the state i.e. `fsm.State{InternalEvents: ..., ExternalEvents: ...}` ?

#### fsm Package

Contains the core finite-state machine types e.g.

+ `State` struct: represents a single state within an automaton.
+ `Input` interface: all inputs into a state machine would implement this interface.





#### netx Package

Contains all the networking types. Built upon types in the `fsm` package.

+ `Connector` interface: establishes connections between endpoints.
+ `Conn` struct: built on top of `net.Conn`. *flyweight* State Machine.
+ `Listener` struct: accepts and creates `netx.Conn`s to remote endpoints. implements `Connector` interface.
+ `Dialer` struct: establishes a `netx.Conn` with remote Service. implements `Connector` interface.
+ `Service` struct: main object through which the consumer interacts. likely to be composed of and/or implement one or more types from the `fsm` package.





## Notes on Architectures & Patterns

### Pub-Sub

> Publisher sends a stream of messages.  
> Subscriber receives messages related to one, or more topics.  

+ REQ-REP subscriber sync
+ Topic envelope / frame
+ Last-value caching (LVC)
+ high water mark (HWM)
+ timestamps; abort threshold / latency

### Service Discovery

#### UDP Surveys

+ message
  - ...
  - ...

*n.b. what about local service discovery?*

### Service Presence

#### UDP Beacons

+ message
  - header: "FOO1"
  - body: service TCP port

### Broker

+ message:
  - proto-sig
  - service name
  - logical address
  - properties
  - correlation id
  - frames / body
+ load balance (queue)
+ interconnect; peering
+ failover; clone
+ no workers available?
  - ignore client requests
  - notify requesting client
  - forward to another broker
+ error code response
+ high water mark
+ timestamps; abort threshold
+ assertions
+ transport bridging
+ heartbeats *(broker-to-worker & worker-to-broker)*
  - liveness
  - exponential backoff *(to a max)*
+ ping-pong *(client-to-server)*
  - TTL 6 secs
+ support for downtime
  - upgrades
  - server crash
+ management
  - name service
  - service discovery
+ retries
+ environments; dev, test, prod etc.
+ async; batch send

*n.b. can interconnect & failover be the same i.e. peer broker becomes clone during downtime?*

### Disconnected / Offline Service

> Receives client requests meant for another service via the Broker.   
> Acts on the clients behalf, calling the service at the appropriate time.  
> Provides results to client upon request.

### Parallel Pipeline

+ ventilator
+ workers
+ sink
+ vent == broker *(?)*
+ sink == client *(akin to GFS chunk servers)*

#### Pipelining

+ credit-based *(async)* flow control
+ e.g. large files
+ compression
+ encryption
+ interruption; resume after disconnect

### P2P / Decentralised

+ UDP
  - beacons *(as above)*
  - multicast to logical group
  - parallel over WiFi AP *(TCP - serial/blocking; at bit rate of slowest/furthest device)*
+ silence -> TCP heartbeats -> DEAD
+ *recovery* channel
+ commands:
  - `HELLO`
    - list of groups.
    - list of other services e.g. logging, file transfer etc.
  - `JOIN` / `LEAVE`
+ direct point-to-point messaging
+ groups
  - track nodes *(join/leave)*
+ replication *(snapshots)* - Pub-Sub:
  1. subscribe
  1. recv 1 upd.
  1. queue updates
  1. request snapshot
  1. recv snapshot
  1. apply queued updates
+ subtrees *(/path/to/res)*
+ peer state
  - change counter; rolling 1 byte buffer
+ *mediator*; elected

#### Problems

+ peer discovery
+ interop w/ existing networks
+ data privacy
+ data integrity *(esp. over WAN, WiFi)*
+ logging & monitoring
+ large scale testing & simulation
+ group messaging
+ wide-area bridging
+ configuration

## Origins *(the zguide)*

Much of this study is inspired by the ZeroMQ '[zguide](http://zguide.zeromq.org/page:all)'.

Here's a list of the more relevant sections from the guide containing code samples:

1. [Ask and Ye Shall Receive](http://zguide.zeromq.org/page:all#Ask-and-Ye-Shall-Receive)
1. [Getting the Message Out](http://zguide.zeromq.org/page:all#Getting-the-Message-Out)
1. [Divide and Conquer](http://zguide.zeromq.org/page:all#Divide-and-Conquer)
1. [Handling Multiple Sockets](http://zguide.zeromq.org/page:all#Handling-Multiple-Sockets)
1. [Shared Queue (DEALER and ROUTER sockets)](http://zguide.zeromq.org/page:all#Shared-Queue-DEALER-and-ROUTER-sockets)
1. [ZeroMQ's Built-In Proxy Function](http://zguide.zeromq.org/page:all#ZeroMQ-s-Built-In-Proxy-Function)
1. [Transport Bridging](http://zguide.zeromq.org/page:all#Transport-Bridging)
1. [Handling Errors and ETERM](http://zguide.zeromq.org/page:all#Handling-Errors-and-ETERM)
1. [Handling Interrupt Signals](http://zguide.zeromq.org/page:all#Handling-Interrupt-Signals)
1. [Multithreading with ZeroMQ](http://zguide.zeromq.org/page:all#Multithreading-with-ZeroMQ)
1. [Signaling Between Threads](http://zguide.zeromq.org/page:all#Signaling-Between-Threads-PAIR-Sockets)
1. [Node Coordination](http://zguide.zeromq.org/page:all#Node-Coordination)
1. [Zero-Copy](http://zguide.zeromq.org/page:all#Zero-Copy)
1. [Pub-Sub Message Envelopes](http://zguide.zeromq.org/page:all#Pub-Sub-Message-Envelopes)
1. [Identities and Addresses](http://zguide.zeromq.org/page:all#Identities-and-Addresses)
1. [ROUTER Broker and REQ Workers](http://zguide.zeromq.org/page:all#ROUTER-Broker-and-REQ-Workers)
1. [ROUTER Broker and DEALER Workers](http://zguide.zeromq.org/page:all#ROUTER-Broker-and-DEALER-Workers)
1. [A Load Balancing Message Broker](http://zguide.zeromq.org/page:all#A-Load-Balancing-Message-Broker)
1. [The Asynchronous Client/Server Pattern](http://zguide.zeromq.org/page:all#The-Asynchronous-Client-Server-Pattern)
1. [Prototyping the State Flow](http://zguide.zeromq.org/page:all#Prototyping-the-State-Flow)
1. [Prototyping the Local and Cloud Flows](http://zguide.zeromq.org/page:all#Prototyping-the-Local-and-Cloud-Flows)
1. [Putting it All Together](http://zguide.zeromq.org/page:all#Putting-it-All-Together)
1. [Client-Side Reliability (Lazy Pirate Pattern)](http://zguide.zeromq.org/page:all#Client-Side-Reliability-Lazy-Pirate-Pattern)
1. [Basic Reliable Queuing (Simple Pirate Pattern)](http://zguide.zeromq.org/page:all#Basic-Reliable-Queuing-Simple-Pirate-Pattern)
1. [Robust Reliable Queuing (Paranoid Pirate Pattern)](http://zguide.zeromq.org/page:all#Robust-Reliable-Queuing-Paranoid-Pirate-Pattern)
1. [Heartbeating for Paranoid Pirate](http://zguide.zeromq.org/page:all#Heartbeating-for-Paranoid-Pirate)
1. [Service-Oriented Reliable Queuing (Majordomo Pattern)](http://zguide.zeromq.org/page:all#Service-Oriented-Reliable-Queuing-Majordomo-Pattern)
1. [Asynchronous Majordomo Pattern](http://zguide.zeromq.org/page:all#Asynchronous-Majordomo-Pattern)
1. [Service Discovery](http://zguide.zeromq.org/page:all#Service-Discovery)
1. [Disconnected Reliability (Titanic Pattern)](http://zguide.zeromq.org/page:all#Disconnected-Reliability-Titanic-Pattern)
1. [Binary Star Implementation](http://zguide.zeromq.org/page:all#Binary-Star-Implementation)
1. [Binary Star Reactor](http://zguide.zeromq.org/page:all#Binary-Star-Reactor)
1. [Model One: Simple Retry and Failover](http://zguide.zeromq.org/page:all#Model-One-Simple-Retry-and-Failover)
1. [Model Two: Brutal Shotgun Massacre](http://zguide.zeromq.org/page:all#Model-Two-Brutal-Shotgun-Massacre)
1. [Model Three: Complex and Nasty](http://zguide.zeromq.org/page:all#Model-Three-Complex-and-Nasty)
1. [Pub-Sub Tracing (Espresso Pattern)](http://zguide.zeromq.org/page:all#Pub-Sub-Tracing-Espresso-Pattern)
1. [Last Value Caching](http://zguide.zeromq.org/page:all#Last-Value-Caching)
1. [Slow Subscriber Detection (Suicidal Snail Pattern)](http://zguide.zeromq.org/page:all#Slow-Subscriber-Detection-Suicidal-Snail-Pattern)
1. [Representing State as Key-Value Pairs](http://zguide.zeromq.org/page:all#Representing-State-as-Key-Value-Pairs)
1. [Getting an Out-of-Band Snapshot](http://zguide.zeromq.org/page:all#Getting-an-Out-of-Band-Snapshot)
1. [Republishing Updates from Clients](http://zguide.zeromq.org/page:all#Republishing-Updates-from-Clients)
1. [Working with Subtrees](http://zguide.zeromq.org/page:all#Working-with-Subtrees)
1. [Ephemeral Values](http://zguide.zeromq.org/page:all#Ephemeral-Values)
1. [Using a Reactor](http://zguide.zeromq.org/page:all#Using-a-Reactor)
1. [Adding the Binary Star Pattern for Reliability](http://zguide.zeromq.org/page:all#Adding-the-Binary-Star-Pattern-for-Reliability)
1. [Building a Multithreaded Stack and API](http://zguide.zeromq.org/page:all#Building-a-Multithreaded-Stack-and-API)
1. [Transferring Files](http://zguide.zeromq.org/page:all#Transferring-Files)
1. [A Self-Healing P2P Network in 30 Seconds](http://zguide.zeromq.org/page:all#A-Self-Healing-P-P-Network-in-Seconds)
1. [Cooperative Discovery Using UDP Broadcasts](http://zguide.zeromq.org/page:all#Cooperative-Discovery-Using-UDP-Broadcasts)
1. [Designing the API](http://zguide.zeromq.org/page:all#Designing-the-API)
1. [Dealing with Blocked Peers](http://zguide.zeromq.org/page:all#Dealing-with-Blocked-Peers)

*I implemented many of the examples above (in C++) as a learning exercise. Although my implementations are 'experimental' in nature and therefore quite *raw*, some of them do resolve issues that I found when trying to run the 'official' examples. Also, each example is self-contained i.e. no external dependencies besides ZeroMQ itself. This should make it easier to understand the code, in contrast with the examples in the zguide which use a collection of opaque helper functions.*

*You can find my zguide examples [here](https://github.com/richardmillen/zguide-examples).*

