package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type Client struct {
	conn *nats.Conn
}

func Connect(uri string) *Client {
	conn, err := nats.Connect(uri)
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}
	fmt.Println("Connected to NATS at:", conn.ConnectedUrl())
	client := Client{conn: conn}

	return &client
}

func (nc *Client) Write(msg []byte, topic string) error {
	err := nc.conn.Publish(topic, msg)
	if err != nil {
		fmt.Printf("Error publishing message to NATS: %v/", err)
	}
	return err
}
