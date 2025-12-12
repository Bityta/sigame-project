package button

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

const RTTCompensationFactor = 2

type PressEntry struct {
	UserID       uuid.UUID
	Username     string
	ReceivedAt   time.Time
	AdjustedTime time.Time
	RTT          time.Duration
}

type Press struct {
	entries      []PressEntry
	pressedUsers map[uuid.UUID]bool
	questionAt   time.Time
	closed       bool
	mu           sync.Mutex
}

func New() *Press {
	return &Press{
		entries:      make([]PressEntry, 0),
		pressedUsers: make(map[uuid.UUID]bool),
		closed:       false,
	}
}

func (b *Press) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = make([]PressEntry, 0)
	b.pressedUsers = make(map[uuid.UUID]bool)
	b.questionAt = time.Now()
	b.closed = false
}

func (b *Press) Press(userID uuid.UUID, username string, rtt time.Duration) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return false
	}

	if b.pressedUsers[userID] {
		return false
	}

	now := time.Now()
	oneWayDelay := rtt / RTTCompensationFactor
	adjustedTime := now.Add(-oneWayDelay)

	entry := PressEntry{
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

func (b *Press) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
}

func (b *Press) IsClosed() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.closed
}

func (b *Press) HasPresses() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries) > 0
}

func (b *Press) GetPressCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.entries)
}

func (b *Press) GetWinner() *PressEntry {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) == 0 {
		return nil
	}

	winner := &b.entries[0]
	for i := 1; i < len(b.entries); i++ {
		if b.entries[i].AdjustedTime.Before(winner.AdjustedTime) {
			winner = &b.entries[i]
		}
	}

	return winner
}

func (b *Press) GetAllPresses() []PressEntry {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.entries) == 0 {
		return nil
	}

	result := make([]PressEntry, len(b.entries))
	copy(result, b.entries)

	sort.Slice(result, func(i, j int) bool {
		return result[i].AdjustedTime.Before(result[j].AdjustedTime)
	})

	return result
}

func (b *Press) GetReactionTime(entry *PressEntry) int64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.questionAt.IsZero() || entry == nil {
		return 0
	}

	return entry.AdjustedTime.Sub(b.questionAt).Milliseconds()
}

func (b *Press) GetQuestionTime() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.questionAt
}
