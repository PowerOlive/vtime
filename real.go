package vtime

import (
	"time"
)

type realClock struct {
}

func (c *realClock) Advance(t time.Time) {
	// do nothing
}

func (c *realClock) Now() time.Time {
	return time.Now()
}

func (c *realClock) NewTicker(d time.Duration) Ticker {
	return &realTicker{time.NewTicker(d)}
}

type realTicker struct {
	t *time.Ticker
}

func (t *realTicker) C() <-chan time.Time {
	return t.t.C
}

func (t *realTicker) Stop() {
	t.t.Stop()
}
