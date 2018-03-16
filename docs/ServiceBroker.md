# Example 9: Service Broker

## Broker

this acts a proxy between one or more clients and any available services.
all client requests are made to the broker which then forwards them on
to an appropriate service instance. all messages from the service are then
forwarded on to the relevant client.

## Client

this app makes calls to several services via a 'service-broker' instance.
first it calls the 'rand-service' a user defined number of times. it then
calls the 'sort-service' to sort the random numbers immediately, or calls
the 'offline-service' which calls the 'sort-service' on the clients behalf.
if the latter option is chosen then the client periodically checks the
'offline-service' for status updates, then eventually gets the result.

### Services

This example comes with several *toy* services which are described below.

### Sort Service

this service registers itself with a 'service-broker' instance then receives
requests from clients to sort random number sequences using a few basic algorithms.

### Random Number Service

this service registers itself with an 'service-broker' instance then acts as
a simple random number generation service for clients.

### Offline Service

this service registers itself with a running 'service-broker' then receives
requests that describe actions to be performed later (offline) on behalf of
clients. it returns a unique 'job' identifier which is then used by the client
to get status updates / results.

this approach can be useful where the requested action can take a long time,
or involves a service that isn't always available. it could also be useful for
scheduling tasks.














