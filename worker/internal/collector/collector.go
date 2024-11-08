package collector

import (
	"fmt"
	"github.com/Irbi/anagog/worker/tools/fwriter"
	"github.com/nats-io/nats.go"
	"time"
)

const DefaultLength = 10
const MaxIdle = time.Second * 10

type Collector struct {
	chVis       chan *nats.Msg
	chAct       chan *nats.Msg
	sVis        []string
	tickerVis   *time.Ticker
	sAct        []string
	tickerAct   *time.Ticker
	storagePath string
	Interval    time.Duration
	period      time.Duration
}

func NewCollector(storage string, interval time.Duration) *Collector {
	return &Collector{
		storagePath: storage,
		sVis:        make([]string, 0, DefaultLength),
		tickerVis:   time.NewTicker(MaxIdle),
		sAct:        make([]string, 0, DefaultLength),
		tickerAct:   time.NewTicker(MaxIdle),
		Interval:    interval,
		period:      interval,
	}
}

func (c *Collector) AppendInputChannel(name string, ch chan *nats.Msg) {
	switch name {
	case "visit":
		c.chVis = ch
	case "activity":
		c.chAct = ch
	}
}

func (c *Collector) Run() {
	for {
		select {
		case <-c.tickerVis.C:
			c.collectBatch("visit", c.sVis)
			c.sVis = make([]string, 0)
			break
		case <-c.tickerAct.C:
			c.collectBatch("activity", c.sAct)
			c.sAct = make([]string, 0)
			break

		case i := <-c.chVis:
			c.tickerVis.Reset(MaxIdle)
			if len(c.sVis) == DefaultLength {
				c.collectBatch(i.Subject, c.sVis)
				c.sVis = make([]string, 0)
			} else {
				c.sVis = append(c.sVis, string(i.Data))
			}
			break

		case i := <-c.chAct:
			c.tickerAct.Reset(MaxIdle)
			if len(c.sAct) == DefaultLength {
				c.collectBatch(i.Subject, c.sAct)
				c.sAct = make([]string, 0)
			} else {
				c.sAct = append(c.sAct, string(i.Data))
			}
			break
		}
	}
}

func (c *Collector) collectBatch(name string, data []string) {
	if len(data) == 0 {
		return
	}
	fName := fmt.Sprintf("%s/%s-%v.json", c.storagePath, name, time.Now().UnixMilli())
	fp, err := fwriter.CreateFile(fName)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", name, err)
	}
	for _, msg := range data {

		fwriter.AppendLines(fp, msg)
	}
}
