package timeglob

import (
	"time"
)

type Ticker struct {
	// Send time values to this on each tick.
	C <-chan time.Time

	// What do we need to be able to do this?
	tg       *TimeGlob
	sendTick chan<- time.Time
	timer    *time.Timer
}

func (tg *TimeGlob) Ticker() *Ticker {
	sendTick := make(chan time.Time, 1)

	result := &Ticker{
		sendTick,
		tg,
		sendTick,
		nil,
	}

	// Create the timer, now that there is a ticker to call tick() on.
	now := time.Now().In(tg.location)
	next := tg.Next(now)
	if next != UNKNOWN {
		result.timer = time.AfterFunc(next.Sub(now), result.tick)
	}
	return result
}

func (t *Ticker) tick() {
	now := time.Now().In(t.tg.location)

	// Never block. The channel already has a buffered value.
	select {
	case t.sendTick <- now:
	default:
	}

	next := t.tg.Next(now)
	if next != UNKNOWN {
		t.timer.Reset(next.Sub(now))
	}
}

func (t *Ticker) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
}
