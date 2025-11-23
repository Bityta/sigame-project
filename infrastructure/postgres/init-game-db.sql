-- =====================================================
-- SIGame 2.0 - Game Service Database Schema
-- =====================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- GAME SESSIONS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS game_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL,
    pack_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'waiting',
    current_round INT NOT NULL DEFAULT 1,
    current_phase VARCHAR(50) NOT NULL DEFAULT 'waiting',
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_game_sessions_room_id ON game_sessions(room_id);
CREATE INDEX IF NOT EXISTS idx_game_sessions_status ON game_sessions(status);
CREATE INDEX IF NOT EXISTS idx_game_sessions_created_at ON game_sessions(created_at);

-- =====================================================
-- GAME PLAYERS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS game_players (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    username VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'player',
    score INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_game_players_game_id ON game_players(game_id);
CREATE INDEX IF NOT EXISTS idx_game_players_user_id ON game_players(user_id);
CREATE INDEX IF NOT EXISTS idx_game_players_game_user ON game_players(game_id, user_id);

-- =====================================================
-- GAME ROUNDS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS game_rounds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    round_number INT NOT NULL,
    round_name VARCHAR(255),
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_game_rounds_game_id ON game_rounds(game_id);
CREATE INDEX IF NOT EXISTS idx_game_rounds_game_round ON game_rounds(game_id, round_number);

-- =====================================================
-- GAME EVENTS TABLE (для аудита)
-- =====================================================
CREATE TABLE IF NOT EXISTS game_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    user_id UUID,
    round_number INT,
    question_id VARCHAR(100),
    data JSONB,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_game_events_game_id ON game_events(game_id);
CREATE INDEX IF NOT EXISTS idx_game_events_event_type ON game_events(event_type);
CREATE INDEX IF NOT EXISTS idx_game_events_timestamp ON game_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_game_events_game_timestamp ON game_events(game_id, timestamp);

-- =====================================================
-- TRIGGER FOR UPDATED_AT
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_game_sessions_updated_at 
    BEFORE UPDATE ON game_sessions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- HELPER FUNCTIONS
-- =====================================================

-- Get game statistics
CREATE OR REPLACE FUNCTION get_game_stats(p_game_id UUID)
RETURNS TABLE (
    total_players INT,
    total_events INT,
    duration_minutes INT,
    winner_id UUID,
    winner_score INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*)::INT FROM game_players WHERE game_id = p_game_id),
        (SELECT COUNT(*)::INT FROM game_events WHERE game_id = p_game_id),
        (SELECT EXTRACT(EPOCH FROM (finished_at - started_at))::INT / 60 
         FROM game_sessions WHERE id = p_game_id),
        (SELECT user_id FROM game_players 
         WHERE game_id = p_game_id 
         ORDER BY score DESC LIMIT 1),
        (SELECT MAX(score)::INT FROM game_players WHERE game_id = p_game_id);
END;
$$ LANGUAGE plpgsql;

-- Cleanup old finished games (older than 30 days)
CREATE OR REPLACE FUNCTION cleanup_old_games()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM game_sessions 
    WHERE status = 'finished' 
    AND finished_at < CURRENT_TIMESTAMP - INTERVAL '30 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- COMPLETION MESSAGE
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE 'Game Service database schema initialized successfully';
END $$;

