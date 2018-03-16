# Pipelining Example

## Goals

Explore the *pipelining* pattern; an efficient mechanism that may be used to 
transfer files, or large volumes of data between two endpoints.

Implement resumable file downloads.

## Implementation

### Client

Requests a large file from the server using an asynchronous credit-based flow 
control to throttle the volume of data being sent to the client.

### Server

Receives requests from client to which it sends a large (fake) file. The client
must instruct the server that it has capacity before the server will send anything. 
This happens asynchronously.

### States

![client and server state diagrams](../images/PipeliningStateDiagrams.png)

### Formal Grammar

The following ABNF defines the protocol:

```abnf

```

## Security

All messages are sent between nodes in plain text.

