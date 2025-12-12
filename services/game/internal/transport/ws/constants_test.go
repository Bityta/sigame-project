package ws

import (
	"testing"
	"time"
)

func TestConstantsAreNotEmpty(t *testing.T) {
	if QueryParamUserID == "" {
		t.Error("QueryParamUserID is empty")
	}
}

func TestErrorConstantsAreNotEmpty(t *testing.T) {
	errors := []string{
		ErrorInvalidGameID,
		ErrorInvalidUserID,
		ErrorUserIDRequired,
		ErrorGameNotFound,
		ErrorInvalidMessageFormat,
		ErrorGameManagerNotFound,
		ErrorMissingServerTime,
		ErrorInvalidServerTimeType,
		ErrorInvalidRTT,
	}

	for _, err := range errors {
		if err == "" {
			t.Errorf("error constant is empty: %s", err)
		}
	}
}

func TestTimeConstantsAreValid(t *testing.T) {
	if WriteWait <= 0 {
		t.Error("WriteWait should be positive")
	}

	if PongWait <= 0 {
		t.Error("PongWait should be positive")
	}

	if JSONPingPeriod <= 0 {
		t.Error("JSONPingPeriod should be positive")
	}

	if MaxRTTDuration <= 0 {
		t.Error("MaxRTTDuration should be positive")
	}

	if MaxRTTDuration > 10*time.Second {
		t.Error("MaxRTTDuration should be reasonable")
	}
}

func TestSizeConstantsAreValid(t *testing.T) {
	if MaxMessageSize <= 0 {
		t.Error("MaxMessageSize should be positive")
	}

	if MaxRTTSamples <= 0 {
		t.Error("MaxRTTSamples should be positive")
	}
}

