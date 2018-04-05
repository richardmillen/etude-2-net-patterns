# Name Lookup Example

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

*n.b. C: client, S: lookup server, P: peer service.*

```abnf
traffic         = *(registration / lookup)

registration    = P:im-here (S:hi-there / S:come-again / S:rtfm)
                / P:im-off

lookup          = C:do-u-know (S:look-here / S:dont-know / S:rtfm)

;       Service registers itself
im-here         = signature %d1 name group address port
name            = string                    ; Service name / alias
group           = string                    ; Service group
address         = 
port            = 

;       Service unregisters itself
im-off          = signature %d2 

;       Lookup server tells Service it's successfully registered
hi-there        = signature %d3 

;       Lookup server tells Service it's already registered
come-again      = signature %d4 

;       Client asks for address of Service
do-u-know       = signature %d5 

;       Lookup server tells Client where to find Service
look-here       = signature %d6 

;       Lookup server tells Client it doesn't know the requested Service
dont-know       = signature %d7 

;       Server tells Client or Service it sent an invalid message
rtfm            = signature %d8 reason
reason          = string                    ; Printable explanation

;       Protocol signature: "NL" (two octets)
signature       = %x4e %x4c

;       Strings are length-prefixed
string          = number-1 *VCHAR

number-1        = 1OCTET
```

## Limitations

Only a single instance of a service may register itself within a group. 

If the lookup service receives a request to register a service that's already
registered within the specified group, the lookup service will respond with an
`inou` message and the service shall not be registered.

Only a single service may be registered per connection.

## Security

All messages are sent between nodes in plain text.

