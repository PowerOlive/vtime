package vtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTicker(t *testing.T) {
	now := time.Date(2015, 5, 6, 7, 8, 9, 10, time.UTC)
	cl := NewVirtualClock(now)
	tk := cl.NewTicker(5 * time.Millisecond)

	readTick := func() time.Time {
		select {
		case t := <-tk.C():
			return t
		default:
			return time.Time{}
		}
	}

	for i := 0; i <= 10; i++ {
		if i < 5 {
			assert.Equal(t, time.Time{}, readTick())
		}
		cl.Advance(now.Add(time.Duration(i) * time.Millisecond))
	}
	assert.Equal(t, now.Add(5*time.Millisecond), readTick())
}
