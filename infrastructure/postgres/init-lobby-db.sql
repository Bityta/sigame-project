-- =====================================================
-- SIGame 2.0 - Lobby Service Database Schema
-- =====================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- NOTE: Using VARCHAR for status/role instead of ENUM
-- for better compatibility with Spring Data R2DBC
-- =====================================================

-- =====================================================
-- GAME ROOMS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS game_rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_code VARCHAR(6) NOT NULL UNIQUE,
    host_id UUID NOT NULL,
    pack_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'waiting',
    max_players INTEGER NOT NULL DEFAULT 6 CHECK (max_players BETWEEN 2 AND 12),
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    password_hash VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    finished_at TIMESTAMP
);

-- Create indexes for game rooms
CREATE INDEX IF NOT EXISTS idx_game_rooms_room_code ON game_rooms(room_code);
CREATE INDEX IF NOT EXISTS idx_game_rooms_status ON game_rooms(status);
CREATE INDEX IF NOT EXISTS idx_game_rooms_host_id ON game_rooms(host_id);
CREATE INDEX IF NOT EXISTS idx_game_rooms_created_at ON game_rooms(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_game_rooms_is_public ON game_rooms(is_public) WHERE is_public = TRUE;

-- =====================================================
-- ROOM PLAYERS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS room_players (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID NOT NULL REFERENCES game_rooms(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    username VARCHAR(50) NOT NULL,
    avatar_url VARCHAR(500),
    role VARCHAR(20) NOT NULL DEFAULT 'player',
    is_ready BOOLEAN NOT NULL DEFAULT FALSE,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    UNIQUE(room_id, user_id)
);

-- Create indexes for room players
CREATE INDEX IF NOT EXISTS idx_room_players_room_id ON room_players(room_id);
CREATE INDEX IF NOT EXISTS idx_room_players_user_id ON room_players(user_id);
CREATE INDEX IF NOT EXISTS idx_room_players_active ON room_players(room_id, user_id) WHERE left_at IS NULL;

-- =====================================================
-- ROOM SETTINGS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS room_settings (
    room_id UUID PRIMARY KEY REFERENCES game_rooms(id) ON DELETE CASCADE,
    time_for_answer INTEGER NOT NULL DEFAULT 30 CHECK (time_for_answer BETWEEN 10 AND 120),
    time_for_choice INTEGER NOT NULL DEFAULT 60 CHECK (time_for_choice BETWEEN 10 AND 180),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- TRIGGERS FOR UPDATED_AT
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_game_rooms_updated_at 
    BEFORE UPDATE ON game_rooms 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_room_settings_updated_at 
    BEFORE UPDATE ON room_settings 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- HELPER FUNCTIONS
-- =====================================================

-- Function to generate unique 6-character room code
CREATE OR REPLACE FUNCTION generate_room_code()
RETURNS VARCHAR(6) AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result VARCHAR(6) := '';
    i INTEGER;
BEGIN
    FOR i IN 1..6 LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::int, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Function to get active player count for a room
CREATE OR REPLACE FUNCTION get_active_player_count(p_room_id UUID)
RETURNS INTEGER AS $$
BEGIN
    RETURN (
        SELECT COUNT(*)
        FROM room_players
        WHERE room_id = p_room_id AND left_at IS NULL
    );
END;
$$ LANGUAGE plpgsql;

-- Function to check if room has available slots
CREATE OR REPLACE FUNCTION room_has_slots(p_room_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    v_current_count INTEGER;
    v_max_players INTEGER;
BEGIN
    SELECT get_active_player_count(p_room_id), max_players
    INTO v_current_count, v_max_players
    FROM game_rooms
    WHERE id = p_room_id;
    
    RETURN v_current_count < v_max_players;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- VIEW: Active Rooms with Player Count
-- =====================================================
CREATE OR REPLACE VIEW v_active_rooms AS
SELECT 
    gr.id,
    gr.room_code,
    gr.name,
    gr.status,
    gr.max_players,
    gr.is_public,
    gr.created_at,
    COUNT(rp.id) FILTER (WHERE rp.left_at IS NULL) as current_players,
    gr.max_players - COUNT(rp.id) FILTER (WHERE rp.left_at IS NULL) as available_slots
FROM game_rooms gr
LEFT JOIN room_players rp ON gr.id = rp.room_id
WHERE gr.status IN ('waiting', 'starting')
GROUP BY gr.id;

-- =====================================================
-- CLEANUP FUNCTION
-- =====================================================
CREATE OR REPLACE FUNCTION cleanup_old_finished_rooms()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Delete rooms finished more than 7 days ago
    DELETE FROM game_rooms 
    WHERE status IN ('finished', 'cancelled') 
    AND finished_at < CURRENT_TIMESTAMP - INTERVAL '7 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- COMPLETION MESSAGE
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE 'Lobby database schema initialized successfully';
END $$;

