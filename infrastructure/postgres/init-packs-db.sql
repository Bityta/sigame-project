-- =====================================================
-- SIGame 2.0 - Pack Management Service Database Schema
-- =====================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For text search

-- =====================================================
-- ENUMS
-- =====================================================
CREATE TYPE pack_status AS ENUM ('processing', 'approved', 'rejected');
CREATE TYPE media_type AS ENUM ('text', 'image', 'audio', 'video');

-- =====================================================
-- PACKS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS packs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    description TEXT,
    s3_path VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    uploaded_by UUID NOT NULL,
    downloads_count INTEGER NOT NULL DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.00 CHECK (rating BETWEEN 0.00 AND 5.00),
    rating_count INTEGER NOT NULL DEFAULT 0,
    status pack_status NOT NULL DEFAULT 'processing',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for packs
CREATE INDEX IF NOT EXISTS idx_packs_name ON packs USING gin(name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_packs_author ON packs(author);
CREATE INDEX IF NOT EXISTS idx_packs_status ON packs(status);
CREATE INDEX IF NOT EXISTS idx_packs_rating ON packs(rating DESC);
CREATE INDEX IF NOT EXISTS idx_packs_created_at ON packs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_packs_uploaded_by ON packs(uploaded_by);

-- =====================================================
-- PACK ROUNDS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS pack_rounds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pack_id UUID NOT NULL REFERENCES packs(id) ON DELETE CASCADE,
    round_number INTEGER NOT NULL,
    round_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pack_id, round_number)
);

-- Create indexes for pack rounds
CREATE INDEX IF NOT EXISTS idx_pack_rounds_pack_id ON pack_rounds(pack_id);

-- =====================================================
-- PACK THEMES TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS pack_themes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    round_id UUID NOT NULL REFERENCES pack_rounds(id) ON DELETE CASCADE,
    theme_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for pack themes
CREATE INDEX IF NOT EXISTS idx_pack_themes_round_id ON pack_themes(round_id);

-- =====================================================
-- PACK QUESTIONS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS pack_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    theme_id UUID NOT NULL REFERENCES pack_themes(id) ON DELETE CASCADE,
    price INTEGER NOT NULL CHECK (price > 0),
    question_text TEXT NOT NULL,
    answer_text TEXT NOT NULL,
    media_type media_type NOT NULL DEFAULT 'text',
    media_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for pack questions
CREATE INDEX IF NOT EXISTS idx_pack_questions_theme_id ON pack_questions(theme_id);
CREATE INDEX IF NOT EXISTS idx_pack_questions_price ON pack_questions(price);

-- =====================================================
-- PACK TAGS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS pack_tags (
    pack_id UUID NOT NULL REFERENCES packs(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (pack_id, tag)
);

-- Create indexes for pack tags
CREATE INDEX IF NOT EXISTS idx_pack_tags_tag ON pack_tags(tag);

-- =====================================================
-- PACK RATINGS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS pack_ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pack_id UUID NOT NULL REFERENCES packs(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review_text TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pack_id, user_id)
);

-- Create indexes for pack ratings
CREATE INDEX IF NOT EXISTS idx_pack_ratings_pack_id ON pack_ratings(pack_id);
CREATE INDEX IF NOT EXISTS idx_pack_ratings_user_id ON pack_ratings(user_id);

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

CREATE TRIGGER update_packs_updated_at 
    BEFORE UPDATE ON packs 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- TRIGGER TO UPDATE PACK RATING
-- =====================================================
CREATE OR REPLACE FUNCTION update_pack_rating()
RETURNS TRIGGER AS $$
BEGIN
    -- Recalculate average rating and count
    UPDATE packs
    SET 
        rating = (SELECT AVG(rating)::DECIMAL(3,2) FROM pack_ratings WHERE pack_id = NEW.pack_id),
        rating_count = (SELECT COUNT(*) FROM pack_ratings WHERE pack_id = NEW.pack_id)
    WHERE id = NEW.pack_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_pack_rating
    AFTER INSERT OR UPDATE ON pack_ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_pack_rating();

