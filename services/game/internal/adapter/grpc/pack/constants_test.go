package pack

import (
	"testing"
	"time"
)

func TestConstantsAreNotEmpty(t *testing.T) {
	if SchemeHTTP == "" {
		t.Error("SchemeHTTP is empty")
	}

	if PathPackContent == "" {
		t.Error("PathPackContent is empty")
	}

	if PathPack == "" {
		t.Error("PathPack is empty")
	}
}

func TestTimeoutConstantIsValid(t *testing.T) {
	if DefaultHTTPTimeout <= 0 {
		t.Error("DefaultHTTPTimeout should be positive")
	}

	if DefaultHTTPTimeout >= 1*time.Minute {
		t.Error("DefaultHTTPTimeout should be reasonable")
	}
}

func TestConstantsHaveExpectedValues(t *testing.T) {
	if SchemeHTTP != "http" {
		t.Errorf("SchemeHTTP = %s, want http", SchemeHTTP)
	}

	if DefaultHTTPTimeout != 10*time.Second {
		t.Errorf("DefaultHTTPTimeout = %v, want 10s", DefaultHTTPTimeout)
	}
}

