# Clone Server Example

the clone server is a simple echo server, but with some important differences.

this app can be started with the address of another server instance from which
it will obtain a list of active instances before it starts listening for client
requests. then whenever a client connects it will be given a list of known good
server addresses.

## clone server

acts as a simple echo server to clients, returning a list of other active servers
to clients when a new connection is established.

note that a clone server must either be acting as 'primary' or 'secondary' and
cannot be both. this limitation makes the example a less realistic implementation,
but reduces complexity.

## primary server

acts as a clone server.

maintains a list of active servers which are sent to 'secondary' servers upon
request. started by running a clone server without the '--primary-server' flag.

if a primary server disappears then the latest active server list is traversed
from top to bottom until a server is able to assume its new role as primary.

## secondary server

acts as a clone server.

sends it's own clone server address to the primary server specified using the
'--primary-server' flag. receives a list of known good clone servers (possibly
including itself) from the primary server.


