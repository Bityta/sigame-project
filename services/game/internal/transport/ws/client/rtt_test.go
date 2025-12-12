package client

import (
	"testing"
	"time"
)

func TestRTTTracker_UpdateRTT(t *testing.T) {
	tracker := newRTTTracker()

	rtt1 := 50 * time.Millisecond
	tracker.UpdateRTT(rtt1, "user1")

	if tracker.GetRTT() == 0 {
		t.Error("GetRTT() should return non-zero after UpdateRTT")
	}

	rtt2 := 100 * time.Millisecond
	tracker.UpdateRTT(rtt2, "user1")

	avg := tracker.GetRTT()
	if avg <= 0 {
		t.Error("GetRTT() should return positive average")
	}
}

func TestRTTTracker_GetRTT_Initial(t *testing.T) {
	tracker := newRTTTracker()

	if tracker.GetRTT() != 0 {
		t.Error("GetRTT() should return 0 initially")
	}
}

func TestRTTTracker_SetGetLastPingSentAt(t *testing.T) {
	tracker := newRTTTracker()

	now := time.Now()
	tracker.SetLastPingSentAt(now)

	got := tracker.GetLastPingSentAt()
	if !got.Equal(now) {
		t.Errorf("GetLastPingSentAt() = %v, want %v", got, now)
	}
}

func TestRTTTracker_MaxSamples(t *testing.T) {
	tracker := newRTTTracker()

	for i := 0; i < MaxRTTSamples+5; i++ {
		tracker.UpdateRTT(time.Duration(i)*time.Millisecond, "user1")
	}

	rtt := tracker.GetRTT()
	if rtt <= 0 {
		t.Error("GetRTT() should return positive value after multiple updates")
	}
}

