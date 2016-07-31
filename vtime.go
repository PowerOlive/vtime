package vtime

import (
	"time"
)

var (
	// RealClock is a Clock that uses real time
	RealClock Clock = &realClock{}
)

type Clock interface {
	Advance(t time.Time)
	Now() time.Time
	NewTicker(d time.Duration) Ticker
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// NewVirtualClock constructs a clock that uses virtual time (time that only
// advances when Advance() is called).
func NewVirtualClock(now time.Time) Clock {
	return &virtualClock{
		now:     now,
		tickers: make(map[int]*virtualTicker, 0),
	}
}
