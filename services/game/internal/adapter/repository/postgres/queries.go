package postgres

const (
	tableGameSessions = "game_sessions"
	tableGamePlayers  = "game_players"
	tableGameEvents   = "game_events"
)

const (
	colID            = "id"
	colRoomID       = "room_id"
	colPackID       = "pack_id"
	colStatus       = "status"
	colCurrentRound = "current_round"
	colCurrentPhase = "current_phase"
	colStartedAt    = "started_at"
	colFinishedAt   = "finished_at"
	colCreatedAt    = "created_at"
	colUpdatedAt    = "updated_at"
	colGameID       = "game_id"
	colUserID       = "user_id"
	colUsername     = "username"
	colRole         = "role"
	colScore        = "score"
	colIsActive     = "is_active"
	colJoinedAt     = "joined_at"
	colLeftAt       = "left_at"
	colEventType    = "event_type"
	colRoundNumber  = "round_number"
	colQuestionID   = "question_id"
	colData         = "data"
	colTimestamp    = "timestamp"
)

const (
	queryInsertGameSession = `
		INSERT INTO game_sessions (id, room_id, pack_id, status, current_round, current_phase, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	queryUpdateGameSession = `
		UPDATE game_sessions 
		SET status = $1, current_round = $2, current_phase = $3, 
		    started_at = $4, finished_at = $5, updated_at = $6
		WHERE id = $7
	`

	querySelectGameSession = `
		SELECT id, room_id, pack_id, status, current_round, current_phase, 
		       started_at, finished_at, created_at, updated_at
		FROM game_sessions
		WHERE id = $1
	`

	queryUpdateGameSessionFinal = `
		UPDATE game_sessions 
		SET status = $1, finished_at = $2, updated_at = $3
		WHERE id = $4
	`

	querySelectGamesByRoomID = `
		SELECT id, room_id, pack_id, status, current_round, current_phase, 
		       started_at, finished_at, created_at, updated_at
		FROM game_sessions
		WHERE room_id = $1
		ORDER BY created_at DESC
	`

	querySelectActiveGameForUser = `
		SELECT gs.id, gs.room_id, gs.pack_id, gs.status, gs.current_round, gs.current_phase, 
		       gs.started_at, gs.finished_at, gs.created_at, gs.updated_at
		FROM game_sessions gs
		INNER JOIN game_players gp ON gs.id = gp.game_id
		WHERE gp.user_id = $1 
		  AND gp.is_active = true
		  AND gs.status NOT IN ('finished', 'cancelled', 'game_end')
		ORDER BY gs.created_at DESC
		LIMIT 1
	`
)

const (
	queryInsertGamePlayer = `
		INSERT INTO game_players (game_id, user_id, username, role, score, is_active, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	queryUpdatePlayerScore = `
		UPDATE game_players 
		SET score = $1, is_active = $2
		WHERE game_id = $3 AND user_id = $4
	`

	querySelectGamePlayers = `
		SELECT user_id, username, role, score, is_active, joined_at, left_at
		FROM game_players
		WHERE game_id = $1
		ORDER BY joined_at
	`

	queryUpdatePlayerFinal = `
		UPDATE game_players 
		SET score = $1, is_active = $2, left_at = $3
		WHERE game_id = $4 AND user_id = $5
	`
)

const (
	queryInsertGameEvent = `
		INSERT INTO game_events (id, game_id, event_type, user_id, round_number, question_id, data, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	querySelectGameEvents = `
		SELECT id, game_id, event_type, user_id, round_number, question_id, data, timestamp
		FROM game_events
		WHERE game_id = $1
		ORDER BY timestamp ASC
	`

	querySelectEventsByType = `
		SELECT id, game_id, event_type, user_id, round_number, question_id, data, timestamp
		FROM game_events
		WHERE game_id = $1 AND event_type = $2
		ORDER BY timestamp ASC
	`

	queryCountGameEvents = `
		SELECT COUNT(*) FROM game_events WHERE game_id = $1
	`
)

