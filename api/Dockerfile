FROM golang:1.23.2

COPY . /go

RUN go get github.com/nats-io/nats.go
RUN go build api.go

EXPOSE 8080
ENTRYPOINT ["/go/api"]