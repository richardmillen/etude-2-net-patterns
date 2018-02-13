// this app makes calls to several services via a 'service-broker' instance.
// first it calls the 'rand-service' a user defined number of times. it then
// calls the 'sort-service' to sort the random numbers immediately, or calls
// the 'offline-service' which calls the 'sort-service' on the clients behalf.
// if the latter option is chosen then the client periodically checks the
// 'offline-service' for status updates, then eventually gets the result.

package main

func main() {

}
