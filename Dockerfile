FROM golang:1.21 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/server/main.go ./main.go
COPY spacecraft.json .

RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver .

FROM bitnami/minideb:stretch

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder ./app/spacecraft.json .
COPY --from=builder /webserver /

CMD [ "/webserver" ]