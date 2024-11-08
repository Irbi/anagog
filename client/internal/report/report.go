package report

import (
	"fmt"
	"github.com/Irbi/anagog/client/internal/schema"
	tools "github.com/Irbi/anagog/client/tools/geo"
	"math/rand"
	"sync"
	"time"
)

type Report struct {
	uids      int
	startTime time.Time
	endTime   time.Time
	ch        []schema.SourceChannels
}

func NewReport(uids int, d time.Time, names []string) Report {

	r := Report{
		uids:      uids,
		startTime: d,
		endTime:   d.Add(24 * time.Hour),
	}

	for _, n := range names {
		ch := make(chan schema.SourceMsg, 12*uids)
		sc := schema.SourceChannels{
			Name: n,
			Ch:   ch,
		}
		r.ch = append(r.ch, sc)
	}

	return r
}

func (r Report) GenerateDay() []schema.SourceChannels {
	wg := &sync.WaitGroup{}
	for _, c := range r.ch {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(c.Ch)
			r.generateQueue(c)
		}()
	}
	wg.Wait()

	fmt.Printf("Day %s ready\n", r.startTime.Format(time.DateOnly))

	return r.ch
}

func (v Report) generateQueue(ch schema.SourceChannels) {
	wg := &sync.WaitGroup{}
	for i := 0; i < v.uids; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.generateClient(i, ch)
		}()
	}
	wg.Wait()
}

func (v Report) generateClient(uid int, ch schema.SourceChannels) {
	d := v.startTime
	for d.Before(v.endTime) {
		switch ch.Name {
		case "visit":
			r := v.generateVisit(uid, d)
			ch.Ch <- schema.SourceMsg{Data: r}
			break
		case "activity":
			r := v.generateActivity(uid, d)
			ch.Ch <- schema.SourceMsg{Data: r}
			break
		default:
			fmt.Printf("Report type %s not supported\n", ch.Name)
		}

		d = d.Add(time.Hour * 2)

	}
}

func (v Report) generateVisit(uid int, timestamp time.Time) schema.VisitReport {
	return schema.VisitReport{
		DataVer:       1,
		UserId:        uid,
		EnterTime:     timestamp.Format(time.DateTime),
		ExitTime:      timestamp.Add(time.Minute * 40).Format(time.DateTime),
		AlgorithmType: rand.Intn(7 - 0),
		PoiId:         rand.Int63(),
		Latitude:      tools.RandGeo(-90, 90),
		Longitude:     tools.RandGeo(-180, 180),
	}
}

func (v Report) generateActivity(uid int, timestamp time.Time) schema.ActivityReport {
	return schema.ActivityReport{
		DataVer:        1,
		UserId:         uid,
		StartTime:      timestamp.Add(time.Hour * 1).Format(time.DateTime),
		EndTime:        timestamp.Add(time.Hour * 1).Add(time.Minute * 30).Format(time.DateTime),
		ActivityType:   rand.Intn(7 - 0),
		StartLatitude:  tools.RandGeo(-90, 90),
		StartLongitude: tools.RandGeo(-180, 180),
		EndLatitude:    tools.RandGeo(-90, 90),
		EndLongitude:   tools.RandGeo(-180, 180),
	}
}
