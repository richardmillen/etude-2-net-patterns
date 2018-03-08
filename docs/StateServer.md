# Example 3: State Server

## Server

The 'state' server example #1 sends a sequence of messages to clients immediately
after they connect, each message corresponding to a different state. Once the
final state is reached and the all messages have been sent the connection to the
client is closed.

An 'extra' state can be enabled via a command line flag which causes the server to
send an unknown / unexpected message to the client. The resulting client behaviour
is dependent upon the version used.

## Client 1

Enters a receiving state (opening a connection to a 'state' server),
receives three specific messages (it doesn't care which order) then the server closes
the connection and the client exits. all other messages will be ignored.

## Client 2

Receives three specific messages from a 'state' server.
the messages must be in the correct order, or an error is displayed and
the client aborts.







