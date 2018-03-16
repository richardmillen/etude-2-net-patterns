# Name Lookup Server Example

Clients contact a central service lookup which returns address information for 
the registered instance. The client then contacts the desired service directly.

## Goals

Explore the following problems:

+ Service discovery.
+ Service liveness / availability.

## Implementation

### Overall Behaviour

A service registers itself with a central name lookup server which service consumers
contact prior to sending requests. The central name lookup server returns the service 
address to calling clients.

Supports arbitrary registration groups *(environments)*, where for example, one instance
of an *echo* service might register itself on the *development* environment and another 
on *test* and so forth.

### States

![client and server state diagrams](../images/NameServerStateDiagrams.png)

### Formal Grammar

The following ABNF grammar defines the protocol:

```abnf

```

## Limitations

Only a single instance of a service may register itself or the previously registered 
instance will be overwritten.

## Security

All messages are sent between nodes in plain text.

