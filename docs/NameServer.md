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

### Registration Collisions

Under what circumstances could a service be registered multiple times?

+ Instance *A* crashes before unregistering. A process/person starts instance *B*.
+ A bug causes instance *A* to shut down without unregistering. A process/person starts instance *B*.
+ A bug causes instance *A* to hang indefinitely before unregistering. A process/person starts instance *B*.
+ A process/person unintentionally starts instances *A* and *B*.
+ A process/person intentionally (but incorrectly) starts instances *A* and *B*.

What data uniquely identifies a service?

+ Service name e.g. *echo*
+ Service group e.g. *development*
+ End point address e.g. *192.168.1.1*
+ End point port e.g. *5678*

> TODO: Consider the pros & cons of having no unique per-instance identifier.

> TODO: Consider the pros & cons of a unique per-instance identifier provided by a) the lookup service and b) the registering service.

### States

> TODO: add state diagrams.

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
address         = ipv4-addr
port            = number-2

;       Lookup server tells Service it's successfully registered
hi-there        = signature %d3 sid
sid             = number-1                  ; Service registration ID

;       Lookup server tells Service it's already registered
come-again      = signature %d4 

;       Server tells Client or Service it sent an invalid message
rtfm            = signature %d8 reason
reason          = string                    ; Printable explanation

;       Service unregisters itself
im-off          = signature %d2 

;       Client asks for address of Service
do-u-know       = signature %d5 

;       Lookup server tells Client where to find Service
look-here       = signature %d6 

;       Lookup server tells Client it doesn't know the requested Service
dont-know       = signature %d7 

;       Protocol signature: "NL" (two octets)
signature       = %x4e %x4c

;       Strings are length-prefixed
string          = number-1 *VCHAR

;       IP v4 address
ipv4-addr       = number-1 number-1 number-1 number-1

number-1        = 1OCTET
number-2        = 2OCTET
```

## Limitations

The design naturally supports only a single lookup service. However, services 
and consumers could conceivably be configured to use more than one.

Only a single instance of a service may register itself within a group. 

If the lookup service receives a request to register a service that's already
registered within the specified group, the lookup service will respond with an
`come-again` message and the service shall not be registered.

Only a single service may be registered per connection.

Service registrations don't expire so there's no way to tell consumers to switch 
to another service instance short of turning the existing one off.

## Security

All messages are sent between nodes in plain text.

