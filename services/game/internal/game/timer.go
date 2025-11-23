package game

import (
	"time"
)

// Timer represents a game timer
type Timer struct {
	C      <-chan time.Time
	timer  *time.Timer
	active bool
}

// NewTimer creates a new Timer
func NewTimer() *Timer {
	return &Timer{
		active: false,
	}
}

// Start starts the timer with the given duration
func (t *Timer) Start(duration time.Duration) {
	t.Stop() // Stop any existing timer

	t.timer = time.NewTimer(duration)
	t.C = t.timer.C
	t.active = true
}

// Stop stops the timer
func (t *Timer) Stop() {
	if t.timer != nil && t.active {
		t.timer.Stop()
		t.active = false
	}
}

// IsActive returns whether the timer is currently active
func (t *Timer) IsActive() bool {
	return t.active
}

