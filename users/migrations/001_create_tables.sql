-- +goose Up
---------------------------------------------------------------------
-- 1) Enable PostGIS (if not already installed in the DB)
---------------------------------------------------------------------
CREATE
EXTENSION IF NOT EXISTS postgis;

---------------------------------------------------------------------
-- 2) Create users table
---------------------------------------------------------------------
CREATE TABLE users
(
    id         TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    username   TEXT,
    firstname  TEXT,
    lastname   TEXT,
    google_id  TEXT                 default '',
    bio        text                 default '',
    privacy    text                 default '',
    background text                 default '',
    lat        float,
    lng        float,
    location   GEOGRAPHY(Point, 4326),
    enabled    BOOL,
    thumbnail  TEXT,
    language   TEXT,
    role       TEXT        NOT NULL DEFAULT 'user',
    last_login TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);


-- Index for filtering enabled users
CREATE INDEX enabled_users_idx ON users (enabled) WHERE enabled;

-- Index on email for better query performance
CREATE INDEX idx_users_email ON users (email);

-- If you do distance queries or bounding queries, create a GIST index on `location`
CREATE INDEX idx_users_location
    ON users
    USING GIST(location);

-- triggers
CREATE TRIGGER created_at_users_trgr
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_users_trgr
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE tokens
(
    id         TEXT        NOT NULL,
    user_id    TEXT        NOT NULL,
    firstname  TEXT,
    lastname   TEXT,
    email      TEXT        NOT NULL,
    token      TEXT,
    enabled    BOOL,
    last_login TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);


---------------------------------------------------------------------
-- 3) Create event-sourcing tables (unchanged)
---------------------------------------------------------------------
CREATE TABLE events
(
    stream_id      TEXT        NOT NULL,
    stream_name    TEXT        NOT NULL,
    stream_version INT         NOT NULL,
    event_id       TEXT        NOT NULL,
    event_name     TEXT        NOT NULL,
    event_data     BYTEA       NOT NULL,
    occurred_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE TABLE snapshots
(
    stream_id      TEXT        NOT NULL,
    stream_name    TEXT        NOT NULL,
    stream_version INT         NOT NULL,
    snapshot_name  TEXT        NOT NULL,
    snapshot_data  BYTEA       NOT NULL,
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name)
);

CREATE TRIGGER updated_at_snapshots_trgr
    BEFORE UPDATE
    ON snapshots
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox
(
    id          TEXT        NOT NULL,
    name        TEXT        NOT NULL,
    subject     TEXT        NOT NULL,
    data        BYTEA       NOT NULL,
    metadata    BYTEA       NOT NULL,
    sent_at     TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE outbox
(
    id           TEXT        NOT NULL,
    name         TEXT        NOT NULL,
    subject      TEXT        NOT NULL,
    data         BYTEA       NOT NULL,
    metadata     BYTEA       NOT NULL,
    sent_at      TIMESTAMPTZ NOT NULL,
    published_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE INDEX users_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

CREATE TABLE sagas
(
    id           TEXT        NOT NULL,
    name         TEXT        NOT NULL,
    data         BYTEA       NOT NULL,
    step         INT         NOT NULL,
    done         BOOL        NOT NULL,
    compensating BOOL        NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, name)
);

CREATE TRIGGER updated_at_sagas_trgr
    BEFORE UPDATE
    ON sagas
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();


-- +goose Down
---------------------------------------------------------------------
-- Drop in reverse order
---------------------------------------------------------------------
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;

-- remove the users table
DROP TABLE IF EXISTS users;

---------------------------------------------------------------------
-- (Optional) drop extension if only used for this DB
-- DROP EXTENSION IF EXISTS postgis;
