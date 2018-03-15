# Example 4: Timeouts

A client app sends a numeric value to a server, the server then sleeps for the
specified amount of time before replying. The client may timeout while waiting
for the server to respond. The client then performs some action based on the 
timeout *(the precise action is configurable)*.

## Goals

1. Explore timeouts using time-based event triggers.
1. Explore different strategies to timeouts.

## Implementation

### Overall Behaviour

Communication between the client and server follows a synchronous request-reply 
pattern, until a timeout occurs at which point the client may send zero or more
messages without receiving a response from the server.

A client is started with the following command line flags:

| Flag           | Type    | Description                                                              |
| :------------- | :------ | :----------------------------------------------------------------------- |
| Sleep          | Integer | *Sleep Duration* value of first request.                                 |
| Step           | Integer | Number of milliseconds to increment *Sleep* for each subsequent request. |
| Timeout        | Integer | Number of milliseconds to wait for server to reply.                      |
| Retry Interval | Integer | Number of milliseconds between each retry attempt.                       |
| Mode           | String  | Specifies how the client/server should behave.                           |

A client sends a maximum of 10 new requests (excluding retries etc) to the server
before terminating.

Each message sent from the client to the server contains the following parameters:

| Parameter      | Type    | Details                                                     |
| :------------- | :------ | :---------------------------------------------------------- |
| Id             | Integer | Starts at 1, incremented for each successful request/reply. |
| Retry Counter  | Integer | Set to 0 (zero) for new requests and incremented per retry. |
| Block          | Boolean | If set, the server connection will not receive client 
requests until the *Sleep Duration* has elapsed. |
| Sleep Duration | Integer | Number of milliseconds the server should sleep before
sending a reply. Ignored by server if *Retry Counter* &ne; 0 (zero). |

The client may be configured to start in various different *modes*, where each
mode may influence both client and server behaviour:

| Client Mode     | Client Behaviour                                                              | Server Behaviour |
| :-------------- | :---------------------------------------------------------------------------- | :--------------- |
| Abort           | On timeout; aborts request and exits.                                         | n/a |
| Retry           | On timeout; resends request with incremented *Retry Counter*.                 | 
Server responds to first retry, but not to original request. |
| Retry Blocked   | Same as *Retry*.                                                              | 
Server blocks, responding to request after timeout, retries ignored. |
| Refresh         | On timeout; resends request up to 3 times with *Retry Counter* set to zero 0. | 
Server replies to each request attempt. |
| Refresh Blocked | Same as *Refresh*.                                                            | 
Server blocks, responding to each request attempt. |
| Backoff         | Same as *Retry*, but *Retry Interval* doubled each time.                      | 
Server blocks, responding to request after timeout, retries ignored. |

### States

![client and server states](../images/Timeouts-StateDiagrams.png)

### Formal Grammar

The following ABNF grammar defines the protocol:

```abnf
...
```






