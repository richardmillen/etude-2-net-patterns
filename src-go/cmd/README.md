# Networking Patterns - Go Example Apps

This folder contains example apps that demonstrate the patterns explored in this study.

## (Random) Word Pub-Sub

+ `word-pub` *(publisher)*
+ `word-lvc` *(proxy)*
+ `word-sub` *(subscriber)*

The basic idea is that a Publisher sends out a bunch of random words one after the other, but they're in a few different languages (English, French & Spanish). A Subscriber can subscribe to one of those languages. In Pub-Sub parlance, these languages are called topics.

If the Proxy is used then Subscribers should subscribe to that instead of the Publisher directly. The Proxy maintains a cache of the last value it receives for each topic (language) which is immediately sent to new Subscribers, rather than them having to wait for a new message on their chosen topic.

### Known Limitations

1. publisher should validate subscriber id's before adding them to the map esp. in case of overwrites.
1. subscriber api enables multiple subscriptions which isn't properly supported.
1. subscriber api enables multiple topics when only one *(the first)* is actually used.
1. handle connection drops properly. *(attempt reconnect)*
1. any need for heartbeats?

## Translation Server

+ `trans-client`
+ `trans-server`

The client app sends requests to the translation server. The requests contain a word, the language and the target language. The server then responds with one or more translations.

## Translation Server Pool

+ `pool-client`
+ `trans-server`

The client app loops through a pool of translation servers attempting to send requests to each one, one after the other until a live server sends back a response. 

## Broker-based Translation Service

+ `trans-client`
+ `svc-broker`
+ `trans-server`

The client app sends requests to one of *n* translation servers with all communication being routed through a broker/proxy server. The broker is transparent to the client.

n.b. The translation servers pause for *n* seconds to simulate work.

## Cloned, Broker-based Translation Service

+ `trans-client`
+ `clone-broker`
+ `trans-server`

This the same as the previous example, except the broker supports *peering* and *failover*. If a client sends a requiest and broker dies, it will resend to the other broker.

## Broker-based Service Discovery

+ `disco-client`
+ `disco-broker`
+ `trans-server`

The translation service registers itself with the broker. The client then asks the broker if the translation service exists, to which the broker responds directly.

## Name Server

+ `name-client`
+ `name-server`
+ `trans-server`

Instances of the translation server register themselves with one or more name servers in a pool. The name servers stay in sync which each other. The client then asks the pool of name servers for the address(es) of known translation servers. It can then send requests to one or more of the translation servers directly.

## Logging / Tracing

+ `log-server`
+ `log-forwarder` *(store & forward)*
+ `log-client`

use cases:

+ store & forward
+ errors
+ performance
+ postmortems
+ audit trail
+ monitor conjestion

##


