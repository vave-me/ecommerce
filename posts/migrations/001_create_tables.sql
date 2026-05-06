-- +goose Up
-----------------------------------------------------------------------
-- 1) Create and enable PostGIS (if not already done)
-----------------------------------------------------------------------
CREATE
EXTENSION IF NOT EXISTS postgis;

-----------------------------------------------------------------------
-- 2) Create posts table (simplified)
-----------------------------------------------------------------------
CREATE TABLE posts
(
    id            TEXT        NOT NULL,
    user_id       TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    description   TEXT        NOT NULL,
    type_of_post     TEXT,
    user_type     TEXT,
    category_id   TEXT,
    category_slug TEXT,
    tags          TEXT,
    status        TEXT        NOT NULL DEFAULT 'published',
    thumbnail     TEXT,
    location      geography(Point, 4326),
    lat           DOUBLE PRECISION     DEFAULT 0.00,
    lng           DOUBLE PRECISION     DEFAULT 0.00,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    PRIMARY KEY (id)
);

-- triggers for created_at / updated_at
CREATE TRIGGER created_at_posts_trgr
    BEFORE UPDATE
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_posts_trgr
    BEFORE UPDATE
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- geospatial index
CREATE INDEX idx_posts_location
    ON posts
    USING GIST(location);

-- example partial index for soft-deletes:
CREATE INDEX idx_posts_deleted_at
    ON posts (deleted_at) WHERE deleted_at IS NOT NULL;

-- indexing user_id for quick filtering
CREATE INDEX idx_posts_user_id
    ON posts (user_id);

-----------------------------------------------------------------------
-- 3) Create events table (for EventStore)
-----------------------------------------------------------------------
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

-----------------------------------------------------------------------
-- 4) Create snapshots table (for EventStore snapshots)
-----------------------------------------------------------------------
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

-----------------------------------------------------------------------
-- 5) Create inbox/outbox for message-driven patterns
-----------------------------------------------------------------------
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

-- Example partial index for unpublished outbox
CREATE INDEX idx_outbox_unpublished
    ON outbox (published_at) WHERE published_at IS NULL;

-----------------------------------------------------------------------
-- 6) Create sagas table (for distributed transactions/workflows)
-----------------------------------------------------------------------
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

-----------------------------------------------------------------------
-- +goose Down
-----------------------------------------------------------------------
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;

-- Optionally drop extension if you only used it for this DB:
-- DROP EXTENSION IF EXISTS postgis;