-- Also handle rating deletion
CREATE OR REPLACE FUNCTION update_pack_rating_on_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE packs
    SET 
        rating = COALESCE((SELECT AVG(rating)::DECIMAL(3,2) FROM pack_ratings WHERE pack_id = OLD.pack_id), 0.00),
        rating_count = (SELECT COUNT(*) FROM pack_ratings WHERE pack_id = OLD.pack_id)
    WHERE id = OLD.pack_id;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_pack_rating_on_delete
    AFTER DELETE ON pack_ratings
    FOR EACH ROW
    EXECUTE FUNCTION update_pack_rating_on_delete();

-- =====================================================
-- HELPER FUNCTIONS
-- =====================================================

-- Function to get pack statistics
CREATE OR REPLACE FUNCTION get_pack_stats(p_pack_id UUID)
RETURNS TABLE (
    total_rounds INTEGER,
    total_themes INTEGER,
    total_questions INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(DISTINCT pr.id)::INTEGER as total_rounds,
        COUNT(DISTINCT pt.id)::INTEGER as total_themes,
        COUNT(pq.id)::INTEGER as total_questions
    FROM packs p
    LEFT JOIN pack_rounds pr ON p.id = pr.pack_id
    LEFT JOIN pack_themes pt ON pr.id = pt.round_id
    LEFT JOIN pack_questions pq ON pt.id = pq.theme_id
    WHERE p.id = p_pack_id;
END;
$$ LANGUAGE plpgsql;

-- Function to increment download count
CREATE OR REPLACE FUNCTION increment_download_count(p_pack_id UUID)
RETURNS VOID AS $$
BEGIN
    UPDATE packs
    SET downloads_count = downloads_count + 1
    WHERE id = p_pack_id;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- VIEW: Pack Details with Statistics
-- =====================================================
CREATE OR REPLACE VIEW v_pack_details AS
SELECT 
    p.*,
    COUNT(DISTINCT pr.id) as rounds_count,
    COUNT(DISTINCT pt.id) as themes_count,
    COUNT(pq.id) as questions_count,
    array_agg(DISTINCT ptags.tag) FILTER (WHERE ptags.tag IS NOT NULL) as tags
FROM packs p
LEFT JOIN pack_rounds pr ON p.id = pr.pack_id
LEFT JOIN pack_themes pt ON pr.id = pt.round_id
LEFT JOIN pack_questions pq ON pt.id = pq.theme_id
LEFT JOIN pack_tags ptags ON p.id = ptags.pack_id
GROUP BY p.id;

-- =====================================================
-- SEARCH FUNCTION
-- =====================================================
CREATE OR REPLACE FUNCTION search_packs(
    search_query TEXT DEFAULT NULL,
    search_author TEXT DEFAULT NULL,
    search_tags TEXT[] DEFAULT NULL,
    min_rating DECIMAL DEFAULT 0.00,
    limit_count INTEGER DEFAULT 20,
    offset_count INTEGER DEFAULT 0
)
RETURNS TABLE (
    id UUID,
    name VARCHAR,
    author VARCHAR,
    description TEXT,
    thumbnail_url VARCHAR,
    rating DECIMAL,
    rating_count INTEGER,
    downloads_count INTEGER,
    rounds_count BIGINT,
    questions_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.name,
        p.author,
        p.description,
        p.thumbnail_url,
        p.rating,
        p.rating_count,
        p.downloads_count,
        COUNT(DISTINCT pr.id) as rounds_count,
        COUNT(pq.id) as questions_count
    FROM packs p
    LEFT JOIN pack_rounds pr ON p.id = pr.pack_id
    LEFT JOIN pack_themes pt ON pr.id = pt.round_id
    LEFT JOIN pack_questions pq ON pt.id = pq.theme_id
    LEFT JOIN pack_tags ptags ON p.id = ptags.pack_id
    WHERE 
        p.status = 'approved'
        AND (search_query IS NULL OR p.name ILIKE '%' || search_query || '%')
        AND (search_author IS NULL OR p.author ILIKE '%' || search_author || '%')
        AND (search_tags IS NULL OR ptags.tag = ANY(search_tags))
        AND p.rating >= min_rating
    GROUP BY p.id
    ORDER BY p.rating DESC, p.created_at DESC
    LIMIT limit_count
    OFFSET offset_count;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- COMPLETION MESSAGE
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE 'Pack Management database schema initialized successfully';
END $$;

