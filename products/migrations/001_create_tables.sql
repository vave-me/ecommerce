-- +goose Up
-----------------------------------------------------------------------
-- 1) Create and enable PostGIS (if not already done)
-----------------------------------------------------------------------
CREATE
EXTENSION IF NOT EXISTS postgis;

-----------------------------------------------------------------------
-- 2) Create products table
-----------------------------------------------------------------------
CREATE TABLE products
(
    id                TEXT        NOT NULL,
    name              TEXT        NOT NULL,
    description       TEXT        NOT NULL,
    base_price        INT         NOT NULL CHECK (base_price >= 0),
    user_seller_id    TEXT        NOT NULL,
    category_id       TEXT        NOT NULL,
    category_slug     TEXT        NOT NULL,
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
    location          geography(Point, 4326),
    lat               DOUBLE PRECISION     DEFAULT 0.00,
    lng               DOUBLE PRECISION     DEFAULT 0.00,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,

    PRIMARY KEY (id)
);

-- triggers for created_at / updated_at
CREATE TRIGGER created_at_products_trgr
    BEFORE UPDATE
    ON products
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_products_trgr
    BEFORE UPDATE
    ON products
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- geospatial index
CREATE INDEX idx_products_location
    ON products
    USING GIST(location);

-- example partial index for soft-deletes:
CREATE INDEX idx_products_deleted_at
    ON products (deleted_at) WHERE deleted_at IS NOT NULL;

-- possible additional indexes for filtering
CREATE INDEX idx_products_user_seller_id ON products (user_seller_id);
CREATE INDEX idx_products_category_id ON products (category_id);

-----------------------------------------------------------------------
-- 3) Create variants table
-----------------------------------------------------------------------
CREATE TABLE variants
(
    id            TEXT        NOT NULL,
    product_id    TEXT        NOT NULL,
    status        TEXT,
    sku           TEXT,
    barcode       TEXT,
    condition     TEXT,
    variant_price INT         NOT NULL CHECK (variant_price >= 0),
    currency_code TEXT,
    stock         INT,
    weight        INT,
    height        INT,
    width         INT,
    depth         INT,
    attributes    TEXT, -- store JSON if you want
    is_available  BOOL,
    has_options   BOOL,
    options       TEXT, -- store JSON
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,

    PRIMARY KEY (id)
);

ALTER TABLE variants
    ADD CONSTRAINT fk_variants_product
        FOREIGN KEY (product_id)
            REFERENCES products (id) ON DELETE CASCADE;

-- triggers
CREATE TRIGGER created_at_variants_trgr
    BEFORE UPDATE
    ON variants
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_variants_trgr
    BEFORE UPDATE
    ON variants
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- partial index for soft-deletes
CREATE INDEX idx_variants_deleted_at
    ON variants (deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX idx_variants_product_id ON variants (product_id);

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
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;

-- Optionally drop extension if you only used it for this DB:
-- DROP EXTENSION IF EXISTS postgis;
