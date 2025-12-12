package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sigame/game/internal/domain/game"
	"github.com/sigame/game/internal/domain/pack"
	"github.com/sigame/game/internal/domain/player"
	"github.com/sigame/game/internal/domain/event"
)

func TestHandleNullTime(t *testing.T) {
	tests := []struct {
		name     string
		nullTime sql.NullTime
		want     *time.Time
	}{
		{
			name: "valid time",
			nullTime: sql.NullTime{
				Valid: true,
				Time:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
			want: func() *time.Time {
				t := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				return &t
			}(),
		},
		{
			name: "invalid time",
			nullTime: sql.NullTime{
				Valid: false,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleNullTime(tt.nullTime)
			if tt.want == nil {
				if got != nil {
					t.Errorf("handleNullTime() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("handleNullTime() = nil, want %v", tt.want)
				} else if !got.Equal(*tt.want) {
					t.Errorf("handleNullTime() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestMarshalEventData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
		},
		{
			name: "empty data",
			data: make(map[string]interface{}),
			wantErr: false,
		},
		{
			name: "valid data",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := marshalEventData(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshalEventData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.data == nil && got != nil {
				t.Errorf("marshalEventData() = %v, want nil for nil data", got)
			}
		})
	}
}

func TestUnmarshalEventData(t *testing.T) {
	tests := []struct {
		name     string
		dataJSON []byte
		wantErr  bool
	}{
		{
			name:     "empty data",
			dataJSON: []byte{},
			wantErr:  false,
		},
		{
			name:     "valid JSON",
			dataJSON: []byte(`{"key1":"value1","key2":123}`),
			wantErr:  false,
		},
		{
			name:     "invalid JSON",
			dataJSON: []byte(`{invalid json}`),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &domain.GameEvent{
				ID:        uuid.New(),
				GameID:    uuid.New(),
				EventType: domain.EventGameCreated,
			}
			err := unmarshalEventData(tt.dataJSON, event)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalEventData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(tt.dataJSON) > 0 && event.Data == nil {
				t.Errorf("unmarshalEventData() event.Data = nil, want non-nil")
			}
		})
	}
}


