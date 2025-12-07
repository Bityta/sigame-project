package game

import (
	"sync"
	"time"
)

// Timer represents a game timer with a stable channel for select
type Timer struct {
	C       chan time.Time // Public channel that never changes
	timer   *time.Timer
	active  bool
	mu      sync.Mutex
	stopped chan struct{}
}

// NewTimer creates a new Timer with a stable channel
func NewTimer() *Timer {
	return &Timer{
		C:       make(chan time.Time, 1), // Buffered to prevent blocking
		active:  false,
		stopped: make(chan struct{}),
	}
}

// Start starts the timer with the given duration
func (t *Timer) Start(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Stop any existing timer
	t.stopInternal()

	// Create new timer
	t.timer = time.NewTimer(duration)
	t.active = true

	// Forward timer events to our stable channel
	go func() {
		select {
		case tick := <-t.timer.C:
			// Forward to our channel (non-blocking)
			select {
			case t.C <- tick:
			default:
				// Channel full, drop the tick
			}
		case <-t.stopped:
			// Timer was stopped
		}
	}()
}

// Stop stops the timer
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopInternal()
}

// stopInternal stops the timer (must be called with lock held)
func (t *Timer) stopInternal() {
	if t.timer != nil && t.active {
		t.timer.Stop()
		t.active = false
		
		// Signal the forwarding goroutine to stop
		select {
		case t.stopped <- struct{}{}:
		default:
		}
		
		// Drain the public channel
		select {
		case <-t.C:
		default:
		}
	}
}

// IsActive returns whether the timer is currently active
func (t *Timer) IsActive() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}
