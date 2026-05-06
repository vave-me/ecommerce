-- +goose Up
-----------------------------------------------------------------------
-- 1) Create and enable required extensions
-----------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-----------------------------------------------------------------------
-- 2) Create streams table
-----------------------------------------------------------------------
CREATE TABLE streams
(
    id                TEXT        NOT NULL,
    title             TEXT        NOT NULL,
    description       TEXT        NOT NULL,
    synopsis          TEXT,
    stream_type       TEXT        NOT NULL,
    status            TEXT        NOT NULL DEFAULT 'draft',
    
    -- Content Details
    stream_url        TEXT        NOT NULL,
    thumbnail_url     TEXT,
    trailer_url       TEXT,
    duration          INT         NOT NULL CHECK (duration > 0),
    release_date      TIMESTAMPTZ,
    content_rating    TEXT        NOT NULL,
    
    -- Technical Details
    available_qualities TEXT[], -- Array of quality options
    default_quality   TEXT,
    subtitles         JSONB       DEFAULT '[]'::JSONB,
    audio_tracks      JSONB       DEFAULT '[]'::JSONB,
    
    -- Access Control
    access_type       TEXT        NOT NULL,
    subscription_tiers TEXT[],
    rental_price      INT         DEFAULT 0 CHECK (rental_price >= 0),
    rental_duration   INT         DEFAULT 48, -- hours
    purchase_price    INT         DEFAULT 0 CHECK (purchase_price >= 0),
    ppv_price         INT         DEFAULT 0 CHECK (ppv_price >= 0),
    
    -- Metadata
    genre             TEXT[],
    tags              TEXT[],
    cast_members      JSONB       DEFAULT '[]'::JSONB,
    directors         TEXT[],
    producers         TEXT[],
    studio            TEXT,
    language          TEXT,
    country           TEXT,
    
    -- Analytics
    view_count        BIGINT      DEFAULT 0,
    like_count        BIGINT      DEFAULT 0,
    dislike_count     BIGINT      DEFAULT 0,
    average_rating    NUMERIC(3,2) DEFAULT 0.00,
    total_revenue     BIGINT      DEFAULT 0,
    
    -- Series Information
    series_id         TEXT,
    season_number     INT,
    episode_number    INT,
    
    -- Timestamps
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at      TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    
    PRIMARY KEY (id)
);

