# Example 9: Pipelining

## Client

Requests a large file from the server using an asynchronous credit-based
flow control to throttle the volume data being sent to the client.

## Server

Receives requests from client and to which it sends a large (fake) file.
The client must instruct the server that it has capacity before the server
will send anything. This happens asynchronously.







