package timer

import (
	"sync"
	"time"
)

type Timer struct {
	C         chan time.Time
	timer     *time.Timer
	active    bool
	mu        sync.Mutex
	stopped   chan struct{}
	startedAt time.Time
	duration  time.Duration
}

func New() *Timer {
	return &Timer{
		C:       make(chan time.Time, ChannelBufferSize),
		active:  false,
		stopped: make(chan struct{}),
	}
}

func (t *Timer) Start(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stopInternal()

	t.timer = time.NewTimer(duration)
	t.active = true
	t.startedAt = time.Now()
	t.duration = duration

	go func() {
		select {
		case tick := <-t.timer.C:
			select {
			case t.C <- tick:
			default:
			}
		case <-t.stopped:
		}
	}()
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopInternal()
}

func (t *Timer) stopInternal() {
	if t.timer != nil && t.active {
		t.timer.Stop()
		t.active = false

		select {
		case t.stopped <- struct{}{}:
		default:
		}

		select {
		case <-t.C:
		default:
		}
	}
}

func (t *Timer) IsActive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}

func (t *Timer) Remaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.active {
		return InactiveRemaining
	}

	elapsed := time.Since(t.startedAt)
	remaining := t.duration - elapsed
	if remaining < 0 {
		return InactiveRemaining
	}
	return int(remaining.Seconds())
}
