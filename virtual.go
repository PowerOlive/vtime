package vtime

import (
	"sync"
	"time"
)

type virtualClock struct {
	now     time.Time
	tickers map[int]*virtualTicker
	mx      sync.RWMutex
}

type virtualTicker struct {
	id   int
	cl   *virtualClock
	d    time.Duration
	last time.Time
	u    chan time.Time
	c    chan time.Time
}

func (cl *virtualClock) Advance(t time.Time) {
	cl.mx.Lock()
	if t.After(cl.now) {
		cl.now = t
		for _, tk := range cl.tickers {
			tk.advance(t)
		}
	}
	cl.mx.Unlock()
}

func (cl *virtualClock) Now() time.Time {
	cl.mx.RLock()
	now := cl.now
	cl.mx.RUnlock()
	return now
}

func (cl *virtualClock) NewTicker(d time.Duration) Ticker {
	cl.mx.Lock()
	tk := &virtualTicker{
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

func (tk *virtualTicker) C() <-chan time.Time {
	return tk.c
}

func (tk *virtualTicker) advance(now time.Time) {
	tk.u <- now
}

func (tk *virtualTicker) run() {
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

func (tk *virtualTicker) Stop() {
	tk.cl.mx.Lock()
	delete(tk.cl.tickers, tk.id)
	close(tk.u)
	tk.cl.mx.Unlock()
}
