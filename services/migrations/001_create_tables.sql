-- +goose Up
-----------------------------------------------------------------------
-- 1) Create and enable PostGIS (if not already done)
-----------------------------------------------------------------------
CREATE
EXTENSION IF NOT EXISTS postgis;

-----------------------------------------------------------------------
-- 2) Create services table
-----------------------------------------------------------------------
CREATE TABLE services
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    description       TEXT        NOT NULL,
    service_type      TEXT,
    base_price        INT         NOT NULL CHECK (base_price >= 0),
    pricing           TEXT,
    availability      TEXT,
    provider_name     TEXT,
    user_id           TEXT        NOT NULL,
    category_id       TEXT        NOT NULL,
    category_slug     TEXT,
    description_short TEXT,
    description_long  TEXT,
    qualifications    TEXT,
    contact           TEXT,
    faq               TEXT,
    tags              TEXT,
    status            TEXT,
    negotiable        BOOL                 DEFAULT FALSE,
    user_type         TEXT,
    middleman_service BOOL                 DEFAULT FALSE,
    has_variants      BOOL                 DEFAULT FALSE,
    attributes        TEXT,
    shipping_cost     INT,
    options           TEXT,
    thumbnail         text,
    location          geography(Point, 4326),
    lat               DOUBLE PRECISION     DEFAULT 0.00,
    lng               DOUBLE PRECISION     DEFAULT 0.00,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    PRIMARY KEY (id)
);

-- triggers for created_at / updated_at
CREATE TRIGGER created_at_services_trgr
    BEFORE UPDATE
    ON services
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_services_trgr
    BEFORE UPDATE
    ON services
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- geospatial index
CREATE INDEX idx_services_location
    ON services
    USING GIST(location);

-- example partial index for soft-deletes:
CREATE INDEX idx_services_deleted_at
    ON services (deleted_at) WHERE deleted_at IS NOT NULL;

-- possible additional indexes for filtering
CREATE INDEX idx_services_user_id ON services (user_id);
CREATE INDEX idx_services_category_id ON services (category_id);


-----------------------------------------------------------------------
-- 4) Event-sourcing / snapshots / sagas / outbox / inbox tables
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

-- Example partial index for unpublished outbox
CREATE INDEX idx_outbox_unpublished
    ON outbox (published_at) WHERE published_at IS NULL;

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
-----------------------------------------------------------------------
-- Drop in reverse order (CASCADE to remove foreign keys)
-----------------------------------------------------------------------
DROP TABLE IF EXISTS variants CASCADE;
DROP TABLE IF EXISTS services CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;

-- Optionally drop extension if you only used it for this DB:
-- DROP EXTENSION IF EXISTS postgis;
