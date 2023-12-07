# spacecraft

## cmd/server

This web server does a couple of things:

1. loads spacecraft from a JSON file
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

It also sleeps for one second to emulate a heavy workload on it.

### How to run it

Navigate to the folder `cmd/server` and issue the command `go run .`. It terminal should wait for any incoming HTTP requests.

### How to test it

Issue the following cURL command: `curl "http://localhost:8080/spacecraft?pageSize=100&pageNumber=0"`

Please note that you need to wrap the address within double quotes. Otherwise, the OS interprets the `&` character as the terminator of the command and you'll get unexpected behavior.

To run the test suite issue the command: `go test -v -cover`
Please be sure to be located within the `internal/handlers` folder.

#### Integration tests

Within the `test` folder at the root level, you can also find integration tests. To run them you need to issue (located within the `test` folder):

`go test -v`

To accomplish these tests a couple of things have been used:

1. the `testcontainers` package to build a temporary image starting from a `Dockerfile`
1. the `testify` library to ease the test experience. Specifically, we used the `assert`, `require`, and `suite` packages

### How to deploy it

To deploy the web server, there is a Dockerfile, located at the root directory of the project, that allows you to build a Docker image. To build it, be sure to be at the same level as the Dockerfile, then run:

`docker build -t <name of the Docker image> .`  

A tiny image should appear between yours. To check it run:

`docker image list`

Finally, to run a container based on this image, run what follows:

`docker run -p 7000:8080 spacecraft-web`

It maps our host `7000` port to the `8080` of the Docker container making it reachable from outside. To manually test the web server from the container, issue:

`curl "http://localhost:7000/spacecraft?pageSize=100&pageNumber=0"`

> Please note that to build this image, it has been used the multi-stage build approach. The final image consists of only a few megabytes (should be less than 10!). We also avoided including the Golang build tools by starting the final image from the `scratch` base image which is empty.

## ## cmd/elasticsearch (TBD)
