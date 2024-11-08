package aggregator

import (
	"bytes"
	"github.com/Irbi/anagog/client/internal/schema"
	"github.com/Irbi/anagog/client/tools/archiver"
	"github.com/simonfrey/jsonl"
	"log"
	"sync"
)

type Aggregator struct {
	chOut   []schema.AggrChannels
	chIn    []schema.SourceChannels
	mCnt    int
	writers []Writer
}

type Writer struct {
	Name   string
	Writer jsonl.Writer
	Buff   *bytes.Buffer
	err    error
}

func NewAggregator(ch []schema.SourceChannels) *Aggregator {

	a := &Aggregator{
		mCnt: 0,
	}
	for _, n := range ch {
		w := a.createWriter(n.Name)
		a.writers = append(a.writers, w)

		ch := make(chan schema.AggrMsg, 1)
		sc := schema.AggrChannels{
			Name: n.Name,
			Ch:   ch,
		}
		a.chOut = append(a.chOut, sc)
		a.chIn = append(a.chIn, n)
	}

	return a
}

func (a *Aggregator) Aggregate() []schema.AggrChannels {
	wg := &sync.WaitGroup{}
	for _, c := range a.chIn {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.writeBatch(c)
		}()
	}
	wg.Wait()

	return a.chOut
}

func (a *Aggregator) writeBatch(ch schema.SourceChannels) {
	w, err := a.getWriter(ch.Name)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range ch.Ch {
		err = w.Writer.Write(msg.Data)
		if err != nil {
			log.Fatal(err)
		}
	}

	gz, err := tools.Zip(w.Buff.String())
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range a.chOut {
		if c.Name == ch.Name {
			c.Ch <- schema.AggrMsg{Data: gz}
			break
		}
	}
}

func (a *Aggregator) getWriter(name string) (*Writer, error) {
	var writer *Writer
	var err error

	for _, w := range a.writers {
		if w.Name == name {
			writer = &w
		}
	}
	if writer == nil {
		log.Fatal("Writer not found", name)
	}

	return writer, err
}

func (a *Aggregator) createWriter(name string) Writer {
	buff := bytes.Buffer{}
	return Writer{
		Name:   name,
		Writer: jsonl.NewWriter(&buff),
		Buff:   &buff,
	}
}
