package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type Subscriber struct {
	name         string
	subscription *nats.Subscription
	err          error
}

type Client struct {
	conn        *nats.Conn
	Subscribers []*Subscriber
	ChVisit     chan *nats.Msg
	ChActivity  chan *nats.Msg
}

func Connect(uri string, errorChan chan error) *Client {
	conn, err := nats.Connect(uri)
	if err != nil {
		log.Fatal("Error establishing connection to NATS:", err)
	}

	fmt.Println("Connected to NATS at:", conn.ConnectedUrl())
	client := Client{
		conn:       conn,
		ChVisit:    make(chan *nats.Msg, 100),
		ChActivity: make(chan *nats.Msg, 100),
	}

	client.conn.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		fmt.Printf("reader error handler %s\n", err.Error())
		errorChan <- err
	})

	client.conn.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		fmt.Printf("reader disconnect error handler: %v", err.Error())
		errorChan <- err
	})

	client.conn.SetClosedHandler(func(_ *nats.Conn) {
		fmt.Printf("reader close handler")
	})

	return &client
}

func (nc *Client) Subscribe(topic string) error {
	var subscription *nats.Subscription
	var err error

	switch topic {
	case "visit":
		subscription, err = nc.conn.ChanSubscribe(topic, nc.ChVisit)
	case "activity":
		subscription, err = nc.conn.ChanSubscribe(topic, nc.ChActivity)
	default:
		log.Fatal("Error subscribing to NATS topic: wrong topic")
	}
	if err != nil {
		log.Fatal("Error subscribing to NATS topic:", err)
	}

	fmt.Printf("Subscribed to: %s\n", topic)

	nc.Subscribers = append(nc.Subscribers, &Subscriber{
		name:         topic,
		subscription: subscription,
	})

	return nil
}
