# Example 2: Hello World Server

> TODO: use ABNF to describe messages

| Client        | Server  |
| :------------ | :------ |
| `HELLO`       | `WORLD` |
| `HI`          | `ERROR` |
| anything else | ignore  |

## Client

The hello world client allows the user to enter string messages to be sent to a 
hello word server. Only "hello" is a valid request message. "hi" will result in
an error response from the server. Any other request will be ignored by the server.
Entering "quit" will shut down the client service and terminate the application.

Note that the client Service must be started somehow in order to spawn background
goroutines etc. this example does so explicitly with a call to Service.Start().




