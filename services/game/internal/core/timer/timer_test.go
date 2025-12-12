package timer

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	timer := New()

	if timer == nil {
		t.Fatal("New returned nil")
	}

	if timer.C == nil {
		t.Error("Timer channel should not be nil")
	}

	if timer.active {
		t.Error("New timer should not be active")
	}

	if timer.stopped == nil {
		t.Error("Timer stopped channel should not be nil")
	}

	if timer.timer != nil {
		t.Error("New timer should not have internal timer set")
	}
}

func TestTimer_Start(t *testing.T) {
	timer := New()
	duration := 100 * time.Millisecond

	timer.Start(duration)

	if !timer.IsActive() {
		t.Error("Timer should be active after Start")
	}

	select {
	case <-timer.C:
		t.Log("Timer tick received")
	case <-time.After(150 * time.Millisecond):
		t.Error("Timer did not fire within expected time")
	}
}

func TestTimer_Start_StopsPrevious(t *testing.T) {
	timer := New()
	firstDuration := 200 * time.Millisecond
	secondDuration := 100 * time.Millisecond

	timer.Start(firstDuration)

	if !timer.IsActive() {
		t.Error("Timer should be active after first Start")
	}

	timer.Start(secondDuration)

	if !timer.IsActive() {
		t.Error("Timer should be active after second Start")
	}

	select {
	case <-timer.C:
		t.Log("Timer tick received")
	case <-time.After(150 * time.Millisecond):
		t.Error("Timer did not fire within expected time")
	}

	select {
	case <-timer.C:
		t.Error("Timer should not fire twice")
	case <-time.After(50 * time.Millisecond):
		t.Log("Timer correctly fired only once")
	}
}

func TestTimer_Stop(t *testing.T) {
	timer := New()
	duration := 200 * time.Millisecond

	timer.Start(duration)

	if !timer.IsActive() {
		t.Error("Timer should be active after Start")
	}

	timer.Stop()

	if timer.IsActive() {
		t.Error("Timer should not be active after Stop")
	}

	select {
	case <-timer.C:
		t.Error("Timer should not fire after Stop")
	case <-time.After(250 * time.Millisecond):
		t.Log("Timer correctly did not fire after Stop")
	}
}

func TestTimer_IsActive(t *testing.T) {
	timer := New()

	if timer.IsActive() {
		t.Error("New timer should not be active")
	}

	timer.Start(100 * time.Millisecond)

	if !timer.IsActive() {
		t.Error("Timer should be active after Start")
	}

	timer.Stop()

	if timer.IsActive() {
		t.Error("Timer should not be active after Stop")
	}
}

func TestTimer_Remaining(t *testing.T) {
	timer := New()

	remaining := timer.Remaining()
	if remaining != InactiveRemaining {
		t.Errorf("Expected remaining %d for inactive timer, got %d", InactiveRemaining, remaining)
	}

	duration := 2 * time.Second
	timer.Start(duration)

	remaining = timer.Remaining()
	if remaining <= 0 || remaining > 2 {
		t.Errorf("Expected remaining between 1 and 2 seconds, got %d", remaining)
	}

	time.Sleep(500 * time.Millisecond)

	remaining = timer.Remaining()
	if remaining <= 0 || remaining > 2 {
		t.Errorf("Expected remaining between 0 and 2 seconds after sleep, got %d", remaining)
	}

	timer.Stop()

	remaining = timer.Remaining()
	if remaining != InactiveRemaining {
		t.Errorf("Expected remaining %d after Stop, got %d", InactiveRemaining, remaining)
	}
}

func TestTimer_Channel(t *testing.T) {
	timer := New()
	duration := 50 * time.Millisecond

	timer.Start(duration)

	select {
	case tick := <-timer.C:
		if tick.IsZero() {
			t.Error("Timer tick should not be zero")
		}
		t.Log("Successfully received timer tick")
	case <-time.After(100 * time.Millisecond):
		t.Error("Timer did not send tick to channel")
	}

	timer.Start(50 * time.Millisecond)

	select {
	case tick := <-timer.C:
		if tick.IsZero() {
			t.Error("Timer tick should not be zero")
		}
		t.Log("Successfully received second timer tick")
	case <-time.After(100 * time.Millisecond):
		t.Error("Timer did not send second tick to channel")
	}
}

func TestTimer_Remaining_Negative(t *testing.T) {
	timer := New()
	duration := 10 * time.Millisecond

	timer.Start(duration)
	time.Sleep(50 * time.Millisecond)

	remaining := timer.Remaining()
	if remaining != InactiveRemaining {
		t.Errorf("Expected remaining %d after timer expired, got %d", InactiveRemaining, remaining)
	}
}

func TestTimer_MultipleStarts(t *testing.T) {
	timer := New()

	for i := 0; i < 5; i++ {
		timer.Start(50 * time.Millisecond)
		if !timer.IsActive() {
			t.Errorf("Timer should be active after Start iteration %d", i)
		}
		time.Sleep(10 * time.Millisecond)
	}

	timer.Stop()

	if timer.IsActive() {
		t.Error("Timer should not be active after Stop")
	}
}

