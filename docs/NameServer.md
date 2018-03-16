# Example 5: Name Lookup Server

A service registers itself with a central name lookup server which service consumers contact prior to sending requests. 
The central name lookup server returns the service address to calling clients.

Supports arbitrary registration groups *(environments)*, where for example, one instance of an *echo* service might register itself on the *development* environment and another on *test* and so forth.









