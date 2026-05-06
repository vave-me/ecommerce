-- +goose Up

CREATE TABLE baskets
(
    id            TEXT        NOT NULL,
    user_id       TEXT        NOT NULL,
    basket_status TEXT        NOT NULL DEFAULT 'open',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
-- Partial index for quick lookups of baskets with status open
CREATE INDEX basket_status_baskets_idx
    ON baskets (basket_status) WHERE basket_status='open';

-- Timestamps triggers for "baskets"
CREATE TRIGGER created_at_baskets_trgr
    BEFORE UPDATE
    ON baskets
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_baskets_trgr
    BEFORE UPDATE
    ON baskets
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE users_cache
(
    id         text        NOT NULL,
    email      text        NOT NULL,
    username   text        NOT NULL,
    first_name text        NOT NULL,
    last_name  text        NOT NULL,
    location   text        NOT NULL,
    enabled    bool        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_users_trgr
    BEFORE UPDATE
    ON users_cache
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_trgr
    BEFORE UPDATE
    ON users_cache
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products_cache
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    description       TEXT        NOT NULL,
    base_price        INT         NOT NULL CHECK (base_price >= 0),
    user_seller_id    TEXT        NOT NULL,
    category_id       TEXT        NOT NULL,
    brand             TEXT,
    condition         TEXT        NOT NULL,
    model             TEXT,
    tags              TEXT,
    manage_stock      BOOL,
    stock             INT                  DEFAULT 0 CHECK (stock >= 0),
    sku               TEXT,
    attributes        TEXT,
    weight            INT,
    height            INT,
    width             INT,
    depth             INT,
    status            TEXT,
    negotiable        BOOL                 DEFAULT FALSE,
    user_type       TEXT,
    middleman_service BOOL,
    shipping_cost     INT,
    has_variants      BOOL,
    options           TEXT,
    thumbnail         text,
    lat               DOUBLE PRECISION DEFAULT 0.00,
    lng               DOUBLE PRECISION DEFAULT 0.00,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_products_trgr
    BEFORE UPDATE
    ON products_cache
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
    BEFORE UPDATE
    ON products_cache
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE events
(
    stream_id      text        NOT NULL,
    stream_name    text        NOT NULL,
    stream_version int         NOT NULL,
    event_id       text        NOT NULL,
    event_name     text        NOT NULL,
    event_data     bytea       NOT NULL,
    occurred_at    timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE TABLE snapshots
(
    stream_id      text        NOT NULL,
    stream_name    text        NOT NULL,
    stream_version int         NOT NULL,
    snapshot_name  text        NOT NULL,
    snapshot_data  bytea       NOT NULL,
    updated_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (stream_id, stream_name)
);

CREATE TRIGGER updated_at_snapshots_trgr
    BEFORE UPDATE
    ON baskets.snapshots
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE inbox
(
    id          text        NOT NULL,
    name        text        NOT NULL,
    subject     text        NOT NULL,
    data        bytea       NOT NULL,
    metadata    bytea       NOT NULL,
    sent_at     timestamptz NOT NULL,
    received_at timestamptz NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE outbox
(
    id           text        NOT NULL,
    name         text        NOT NULL,
    subject      text        NOT NULL,
    data         bytea       NOT NULL,
    metadata     bytea       NOT NULL,
    sent_at      timestamptz NOT NULL,
    published_at timestamptz,
    PRIMARY KEY (id)
);
CREATE INDEX basket_unpublished_idx ON baskets.outbox (published_at) WHERE published_at IS NULL;
CREATE TABLE sagas
(
    id           text        NOT NULL,
    name         text        NOT NULL,
    data         bytea       NOT NULL,
    step         int         NOT NULL,
    done         bool        NOT NULL,
    compensating bool        NOT NULL,
    updated_at   timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, name)
);

CREATE TRIGGER updated_at_sagas_trgr
    BEFORE UPDATE
    ON sagas
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
DROP TABLE IF EXISTS users_cache;
DROP TABLE IF EXISTS products_cache;

DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;


