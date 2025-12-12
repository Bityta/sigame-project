package button

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewButtonPress(t *testing.T) {
	bp := NewButtonPress()

	if bp == nil {
		t.Fatal("NewButtonPress() returned nil")
	}
	if bp.closed {
		t.Error("NewButtonPress() closed = true, want false")
	}
	if len(bp.entries) != 0 {
		t.Errorf("NewButtonPress() entries length = %d, want 0", len(bp.entries))
	}
	if len(bp.pressedUsers) != 0 {
		t.Errorf("NewButtonPress() pressedUsers length = %d, want 0", len(bp.pressedUsers))
	}
}

func TestButtonPress_Reset(t *testing.T) {
	bp := NewButtonPress()
	userID := uuid.New()

	bp.Press(userID, "user1", 100*time.Millisecond)
	bp.Close()

	bp.Reset()

	if bp.closed {
		t.Error("Reset() closed = true, want false")
	}
	if len(bp.entries) != 0 {
		t.Errorf("Reset() entries length = %d, want 0", len(bp.entries))
	}
	if len(bp.pressedUsers) != 0 {
		t.Errorf("Reset() pressedUsers length = %d, want 0", len(bp.pressedUsers))
	}
	if bp.questionAt.IsZero() {
		t.Error("Reset() questionAt is zero")
	}
}

func TestButtonPress_Press(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	userID := uuid.New()
	username := "testuser"
	rtt := 50 * time.Millisecond

	success := bp.Press(userID, username, rtt)

	if !success {
		t.Error("Press() returned false, want true")
	}
	if !bp.HasPresses() {
		t.Error("Press() HasPresses() = false, want true")
	}
	if bp.GetPressCount() != 1 {
		t.Errorf("Press() GetPressCount() = %d, want 1", bp.GetPressCount())
	}
}

func TestButtonPress_Press_Duplicate(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	userID := uuid.New()
	rtt := 50 * time.Millisecond

	bp.Press(userID, "user1", rtt)
	success := bp.Press(userID, "user1", rtt)

	if success {
		t.Error("Press() duplicate returned true, want false")
	}
	if bp.GetPressCount() != 1 {
		t.Errorf("Press() duplicate GetPressCount() = %d, want 1", bp.GetPressCount())
	}
}

func TestButtonPress_Press_Closed(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()
	bp.Close()

	userID := uuid.New()
	success := bp.Press(userID, "user1", 50*time.Millisecond)

	if success {
		t.Error("Press() after Close() returned true, want false")
	}
	if bp.HasPresses() {
		t.Error("Press() after Close() HasPresses() = true, want false")
	}
}

func TestButtonPress_Close(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	bp.Close()

	if !bp.IsClosed() {
		t.Error("Close() IsClosed() = false, want true")
	}
}

func TestButtonPress_IsClosed(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.IsClosed() {
		t.Error("IsClosed() = true, want false")
	}

	bp.Close()

	if !bp.IsClosed() {
		t.Error("IsClosed() after Close() = false, want true")
	}
}

func TestButtonPress_HasPresses(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.HasPresses() {
		t.Error("HasPresses() empty = true, want false")
	}

	bp.Press(uuid.New(), "user1", 50*time.Millisecond)

	if !bp.HasPresses() {
		t.Error("HasPresses() with presses = false, want true")
	}
}

func TestButtonPress_GetPressCount(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.GetPressCount() != 0 {
		t.Errorf("GetPressCount() empty = %d, want 0", bp.GetPressCount())
	}

	bp.Press(uuid.New(), "user1", 50*time.Millisecond)
	bp.Press(uuid.New(), "user2", 50*time.Millisecond)

	if bp.GetPressCount() != 2 {
		t.Errorf("GetPressCount() = %d, want 2", bp.GetPressCount())
	}
}

