package main

import (
	"github.com/Irbi/anagog/worker/internal/collector"
	"github.com/Irbi/anagog/worker/internal/nats"
	"log"
	"os"
	"time"
)

func main() {
	natsUri := os.Getenv("NATS_URI")
	if natsUri == "" {
		natsUri = "nats://0.0.0.0:4222"
	}

	storage := os.Getenv("FILESTORAGE_PATH")
	if storage == "" {
		log.Fatal("FILESTORAGE_PATH environment variable not set")
	}
	errorChan := make(chan error, 1)
	natsClient := nats.Connect(natsUri, errorChan)

	if err := natsClient.Subscribe("visit"); err != nil {
		log.Fatal(err)
	}
	if err := natsClient.Subscribe("activity"); err != nil {
		log.Fatal(err)
	}

	c := collector.NewCollector(storage, 5*time.Second)
	c.AppendInputChannel("visit", natsClient.ChVisit)
	c.AppendInputChannel("activity", natsClient.ChActivity)

	c.Run()
}
