package vtime

import (
	"sync"
	"time"
)

type Clock struct {
	now     time.Time
	tickers map[int]*ticker
	mx      sync.RWMutex
}

func NewClock(now time.Time) *Clock {
	return &Clock{
		now:     now,
		tickers: make(map[int]*ticker, 0),
	}
}

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type ticker struct {
	id   int
	cl   *Clock
	d    time.Duration
	last time.Time
	u    chan time.Time
	c    chan time.Time
}

func (cl *Clock) Advance(t time.Time) {
	cl.mx.Lock()
	if t.After(cl.now) {
		cl.now = t
		for _, tk := range cl.tickers {
			tk.advance(t)
		}
	}
	cl.mx.Unlock()
}

func (cl *Clock) Now() time.Time {
	cl.mx.RLock()
	now := cl.now
	cl.mx.RUnlock()
	return now
}

func (cl *Clock) NewTicker(d time.Duration) Ticker {
	cl.mx.Lock()
	tk := &ticker{
		id:   len(cl.tickers),
		d:    d,
		cl:   cl,
		last: cl.now,
		u:    make(chan time.Time),
		c:    make(chan time.Time, 1),
	}
	cl.tickers[tk.id] = tk
	cl.mx.Unlock()
	go tk.run()
	return tk
}

func (tk *ticker) C() <-chan time.Time {
	return tk.c
}

func (tk *ticker) advance(now time.Time) {
	tk.u <- now
}

func (tk *ticker) run() {
	for now := range tk.u {
		if tk.last.IsZero() {
			// initialize as soon as we get our first time
			tk.last = now
			continue
		}
		if now.Sub(tk.last) >= tk.d {
			// New tick
			select {
			case tk.c <- now:
				// submitted tick
				tk.last = now
			default:
				// no on listening, drop it on the floor
			}
		}
	}
}

func (tk *ticker) Stop() {
	tk.cl.mx.Lock()
	delete(tk.cl.tickers, tk.id)
	close(tk.u)
	tk.cl.mx.Unlock()
}
