FROM golang:1.23.2

COPY . /go

RUN go get github.com/nats-io/nats.go
RUN go build worker.go

EXPOSE 8082
ENTRYPOINT ["/go/worker"]