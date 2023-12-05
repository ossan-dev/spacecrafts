FROM golang:alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY cmd/server/main.go ./main.go
COPY cmd/server/store/spacecraft.json store/spacecraft.json

RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver .

FROM scratch

WORKDIR /app

COPY --from=builder /webserver /

CMD [ "/webserver" ]