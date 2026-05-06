-- +goose Up
-----------------------------------------------------------------------
-- 1) Create and enable required extensions
-----------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-----------------------------------------------------------------------
-- 2) Create matches table
-----------------------------------------------------------------------
CREATE TABLE matches
(
    id                TEXT        NOT NULL,
    
    -- Teams and Competition
    home_team_id      TEXT        NOT NULL,
    home_team_name    TEXT        NOT NULL,
    away_team_id      TEXT        NOT NULL,
    away_team_name    TEXT        NOT NULL,
    competition_id    TEXT        NOT NULL,
    competition_name  TEXT        NOT NULL,
    competition_type  TEXT        NOT NULL,
    season            TEXT,
    round             INT,
    stage             TEXT,
    
    -- Match Details
    match_date        TIMESTAMPTZ NOT NULL,
    status            TEXT        NOT NULL DEFAULT 'draft',
    
    -- Stadium Information
    stadium_id        TEXT        NOT NULL,
    stadium_name      TEXT        NOT NULL,
    stadium_city      TEXT        NOT NULL,
    stadium_country   TEXT        NOT NULL,
    stadium_capacity  INT         NOT NULL CHECK (stadium_capacity > 0),
    stadium_location  geography(Point, 4326),
    
    -- Capacity and Sales
    total_capacity    INT         NOT NULL,
    available_tickets INT         NOT NULL,
    sold_tickets      INT         DEFAULT 0,
    
    -- Sales Dates
    sales_start_date  TIMESTAMPTZ NOT NULL,
    sales_end_date    TIMESTAMPTZ NOT NULL,
    early_sales_date  TIMESTAMPTZ,
    
    -- Match Officials
    referee           TEXT,
    linesmen          TEXT[], -- Array of linesmen names
    var_official      TEXT,
    
    -- Match Conditions
    weather           TEXT,
    temperature       INT,
    attendance        INT,
    
    -- Media
    thumbnail_url     TEXT,
    banner_url        TEXT,
    
    -- Pricing
    dynamic_pricing   BOOL        DEFAULT FALSE,
    
    -- Timestamps
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at      TIMESTAMPTZ,
    deleted_at        TIMESTAMPTZ,
    
    PRIMARY KEY (id)
);

