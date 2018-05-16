# Name Lookup Example

Consumers obtain service endpoint addresses from a central service lookup. The
consumer then contacts the service endpoint(s) directly.

## Goals

Explore the following problems:

+ Service discovery.
+ Service liveness / availability.

## Implementation

### Overall Behaviour

A service registers itself with a central name lookup server.

A service consumer sends requests to the lookup service, to which the service
responds with a list of all matching service endpoints.

The consumer then communicates directly with one or more of the service endpoints.

Supports arbitrary registration groups *(environments)*, where for example, one 
instance of an *echo* service might register itself on the *development* environment 
and another on *test* and so forth.

**IMPORTANT**: If a lookup service endpoint crashes, all registration entries 
will be lost.

### PUSH Notifications

This may be a good use case for a PUSH notification style architecture, where
the lookup service sends realtime updates to listening service consumers. The
approach has a lot going for it, such as:

+ reduces responsibility of the service consumer.
+ the lookup service could send updates to consumers only when required and in realtime.
+ update (PUSH) messages could be smaller than the initial lookup i.e. returning difference data<sup>1</sup>.
+ the consumer could monitor the liveness of the lookup service.

Despite this, I've opted for a request/reply model due to its relative simplicity.

<sup>1. This isn't strictly intrinsic to PUSH notifications, but it seems more natural than in a request/reply model.</sup>

### Service Identification

What data may be used to identify a service?

+ Service name e.g. *echo*
+ Service group e.g. *development*
+ Endpoint address e.g. *192.168.1.1*
+ Endpoint local/ephemeral port e.g. *61000*

### Lookup & QoS

One of the most obvious drawbacks to service consumers is the potential delay 
while waiting for a lookup request to complete. 

Lookup service instances have no knowledge of one another, nor do they maintain
a connection to consumer endpoints. This design has the following consequences:

+ lookup requests must be made by the consumer.
+ consumers will be unaware that a lookup service endpoint has become unavailable until their next lookup attempt.
+ consumers will be configured with one or more lookup service addresses.

If a consumer needs to track service availability e.g. in order to notify a user
when a service becomes unavailable, it should perform regular & frequent lookups. 
If this isn't required, it can perform lookups as close as possible/reasonable to
its next anticipated interaction with the service so that it has up-to-date 
service endpoint information and minimal interaction with the lookup service.

### Registration

> What are the pros & cons of maintaining connections between a lookup service 
> endpoint and its registered service endpoints?

Pros:

1. Service endpoints would know if the lookup service goes offline and take appropriate action.
1. Would allow the *liveness* of service endpoints to be tracked and therefore prioritised.
1. No need for services to unregister; a dropped connection results in removal of registration.
1. Would allow the lookup service to revoke/unregister services and notify them.

Cons:

1. A connection per lookup service equates to more work/load on the service endpoint.

> What are the pros & cons of services registering on startup, disconnecting then 
> unregistering when they shut down?

Pros:

1. Greatly simplifies the lookup service (no liveness tracking).
1. Reduces the load on the lookup service (no liveness tracking).
1. It may be left to the consumer & service to track liveness.

Cons:

1. The lookup service could be shut down or crash (losing all registration data) and then be restarted without service endpoints ever knowing that they're no longer registered.<sup>1</sup>
1. It may not be possible for the lookup service to send notifications to service endpoints e.g. "you've been unregistered".

<sup>1. This can't be resolved by caching to disk because a service endpoint could shutdown while the lookup service is offline. The lookup service would then have to ping all cached service endpoints when it starts.</sup>

> What are the pros & cons of caching registration data so that it's not lost due to a restart / crash?

Pros:

> If services are required to unregister themselves, under what circumstances could a service be registered multiple times?

1. Instance *A* crashes before unregistering. A process/person starts instance *B*.
1. A bug causes instance *A* to shut down without unregistering. A process/person starts instance *B*.
1. A bug causes instance *A* to hang indefinitely before unregistering. A process/person starts instance *B*.
1. A process/person unintentionally starts instances *A* and *B*.
1. A process/person intentionally (but incorrectly) starts instances *A* and *B*.

> How might the lookup service behave when it receives a registration request for a service that's already registered?

1. Do nothing.
1. Return an error code to the registering service.
1. 
1. 

> TODO: Consider the pros & cons of having vs. not having a universally unique identifier per service instance. Also consider its origin; i.e. whether it's created (and perhaps returned) by the lookup server vs. created by the registering service and passed to the lookup service.



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

If a lookup service endpoint crashes, all registration data will be lost and
service endpoints must re-register.

The design only naturally supports a single lookup service. However, services 
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

