package redis

import "testing"

func TestParseUUIDs(t *testing.T) {
	tests := []struct {
		name    string
		strs    []string
		wantErr bool
	}{
		{
			name:    "valid UUIDs",
			strs:    []string{"550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"},
			wantErr: false,
		},
		{
			name:    "empty slice",
			strs:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			strs:    []string{"invalid-uuid"},
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid",
			strs:    []string{"550e8400-e29b-41d4-a716-446655440000", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseUUIDs(tt.strs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUUIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) != len(tt.strs) {
				t.Errorf("parseUUIDs() returned %d UUIDs, want %d", len(result), len(tt.strs))
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "string value",
			value:   "test",
			wantErr: false,
		},
		{
			name:    "int value",
			value:   42,
			wantErr: false,
		},
		{
			name:    "bool value true",
			value:   true,
			wantErr: false,
		},
		{
			name:    "bool value false",
			value:   false,
			wantErr: false,
		},
		{
			name:    "struct value",
			value:   struct{ Name string }{Name: "test"},
			wantErr: false,
		},
		{
			name:    "map value",
			value:   map[string]int{"key": 1},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToString(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == "" {
				t.Error("convertToString() returned empty string")
			}

			if tt.name == "int value" && result != "42" {
				t.Errorf("convertToString() = %v, want %v", result, "42")
			}
			if tt.name == "bool value true" && result != "true" {
				t.Errorf("convertToString() = %v, want %v", result, "true")
			}
			if tt.name == "string value" && result != "test" {
				t.Errorf("convertToString() = %v, want %v", result, "test")
			}
		})
	}
}

func TestConvertToString_InvalidJSON(t *testing.T) {
	invalidValue := make(chan int)
	_, err := convertToString(invalidValue)
	if err == nil {
		t.Error("convertToString() expected error for channel, got nil")
	}
}