-- Indexes for performance
CREATE INDEX idx_streams_status ON streams (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_streams_stream_type ON streams (stream_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_streams_access_type ON streams (access_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_streams_genre ON streams USING GIN (genre) WHERE deleted_at IS NULL;
CREATE INDEX idx_streams_tags ON streams USING GIN (tags) WHERE deleted_at IS NULL;
CREATE INDEX idx_streams_series_id ON streams (series_id) WHERE series_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_streams_published_at ON streams (published_at DESC) WHERE status = 'published' AND deleted_at IS NULL;
CREATE INDEX idx_streams_view_count ON streams (view_count DESC) WHERE status = 'published' AND deleted_at IS NULL;
CREATE INDEX idx_streams_created_at ON streams (created_at DESC) WHERE deleted_at IS NULL;

-- Full text search index
CREATE INDEX idx_streams_search ON streams 
    USING GIN (to_tsvector('english', title || ' ' || COALESCE(description, '') || ' ' || COALESCE(synopsis, '')))
    WHERE deleted_at IS NULL;

-- Triggers
CREATE TRIGGER created_at_streams_trgr
    BEFORE UPDATE
    ON streams
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_streams_trgr
    BEFORE UPDATE
    ON streams
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-----------------------------------------------------------------------
-- 3) Create series table
-----------------------------------------------------------------------
CREATE TABLE series
(
    id                TEXT        NOT NULL,
    title             TEXT        NOT NULL,
    description       TEXT        NOT NULL,
    thumbnail_url     TEXT,
    genre             TEXT[],
    studio            TEXT,
    total_seasons     INT         DEFAULT 0,
    status            TEXT        DEFAULT 'active',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    
    PRIMARY KEY (id)
);

-- Indexes
CREATE INDEX idx_series_genre ON series USING GIN (genre) WHERE deleted_at IS NULL;
CREATE INDEX idx_series_studio ON series (studio) WHERE deleted_at IS NULL;
CREATE INDEX idx_series_status ON series (status) WHERE deleted_at IS NULL;

-- Triggers
CREATE TRIGGER created_at_series_trgr
    BEFORE UPDATE
    ON series
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_series_trgr
    BEFORE UPDATE
    ON series
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- Foreign key from streams to series
ALTER TABLE streams
    ADD CONSTRAINT fk_streams_series
        FOREIGN KEY (series_id)
            REFERENCES series (id) ON DELETE SET NULL;

-----------------------------------------------------------------------
-- 4) Create stream_access table for user access management
-----------------------------------------------------------------------
CREATE TABLE stream_access
(
    stream_id         TEXT        NOT NULL,
    user_id           TEXT        NOT NULL,
    access_type       TEXT        NOT NULL,
    granted_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at        TIMESTAMPTZ NOT NULL,
    last_watched_at   TIMESTAMPTZ,
    watch_progress    INT         DEFAULT 0, -- seconds
    completed         BOOL        DEFAULT FALSE,
    
    PRIMARY KEY (stream_id, user_id),
    FOREIGN KEY (stream_id) REFERENCES streams (id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_stream_access_user_id ON stream_access (user_id);
CREATE INDEX idx_stream_access_expires_at ON stream_access (expires_at);
CREATE INDEX idx_stream_access_last_watched ON stream_access (user_id, last_watched_at DESC) 
    WHERE last_watched_at IS NOT NULL AND completed = FALSE;

-----------------------------------------------------------------------
-- 5) Create stream_ratings table
-----------------------------------------------------------------------
CREATE TABLE stream_ratings
(
    stream_id         TEXT        NOT NULL,
    user_id           TEXT        NOT NULL,
    rating            INT         NOT NULL CHECK (rating >= 1 AND rating <= 5),
    is_like           BOOL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (stream_id, user_id),
    FOREIGN KEY (stream_id) REFERENCES streams (id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_stream_ratings_stream_id ON stream_ratings (stream_id);

-- Trigger for updated_at
CREATE TRIGGER updated_at_stream_ratings_trgr
    BEFORE UPDATE
    ON stream_ratings
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-----------------------------------------------------------------------
-- 6) Create watchlist table
-----------------------------------------------------------------------
CREATE TABLE watchlists
(
    user_id           TEXT        NOT NULL,
    stream_id         TEXT        NOT NULL,
    added_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (user_id, stream_id),
    FOREIGN KEY (stream_id) REFERENCES streams (id) ON DELETE CASCADE
);

-- Index
CREATE INDEX idx_watchlists_user_added ON watchlists (user_id, added_at DESC);

-----------------------------------------------------------------------
-- 7) Event-sourcing tables
-----------------------------------------------------------------------
CREATE TABLE events
(
    stream_id         TEXT        NOT NULL,
    stream_name       TEXT        NOT NULL,
    stream_version    INT         NOT NULL,
    event_id          TEXT        NOT NULL,
    event_name        TEXT        NOT NULL,
    event_data        BYTEA       NOT NULL,
    occurred_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE INDEX idx_events_occurred_at ON events (occurred_at DESC);

CREATE TABLE snapshots
(
    stream_id         TEXT        NOT NULL,
    stream_name       TEXT        NOT NULL,
    stream_version    INT         NOT NULL,
    snapshot_name     TEXT        NOT NULL,
    snapshot_data     BYTEA       NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name)
);

CREATE TRIGGER updated_at_snapshots_trgr
    BEFORE UPDATE
    ON snapshots
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    subject           TEXT        NOT NULL,
    data              BYTEA       NOT NULL,
    metadata          BYTEA       NOT NULL,
    sent_at           TIMESTAMPTZ NOT NULL,
    received_at       TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE outbox
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    subject           TEXT        NOT NULL,
    data              BYTEA       NOT NULL,
    metadata          BYTEA       NOT NULL,
    sent_at           TIMESTAMPTZ NOT NULL,
    published_at      TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE INDEX idx_outbox_unpublished
    ON outbox (published_at) WHERE published_at IS NULL;

CREATE TABLE sagas
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    data              BYTEA       NOT NULL,
    step              INT         NOT NULL,
    done              BOOL        NOT NULL,
    compensating      BOOL        NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, name)
);

CREATE TRIGGER updated_at_sagas_trgr
    BEFORE UPDATE
    ON sagas
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
-----------------------------------------------------------------------
-- Drop in reverse order
-----------------------------------------------------------------------
DROP TABLE IF EXISTS watchlists CASCADE;
DROP TABLE IF EXISTS stream_ratings CASCADE;
DROP TABLE IF EXISTS stream_access CASCADE;
DROP TABLE IF EXISTS streams CASCADE;
DROP TABLE IF EXISTS series CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;