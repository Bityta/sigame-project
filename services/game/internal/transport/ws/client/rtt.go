package client

import (
	"sync"
	"time"

	"sigame/game/internal/infrastructure/logger"
)

const MaxRTTSamples = 10

type RTTTracker struct {
	samples     []time.Duration
	avgRTT      time.Duration
	lastPingAt  time.Time
	mu          sync.RWMutex
}

func newRTTTracker() *RTTTracker {
	return &RTTTracker{
		samples: make([]time.Duration, 0, MaxRTTSamples),
	}
}

func (r *RTTTracker) UpdateRTT(rtt time.Duration, userID interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.samples = append(r.samples, rtt)

	if len(r.samples) > MaxRTTSamples {
		r.samples = r.samples[1:]
	}

	var total time.Duration
	for _, sample := range r.samples {
		total += sample
	}
	r.avgRTT = total / time.Duration(len(r.samples))

	logger.Debugf(nil, "[RTT] User %v: new sample=%v, avg=%v (samples=%d)",
		userID, rtt, r.avgRTT, len(r.samples))
}

func (r *RTTTracker) GetRTT() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.avgRTT
}

func (r *RTTTracker) SetLastPingSentAt(t time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastPingAt = t
}

func (r *RTTTracker) GetLastPingSentAt() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.lastPingAt
}

