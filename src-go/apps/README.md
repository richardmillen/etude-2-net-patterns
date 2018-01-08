# Networking Patterns - Go Apps

This folder contains example apps that demonstrate the patterns explored by this study.

## PUB-SUB

notes:-

+ REQ-REP subscriber sync
+ KEY envelope / frame
+ Last-value caching (LVC)
+ high water mark (HWM)
+ timestamps; abort threshold / latency

### Weather Updates Example

+ `weather-server`
+ `weather-proxy` *(LVC)*
+ `weather-client`

### Logging / Tracing Example

+ store & forward
+ errors
+ performance
+ postmortems
+ audit trail
+ monitor conjestion

---

+ `log-server`
+ `log-forwarder` *(store & forward)*
+ `log-client`

## Broker

notes:-

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
  - exponential backoff *(to max)*
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

## Parallel Pipeline

+ ventilator
+ workers
+ sink
+ vent == broker *(?)*
+ sink == client *(akin to GFS chunk servers)*

### Pipelining

+ credit-based *(async)* flow control
+ e.g. large files
+ compression
+ encryption
+ interruption; resume after disconnect

## Service Presence

### UDP beacons

+ message
  - header: "FOO1"
  - body: service TCP port

## P2P / Decentralised












