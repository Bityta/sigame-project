package game

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ButtonPress handles atomic button press functionality
type ButtonPress struct {
	pressed    bool
	pressedBy  uuid.UUID
	pressedAt  time.Time
	questionAt time.Time
	mu         sync.Mutex
}

// NewButtonPress creates a new ButtonPress handler
func NewButtonPress() *ButtonPress {
	return &ButtonPress{
		pressed: false,
	}
}

// Reset resets the button press state
func (b *ButtonPress) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pressed = false
	b.pressedBy = uuid.Nil
	b.pressedAt = time.Time{}
	b.questionAt = time.Now()
}

// Press attempts to press the button
// Returns true if this player was first to press, false otherwise
func (b *ButtonPress) Press(userID uuid.UUID) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.pressed {
		return false
	}

	b.pressed = true
	b.pressedBy = userID
	b.pressedAt = time.Now()

	return true
}

// GetPressedBy returns the user ID who pressed the button first
func (b *ButtonPress) GetPressedBy() uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pressedBy
}

// GetLatency returns the latency in milliseconds from question shown to button press
func (b *ButtonPress) GetLatency() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.pressedAt.IsZero() || b.questionAt.IsZero() {
		return 0
	}

	return b.pressedAt.Sub(b.questionAt).Milliseconds()
}

// IsPressed returns whether the button has been pressed
func (b *ButtonPress) IsPressed() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pressed
}

