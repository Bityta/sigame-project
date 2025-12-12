package postgres

import (
	"strings"
	"testing"
)

func TestQueriesContainTableNames(t *testing.T) {
	tests := []struct {
		name  string
		query string
		table string
	}{
		{
			name:  "insert game session",
			query: queryInsertGameSession,
			table: tableGameSessions,
		},
		{
			name:  "update game session",
			query: queryUpdateGameSession,
			table: tableGameSessions,
		},
		{
			name:  "select game session",
			query: querySelectGameSession,
			table: tableGameSessions,
		},
		{
			name:  "insert game player",
			query: queryInsertGamePlayer,
			table: tableGamePlayers,
		},
		{
			name:  "update player score",
			query: queryUpdatePlayerScore,
			table: tableGamePlayers,
		},
		{
			name:  "select game players",
			query: querySelectGamePlayers,
			table: tableGamePlayers,
		},
		{
			name:  "insert game event",
			query: queryInsertGameEvent,
			table: tableGameEvents,
		},
		{
			name:  "select game events",
			query: querySelectGameEvents,
			table: tableGameEvents,
		},
		{
			name:  "select events by type",
			query: querySelectEventsByType,
			table: tableGameEvents,
		},
		{
			name:  "count game events",
			query: queryCountGameEvents,
			table: tableGameEvents,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(strings.ToLower(tt.query), strings.ToLower(tt.table)) {
				t.Errorf("query %s does not contain table name %s", tt.name, tt.table)
			}
		})
	}
}

func TestQueriesAreNotEmpty(t *testing.T) {
	queries := []struct {
		name  string
		query string
	}{
		{"insert game session", queryInsertGameSession},
		{"update game session", queryUpdateGameSession},
		{"select game session", querySelectGameSession},
		{"update game session final", queryUpdateGameSessionFinal},
		{"select games by room id", querySelectGamesByRoomID},
		{"select active game for user", querySelectActiveGameForUser},
		{"insert game player", queryInsertGamePlayer},
		{"update player score", queryUpdatePlayerScore},
		{"select game players", querySelectGamePlayers},
		{"update player final", queryUpdatePlayerFinal},
		{"insert game event", queryInsertGameEvent},
		{"select game events", querySelectGameEvents},
		{"select events by type", querySelectEventsByType},
		{"count game events", queryCountGameEvents},
	}

	for _, tt := range queries {
		t.Run(tt.name, func(t *testing.T) {
			if strings.TrimSpace(tt.query) == "" {
				t.Errorf("query %s is empty", tt.name)
			}
		})
	}
}

func TestTableConstants(t *testing.T) {
	tables := []string{
		tableGameSessions,
		tableGamePlayers,
		tableGameEvents,
	}

	for _, table := range tables {
		if table == "" {
			t.Errorf("table constant is empty")
		}
	}
}

func TestColumnConstants(t *testing.T) {
	columns := []string{
		colID,
		colRoomID,
		colPackID,
		colStatus,
		colCurrentRound,
		colCurrentPhase,
		colStartedAt,
		colFinishedAt,
		colCreatedAt,
		colUpdatedAt,
		colGameID,
		colUserID,
		colUsername,
		colRole,
		colScore,
		colIsActive,
		colJoinedAt,
		colLeftAt,
		colEventType,
		colRoundNumber,
		colQuestionID,
		colData,
		colTimestamp,
	}

	for _, col := range columns {
		if col == "" {
			t.Errorf("column constant is empty")
		}
	}
}

