// this service registers itself with a running 'service-broker' then receives
// requests that describe actions to be performed later (offline) on behalf of
// clients. it returns a unique 'job' identifier which is then used by the client
// to get status updates / results.
//
// this approach can be useful where the requested action can take a long time,
// or involves a service that isn't always available. it could also be useful for
// scheduling tasks.

package main

func main() {

}
