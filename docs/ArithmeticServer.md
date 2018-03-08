# Example 4: Arithmetic Server

This example follows the common client-server architecture.

The arithmetic client sends a series of basic arithmetical expressions to a remote
server which returns the result.

Each arithmetical expression is sent piece by piece i.e.

1. operand
1. operator
1. operand

the server then returns the result.

this implementation only supports very simple arithmetic operations i.e. n+n, n/n etc.
where 'n' is a 32-bit float.

n.b. type safety is moot (if not misplaced) at the point where data is passed into the
Server because validation (and serialisation) would be performed within the Service by
the current State, or more accurately by the fsm.Input on the associated event (fsm.Event).
so code such as the following which copies values to the Service using the relevant
input types resembles what would happen within a simple call to Service.Send()/.Execute()
(or whatever the API ends up looking like):

```go
msgs.Num.Copy(a, svc)
msgs.Op.Copy(op, svc)
msgs.Num.Copy(b, svc)
```








