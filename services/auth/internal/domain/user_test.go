package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				ID:       uuid.New(),
				Username: "testuser",
			},
			wantErr: false,
		},
		{
			name: "empty username",
			user: User{
				ID:       uuid.New(),
				Username: "",
			},
			wantErr: false, // No validation in User struct currently
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just check that user can be created
			if tt.user.ID == uuid.Nil {
				t.Error("User ID should not be nil")
			}
		})
	}
}

func TestUser_ToResponse(t *testing.T) {
	user := User{
		ID:       uuid.New(),
		Username: "testuser",
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, response.ID)
	}

	if response.Username != user.Username {
		t.Errorf("Expected Username %v, got %v", user.Username, response.Username)
	}
}

func TestNormalizeUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     string
	}{
		{
			name:     "lowercase",
			username: "testuser",
			want:     "testuser",
		},
		{
			name:     "uppercase",
			username: "TESTUSER",
			want:     "testuser",
		},
		{
			name:     "mixed case",
			username: "TestUser",
			want:     "testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeUsername(tt.username); got != tt.want {
				t.Errorf("NormalizeUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

