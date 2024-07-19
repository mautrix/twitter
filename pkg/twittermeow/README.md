# x-go
A Go library for interacting with X's API

### Test steps

1. Create cookies.txt in this directory, grab your cookie string from x.com and paste it there
2. Run `go test client_test.go -v`

### Testing functionality

```go
	_, _, err = cli.LoadMessagesPage()
	if err != nil {
		log.Fatal(err)
	}
```
The `LoadMessagesPage` method makes a request to `https://x.com/messages` then makes 2 calls:
```go
	data, err := c.GetAccountSettings(...)
	initialInboxState, err := c.GetInitialInboxState(...)
```
it sets up the current "page" session for the client, fetches the current authenticated user info as well as the initial inbox state (the very starting inbox information you see when u load `/messages`) then returns the parsed data.

To easily test with the available functions I have made, lets say you wanna test uploading an image and sending it to the top conversation in your inbox you could simply do something like:
```go
	initialInboxData, _, err := cli.LoadMessagesPage()
	if err != nil {
		log.Fatal(err)
	}
    	uploadAndSendImageTest(initialInboxData)
```
Or feel free to try it out yourself! All the methods are available on the client instance.
