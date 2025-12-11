package game

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ButtonPressEntry represents a single button press with RTT compensation
type ButtonPressEntry struct {
	UserID       uuid.UUID
	Username     string
	ReceivedAt   time.Time     // When server received the press
	AdjustedTime time.Time     // ReceivedAt - RTT/2 (compensated time)
	RTT          time.Duration // Client's RTT at time of press
}

// ButtonPress handles button press collection with RTT compensation
type ButtonPress struct {
	entries      []ButtonPressEntry
	pressedUsers map[uuid.UUID]bool // Track who already pressed (prevent duplicates)
	questionAt   time.Time          // When the question was shown / button enabled
	closed       bool               // Whether collection window is closed
	mu           sync.Mutex
}

// NewButtonPress creates a new ButtonPress handler
func NewButtonPress() *ButtonPress {
	return &ButtonPress{
		entries:      make([]ButtonPressEntry, 0),
		pressedUsers: make(map[uuid.UUID]bool),
		closed:       false,
	}
}

// Reset resets the button press state for a new question
func (b *ButtonPress) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = make([]ButtonPressEntry, 0)
	b.pressedUsers = make(map[uuid.UUID]bool)
	b.questionAt = time.Now()
	b.closed = false
}

// Press records a button press with RTT compensation
// Returns true if this is a valid new press, false if duplicate or closed
func (b *ButtonPress) Press(userID uuid.UUID, username string, rtt time.Duration) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if collection is closed
	if b.closed {
		return false
	}

	// Check if user already pressed (prevent duplicates)
	if b.pressedUsers[userID] {
		return false
	}

	now := time.Now()

	// Calculate adjusted time: server receive time - one-way delay (RTT/2)
	// This compensates for network latency
	oneWayDelay := rtt / 2
	adjustedTime := now.Add(-oneWayDelay)

	entry := ButtonPressEntry{
		UserID:       userID,
		Username:     username,
		ReceivedAt:   now,
		AdjustedTime: adjustedTime,
		RTT:          rtt,
	}

	b.entries = append(b.entries, entry)
	b.pressedUsers[userID] = true

	return true
}

// Close closes the collection window - no more presses accepted
func (b *ButtonPress) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
}

// IsClosed returns whether collection is closed
func (b *ButtonPress) IsClosed() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closed
}

// HasPresses returns whether any button presses have been recorded
func (b *ButtonPress) HasPresses() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries) > 0
}

// GetPressCount returns the number of presses recorded
func (b *ButtonPress) GetPressCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

// GetWinner returns the winner (earliest adjusted time) or nil if no presses
func (b *ButtonPress) GetWinner() *ButtonPressEntry {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) == 0 {
		return nil
	}

	// Find entry with earliest adjusted time
	winner := &b.entries[0]
	for i := 1; i < len(b.entries); i++ {
		if b.entries[i].AdjustedTime.Before(winner.AdjustedTime) {
			winner = &b.entries[i]
		}
	}

	return winner
}

// GetAllPresses returns all presses sorted by adjusted time (earliest first)
func (b *ButtonPress) GetAllPresses() []ButtonPressEntry {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) == 0 {
		return nil
	}

	// Make a copy to avoid returning reference to internal slice
	result := make([]ButtonPressEntry, len(b.entries))
	copy(result, b.entries)

	// Sort by adjusted time (earliest first)
	sort.Slice(result, func(i, j int) bool {
		return result[i].AdjustedTime.Before(result[j].AdjustedTime)
	})

	return result
}

// GetReactionTime returns the reaction time in milliseconds for a press entry
// (time from question shown to adjusted press time)
func (b *ButtonPress) GetReactionTime(entry *ButtonPressEntry) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.questionAt.IsZero() || entry == nil {
		return 0
	}

	return entry.AdjustedTime.Sub(b.questionAt).Milliseconds()
}

// GetQuestionTime returns when the question/button was enabled
func (b *ButtonPress) GetQuestionTime() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.questionAt
}
