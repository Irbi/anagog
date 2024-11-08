package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/simonfrey/jsonl"
	"io"
	"log"
	"net/http"
	"strings"
)

type Msg struct {
	Name string
	Data []byte
}

func Run(port string, version string, ch chan Msg) {
	http.HandleFunc("/api/visit"+version, func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, ch, "visit")
	})
	http.HandleFunc("/api/activity"+version, func(w http.ResponseWriter, r *http.Request) {
		handle(w, r, ch, "activity")
	})

	fmt.Println("Server listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("could not open HTTP server", err)
	}
}

func handle(w http.ResponseWriter, r *http.Request, ch chan Msg, topic string) {
	reader, err := gzip.NewReader(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer reader.Close()
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}

	var data string
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}
	jReader := jsonl.NewReader(strings.NewReader(data))
	err = jReader.ReadLines(func(data []byte) error {
		ch <- Msg{
			Name: topic,
			Data: data,
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
