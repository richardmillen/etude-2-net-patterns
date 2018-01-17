# Networking Patterns - Go Example Apps

This folder contains example apps that demonstrate the patterns explored in this study.

## (Random) Word Pub-Sub

+ `word-pub` *(publisher)*
+ `word-lvc` *(proxy)*
+ `word-sub` *(subscriber)*

The basic idea is that a Publisher sends out a bunch of random words one after the other, but they're in a few different languages (English, French & Spanish). A Subscriber can subscribe to one of those languages. In Pub-Sub parlance, these languages are called topics.

If the Proxy is used then Subscribers should subscribe to that instead of the Publisher directly. The Proxy maintains a cache of the last value it receives for each topic (language) which is immediately sent to new Subscribers, rather than them having to wait for a new message on their chosen topic.

### To Do

1. publisher should validate subscriber id's before adding them to the map.
1. subscriber api enables multiple subscriptions which isn't properly supported.

## Logging / Tracing Example

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