func TestButtonPress_GetWinner(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.GetWinner() != nil {
		t.Error("GetWinner() empty = non-nil, want nil")
	}

	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	bp.Press(user1, "user1", 100*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	bp.Press(user2, "user2", 50*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	bp.Press(user3, "user3", 200*time.Millisecond)

	winner := bp.GetWinner()

	if winner == nil {
		t.Fatal("GetWinner() returned nil")
	}

	allPresses := bp.GetAllPresses()
	if len(allPresses) == 0 {
		t.Fatal("No presses recorded")
	}

	expectedWinner := allPresses[0]
	if winner.UserID != expectedWinner.UserID {
		t.Errorf("GetWinner() UserID = %v, want %v (user with earliest adjusted time)", winner.UserID, expectedWinner.UserID)
	}
	if !winner.AdjustedTime.Equal(expectedWinner.AdjustedTime) {
		t.Errorf("GetWinner() AdjustedTime = %v, want %v", winner.AdjustedTime, expectedWinner.AdjustedTime)
	}
}

func TestButtonPress_GetAllPresses(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.GetAllPresses() != nil {
		t.Error("GetAllPresses() empty = non-nil, want nil")
	}

	user1 := uuid.New()
	user2 := uuid.New()
	user3 := uuid.New()

	bp.Press(user1, "user1", 100*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	bp.Press(user2, "user2", 50*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	bp.Press(user3, "user3", 200*time.Millisecond)

	presses := bp.GetAllPresses()

	if len(presses) != 3 {
		t.Errorf("GetAllPresses() length = %d, want 3", len(presses))
	}

	for i := 1; i < len(presses); i++ {
		if presses[i].AdjustedTime.Before(presses[i-1].AdjustedTime) {
			t.Error("GetAllPresses() not sorted by adjusted time")
		}
	}
}

func TestButtonPress_GetReactionTime(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	entry := &ButtonPressEntry{
		AdjustedTime: time.Now(),
	}

	if bp.GetReactionTime(entry) != 0 {
		t.Error("GetReactionTime() with zero questionAt = non-zero, want 0")
	}

	bp.Reset()
	time.Sleep(100 * time.Millisecond)

	userID := uuid.New()
	bp.Press(userID, "user1", 50*time.Millisecond)

	presses := bp.GetAllPresses()
	if len(presses) == 0 {
		t.Fatal("No presses recorded")
	}

	reactionTime := bp.GetReactionTime(&presses[0])

	if reactionTime <= 0 {
		t.Errorf("GetReactionTime() = %d, want > 0", reactionTime)
	}
}

func TestButtonPress_GetReactionTime_NilEntry(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	if bp.GetReactionTime(nil) != 0 {
		t.Error("GetReactionTime() with nil entry = non-zero, want 0")
	}
}

func TestButtonPress_GetQuestionTime(t *testing.T) {
	bp := NewButtonPress()

	if !bp.GetQuestionTime().IsZero() {
		t.Error("GetQuestionTime() before Reset() = non-zero, want zero")
	}

	bp.Reset()

	if bp.GetQuestionTime().IsZero() {
		t.Error("GetQuestionTime() after Reset() = zero, want non-zero")
	}
}

func TestButtonPress_RTTCompensation(t *testing.T) {
	bp := NewButtonPress()
	bp.Reset()

	userID := uuid.New()
	rtt := 100 * time.Millisecond
	beforePress := time.Now()

	bp.Press(userID, "user1", rtt)

	presses := bp.GetAllPresses()
	if len(presses) == 0 {
		t.Fatal("No presses recorded")
	}

	entry := presses[0]
	afterPress := time.Now()

	if entry.AdjustedTime.Before(beforePress.Add(-rtt/RTTCompensationFactor - 10*time.Millisecond)) {
		t.Error("AdjustedTime is too early")
	}
	if entry.AdjustedTime.After(afterPress) {
		t.Error("AdjustedTime is after received time")
	}
	if entry.ReceivedAt.Before(beforePress) || entry.ReceivedAt.After(afterPress.Add(10*time.Millisecond)) {
		t.Error("ReceivedAt is not within expected range")
	}
}

