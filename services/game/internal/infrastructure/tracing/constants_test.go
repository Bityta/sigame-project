package tracing

import (
	"testing"
	"time"
)

func TestConstantsAreNotEmpty(t *testing.T) {
	if EnvOTLPEndpoint == "" {
		t.Error("EnvOTLPEndpoint is empty")
	}

	if DefaultOTLPEndpoint == "" {
		t.Error("DefaultOTLPEndpoint is empty")
	}

	if DefaultServiceVersion == "" {
		t.Error("DefaultServiceVersion is empty")
	}
}

func TestTimeoutConstantsAreValid(t *testing.T) {
	if ConnectionTimeout <= 0 {
		t.Error("ConnectionTimeout should be positive")
	}

	if ShutdownTimeout <= 0 {
		t.Error("ShutdownTimeout should be positive")
	}

	if ConnectionTimeout >= 1*time.Minute {
		t.Error("ConnectionTimeout should be reasonable")
	}

	if ShutdownTimeout >= 1*time.Minute {
		t.Error("ShutdownTimeout should be reasonable")
	}
}

func TestRetryConstantsAreValid(t *testing.T) {
	if MaxRetries <= 0 {
		t.Error("MaxRetries should be positive")
	}

	if RetryBackoffBase <= 0 {
		t.Error("RetryBackoffBase should be positive")
	}
}