-- Indexes for performance
CREATE INDEX idx_matches_status ON matches (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_match_date ON matches (match_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_home_team ON matches (home_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_away_team ON matches (away_team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_competition ON matches (competition_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_stadium ON matches (stadium_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_sales_dates ON matches (sales_start_date, sales_end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_matches_published ON matches (published_at DESC) WHERE status = 'published' AND deleted_at IS NULL;

-- Geospatial index
CREATE INDEX idx_matches_stadium_location
    ON matches
    USING GIST(stadium_location);

-- Triggers
CREATE TRIGGER created_at_matches_trgr
    BEFORE UPDATE
    ON matches
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_matches_trgr
    BEFORE UPDATE
    ON matches
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-----------------------------------------------------------------------
-- 3) Create sectors table
-----------------------------------------------------------------------
CREATE TABLE sectors
(
    id                TEXT        NOT NULL,
    match_id          TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    level             INT         NOT NULL DEFAULT 0, -- 0=ground, 1=first tier, etc
    category          TEXT        NOT NULL, -- VIP, Premium, Standard, Away
    total_seats       INT         NOT NULL CHECK (total_seats > 0),
    available_seats   INT         NOT NULL,
    sold_seats        INT         DEFAULT 0,
    base_price        INT         NOT NULL CHECK (base_price >= 0),
    amenities         TEXT[],
    entrance_gates    TEXT[],
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (id),
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_sectors_match_id ON sectors (match_id);
CREATE INDEX idx_sectors_category ON sectors (category);
CREATE INDEX idx_sectors_available ON sectors (match_id, available_seats) WHERE available_seats > 0;

-- Triggers
CREATE TRIGGER updated_at_sectors_trgr
    BEFORE UPDATE
    ON sectors
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-----------------------------------------------------------------------
-- 4) Create rows table
-----------------------------------------------------------------------
CREATE TABLE rows
(
    id                TEXT        NOT NULL,
    sector_id         TEXT        NOT NULL,
    row_number        TEXT        NOT NULL,
    total_seats       INT         NOT NULL CHECK (total_seats > 0),
    
    PRIMARY KEY (id),
    FOREIGN KEY (sector_id) REFERENCES sectors (id) ON DELETE CASCADE
);

-- Index
CREATE INDEX idx_rows_sector_id ON rows (sector_id);

-----------------------------------------------------------------------
-- 5) Create seats table
-----------------------------------------------------------------------
CREATE TABLE seats
(
    id                TEXT        NOT NULL,
    row_id            TEXT        NOT NULL,
    seat_number       INT         NOT NULL,
    status            TEXT        NOT NULL DEFAULT 'available',
    category          TEXT,
    price             INT         NOT NULL CHECK (price >= 0),
    features          TEXT[],
    ticket_id         TEXT,
    reserved_until    TIMESTAMPTZ,
    
    PRIMARY KEY (id),
    FOREIGN KEY (row_id) REFERENCES rows (id) ON DELETE CASCADE,
    UNIQUE (row_id, seat_number)
);

-- Indexes
CREATE INDEX idx_seats_row_id ON seats (row_id);
CREATE INDEX idx_seats_status ON seats (status);
CREATE INDEX idx_seats_reserved ON seats (reserved_until) WHERE reserved_until IS NOT NULL;

-----------------------------------------------------------------------
-- 6) Create tickets table
-----------------------------------------------------------------------
CREATE TABLE tickets
(
    id                TEXT        NOT NULL,
    
    -- Match Information
    match_id          TEXT        NOT NULL,
    match_date        TIMESTAMPTZ NOT NULL,
    home_team         TEXT        NOT NULL,
    away_team         TEXT        NOT NULL,
    competition       TEXT        NOT NULL,
    
    -- Seat Information
    stadium_name      TEXT        NOT NULL,
    sector_id         TEXT        NOT NULL,
    sector_name       TEXT        NOT NULL,
    row_id            TEXT        NOT NULL,
    row_number        TEXT        NOT NULL,
    seat_number       INT         NOT NULL,
    entrance_gate     TEXT        NOT NULL,
    
    -- Ticket Details
    type              TEXT        NOT NULL, -- regular, season, vip, etc
    status            TEXT        NOT NULL DEFAULT 'active',
    category          TEXT        NOT NULL,
    price             INT         NOT NULL CHECK (price >= 0),
    
    -- Owner Information
    owner_id          TEXT        NOT NULL,
    owner_name        TEXT        NOT NULL,
    owner_email       TEXT        NOT NULL,
    owner_phone       TEXT,
    
    -- Purchase Information
    purchaser_id      TEXT        NOT NULL,
    purchase_date     TIMESTAMPTZ NOT NULL,
    payment_id        TEXT        NOT NULL,
    order_id          TEXT        NOT NULL,
    
    -- Security
    qr_code           TEXT        NOT NULL UNIQUE,
    barcode           TEXT        NOT NULL UNIQUE,
    security_code     TEXT        NOT NULL,
    
    -- Transfer Information
    transferable      BOOL        DEFAULT TRUE,
    transfer_count    INT         DEFAULT 0,
    
    -- Usage Information
    used_at           TIMESTAMPTZ,
    used_gate         TEXT,
    
    -- Timestamps
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at        TIMESTAMPTZ NOT NULL,
    
    PRIMARY KEY (id),
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE RESTRICT,
    FOREIGN KEY (sector_id) REFERENCES sectors (id) ON DELETE RESTRICT
);

-- Indexes
CREATE INDEX idx_tickets_match_id ON tickets (match_id);
CREATE INDEX idx_tickets_owner_id ON tickets (owner_id);
CREATE INDEX idx_tickets_purchaser_id ON tickets (purchaser_id);
CREATE INDEX idx_tickets_order_id ON tickets (order_id);
CREATE INDEX idx_tickets_status ON tickets (status);
CREATE INDEX idx_tickets_qr_code ON tickets (qr_code);
CREATE INDEX idx_tickets_barcode ON tickets (barcode);
CREATE INDEX idx_tickets_match_date ON tickets (match_date);
CREATE INDEX idx_tickets_used ON tickets (used_at) WHERE used_at IS NOT NULL;

-- Triggers
CREATE TRIGGER created_at_tickets_trgr
    BEFORE UPDATE
    ON tickets
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_tickets_trgr
    BEFORE UPDATE
    ON tickets
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-----------------------------------------------------------------------
-- 7) Create ticket_transfers table
-----------------------------------------------------------------------
CREATE TABLE ticket_transfers
(
    id                TEXT        NOT NULL,
    ticket_id         TEXT        NOT NULL,
    from_user_id      TEXT        NOT NULL,
    to_user_id        TEXT        NOT NULL,
    reason            TEXT,
    transferred_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (id),
    FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE
);

-- Index
CREATE INDEX idx_ticket_transfers_ticket_id ON ticket_transfers (ticket_id);
CREATE INDEX idx_ticket_transfers_from_user ON ticket_transfers (from_user_id);
CREATE INDEX idx_ticket_transfers_to_user ON ticket_transfers (to_user_id);

-----------------------------------------------------------------------
-- 8) Create ticket_validations table
-----------------------------------------------------------------------
CREATE TABLE ticket_validations
(
    id                TEXT        NOT NULL,
    ticket_id         TEXT        NOT NULL,
    gate              TEXT        NOT NULL,
    result            TEXT        NOT NULL,
    reason            TEXT,
    validated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (id),
    FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE
);

-- Index
CREATE INDEX idx_ticket_validations_ticket_id ON ticket_validations (ticket_id);
CREATE INDEX idx_ticket_validations_validated_at ON ticket_validations (validated_at DESC);

-----------------------------------------------------------------------
-- 9) Create price_categories table
-----------------------------------------------------------------------
CREATE TABLE price_categories
(
    id                TEXT        NOT NULL,
    match_id          TEXT        NOT NULL,
    category          TEXT        NOT NULL,
    base_price        INT         NOT NULL CHECK (base_price >= 0),
    
    PRIMARY KEY (id),
    FOREIGN KEY (match_id) REFERENCES matches (id) ON DELETE CASCADE,
    UNIQUE (match_id, category)
);

-- Index
CREATE INDEX idx_price_categories_match_id ON price_categories (match_id);

-----------------------------------------------------------------------
-- 10) Event-sourcing tables
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
DROP TABLE IF EXISTS ticket_validations CASCADE;
DROP TABLE IF EXISTS ticket_transfers CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS seats CASCADE;
DROP TABLE IF EXISTS rows CASCADE;
DROP TABLE IF EXISTS sectors CASCADE;
DROP TABLE IF EXISTS price_categories CASCADE;
DROP TABLE IF EXISTS matches CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;