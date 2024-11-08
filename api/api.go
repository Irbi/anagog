package main

import (
	"github.com/Irbi/anagog/api/internal/nats"
	"github.com/Irbi/anagog/api/internal/server"
	"log"
	"os"
)

func main() {
	natsUri := os.Getenv("NATS_URI")
	apiPort := os.Getenv("API_PORT")
	apiVersion := os.Getenv("API_VERSION")

	if apiPort == "" {
		apiPort = "8080"
	}
	if apiVersion == "" {
		apiVersion = "/v1/"
	}

	natsClient := nats.Connect(natsUri)

	messages := make(chan server.Msg, 100)

	go func() {
		for {
			msg := <-messages
			if err := natsClient.Write(msg.Data, msg.Name); err != nil {
				log.Fatal("Error publishing to NATS", err)
			}
		}
	}()

	server.Run(apiPort, apiVersion, messages)
}
