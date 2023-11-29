# spacecraft

## cmd/server

This web server does a couple of things:

1. loads spacecraft from a JSON file while booting up
1. launches a web server exposing a single endpoint reachable via the address `/spacecraft` with a GET request

### The Getspacecraft handler

The request will accept two parameters provided via the query string:

- the `pageSize` option specifying how many spacecraft the requested page should have
- the `pageNumber` option specifying which page number you care about

Both of the above-mentioned parameters could be omitted. If that's the case, two default values will be used.

The endpoint will reply with the following HTTP responses:

1. `200` if everything goes right
1. `400` in case of a malformed incoming HTTP request
1. `500` for any other unexpected errors

It also sleeps for two seconds to emulate a heavy workload on it.

### How to run it

Navigate to the folder `cmd/server` and issue the command `go run .`. It terminal should wait for any incoming HTTP requests.

### How to test it

Issue the following cURL command: `curl "http://localhost:8080/spacecraft?pageSize=100&pageNumber=0"`

Please note that you need to wrap the address within double quotes. Otherwise, the OS interprets the `&` character as the terminator of the command and you'll get unexpected behavior.

To run the test suite issue the command: `go test -v -cover ./...`
Please be sure to be located within the `cmd/server` folder.

## ## cmd/elasticsearch (TBD)
