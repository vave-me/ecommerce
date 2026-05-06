-- +goose Up


CREATE TABLE categories
(
    id                 SERIAL PRIMARY KEY,
    category_id        TEXT        NOT NULL,
    description        TEXT        NOT NULL,
    parent_id          TEXT,
    google_category_id TEXT,
    tags               TEXT,
    is_active          BOOL                 DEFAULT TRUE,
    slug               TEXT,
    seo_title          TEXT,
    seo_keywords       TEXT,
    seo_desc           TEXT,
    lang               TEXT,
    category_type      TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);
CREATE INDEX active_categories_idx ON categories (is_active) WHERE is_active;

CREATE TRIGGER created_at_categories_trgr
    BEFORE UPDATE
    ON categories
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_categories_trgr
    BEFORE UPDATE
    ON categories
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE filters
(
    id          TEXT        NOT NULL,
    category_id TEXT        NOT NULL,
    name        TEXT,
    filter_type TEXT,
    values      TEXT,
    is_active   BOOL                 DEFAULT TRUE,
    lang        text,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE INDEX categories_filters_idx ON filters (category_id);

CREATE TRIGGER created_at_filters_trgr
    BEFORE UPDATE
    ON filters
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_filters_trgr
    BEFORE UPDATE
    ON filters
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

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

-------------------------------------------------------------------------------
-- +goose Down
-------------------------------------------------------------------------------
DROP TABLE IF EXISTS filters;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS item_interaction_counts;
