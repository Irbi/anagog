package publisher

import (
	"fmt"
	"github.com/Irbi/anagog/client/internal/schema"
	"net/http"
	"sync"
	"time"
)

type Publisher struct {
	client  *http.Client
	url     string
	version string
	chIn    []schema.AggrChannels
}

func NewPublisher(url string, version string, ch []schema.AggrChannels) *Publisher {
	p := &Publisher{
		client:  &http.Client{Timeout: 10 * time.Second},
		url:     url,
		version: version,
	}
	for _, n := range ch {
		p.chIn = append(p.chIn, n)
	}

	return p
}

func (p *Publisher) Publish() {
	wg := &sync.WaitGroup{}
	for _, c := range p.chIn {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(c.Ch)
			p.post(c)
		}()
	}
	wg.Wait()
}

func (p *Publisher) post(ch schema.AggrChannels) {
	msg := <-ch.Ch
	url := p.url + ch.Name + p.version
	req, err := http.NewRequest("POST", url, msg.Data)
	if err != nil {
		fmt.Printf("Impossible to build request: %s\n", err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	_, err = p.client.Do(req)
	if err != nil {
		fmt.Printf("Impossible to send request: %s\n", err)
	}
}
