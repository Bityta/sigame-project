package domain

import (
	"testing"
	"time"

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
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty username",
			user: User{
				ID:       uuid.New(),
				Username: "",
				Email:    "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			user: User{
				ID:       uuid.New(),
				Username: "testuser",
				Email:    "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsEmailVerified(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "email verified",
			user: User{
				EmailVerifiedAt: &now,
			},
			want: true,
		},
		{
			name: "email not verified",
			user: User{
				EmailVerifiedAt: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsEmailVerified(); got != tt.want {
				t.Errorf("User.IsEmailVerified() = %v, want %v", got, tt.want)
			}
		})
	}
}

