# Example 10: Clone Server

The clone server is a simple echo server, but with some important differences.

The server app may be started as a **Helper**, where the address of a peer server
instance is specified via the commandline flag `--tbc`.

The server given to a **Helper** is referred to as the **Primary Lookup** server
for the **Helper**. The **Helper** sends a message to the **Primary Lookup**
telling it that it's available to help and then periodically requests an up-to-date
**Active List** of peers which includes other **Helper**s known to the **Primary Lookup**
server.

The client app receives the latest **Active List** immediately after connecting 
to a server. It then requests updates periodically.

The example is designed to support basic redundancy, where

1. The most active **Helper** will be promoted to **Primary Lookup** server if the **Primary Lookup** server goes offline.
1. The client app sends requests to the current active server.......

The list of *active* servers is returned directly, or indirectly by a **Primary Lookup**
is sorted by server *liveness* (how active the servers are) and hence resembles a 
priority queue.

## Modeling / Notes

### Start Server

+ have peer address?
    - send `HERETOHELP` (inc. *this* server address)
    - receive `THX` (inc. **Active List**)
    - run "Refresh Lookup List" *(run concurrently)*
+ listen *(loop)*
    - receive `GETSERVERS`
        - from server in **Helper List**?
            - update *liveness*
        - send **Active List** (*this* + **Lookup List** + **Helper List**)
    - receive `HERETOHELP`
        - from server in **Lookup List**?
            - remove from list
        - add server to **Helper List**
        - send **Active List** (*this* + **Lookup List** + **Helper List**)
    - receive `ECHO`
        - send `ECHO`

### Refresh Lookup List

+ wait (based on *liveness* of **Active Lookup** server; exponential back-off)
+ start refresh
    - send `GETSERVERS` (to **Active Lookup** server)
    - receive **Lookup List**
    - go to 'wait'
    - **ERROR**:
        - update server *liveness*
        - go to 'start refresh'







