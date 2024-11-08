FROM golang:1.23.2

COPY . /go

RUN go mod download
RUN go build client.go

EXPOSE 8081
ENTRYPOINT ["/go/client"]