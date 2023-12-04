FROM golang:1.21

RUN mkdir /app
ADD . /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/server/main.go ./main.go
COPY spacecraft.json .

RUN CGO_ENABLED=0 GOOS=linux go build -o /webserver .

# nit: use multistage builds for prod images:
# https://docs.bitnami.com/tutorials/optimize-docker-images-multistage-builds/

CMD [ "/webserver" ]