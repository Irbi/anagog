package main

import (
	"github.com/Irbi/anagog/client/internal/aggregator"
	"github.com/Irbi/anagog/client/internal/publisher"
	"github.com/Irbi/anagog/client/internal/report"
	"os"
	"strconv"
	"time"
)

var names = []string{"visit", "activity"}

func main() {
	clients, err := strconv.Atoi(os.Getenv("CLIENTS"))
	if err != nil {
		clients = 1
	}

	startDay := time.Date(2018, 1, 1, 0, 0, 0, 100, time.UTC)
	endDay := time.Now()
	curDay := startDay

	for !curDay.Equal(endDay) {
		time.Sleep(5 * time.Second)
		rep := report.NewReport(clients, curDay, names)
		rchs := rep.GenerateDay()

		agg := aggregator.NewAggregator(rchs)
		achs := agg.Aggregate()

		pub := publisher.NewPublisher(os.Getenv("API_URL"), os.Getenv("API_VERSION"), achs)
		pub.Publish()

		curDay = curDay.AddDate(0, 0, 1)
	}
}
