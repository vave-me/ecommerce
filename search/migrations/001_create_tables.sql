-- +goose Up
CREATE TABLE users_cache
(
    id         text        NOT NULL,
    first_name text        NOT NULL,
    last_name  text        NOT NULL,
    username   text,
    location   text,
    email      text        NOT NULL,
    enabled    bool,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_users_trgr
    BEFORE UPDATE
    ON users_cache
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_trgr
    BEFORE UPDATE
    ON users_cache
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products_cache
(
    id             text          NOT NULL,
    user_seller_id text          NOT NULL,
    name           text          NOT NULL,
    description    text          NOT NULL,
    price          NUMERIC(9, 2) NOT NULL CHECK (price >= 0),
    created_at     timestamptz   NOT NULL DEFAULT NOW(),
    updated_at     timestamptz   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_products_trgr
    BEFORE UPDATE
    ON products_cache
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
    BEFORE UPDATE
    ON products_cache
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE orders
(
    order_id           text        NOT NULL,
    user_customer_id   text        NOT NULL,
    user_customer_name text        NOT NULL,
    items              bytea       NOT NULL,
    status             text        NOT NULL,
    product_ids        text ARRAY  NOT NULL,
    user_sellers       text ARRAY  NOT NULL,
    created_at         timestamptz NOT NULL DEFAULT NOW(),
    updated_at         timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (order_id)
);

CREATE TRIGGER updated_at_sorders_trgr
    BEFORE UPDATE
    ON orders
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

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
    ON snapshots
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

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

CREATE INDEX unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;

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
-- Add indexes for foreign key lookups
CREATE INDEX idx_products_cache_user_seller_id ON products_cache (user_seller_id);
CREATE INDEX idx_orders_user_customer_id ON orders (user_customer_id);

-- Add composite indexes for common query patterns
CREATE INDEX idx_orders_status_created_at ON orders (status, created_at DESC);
CREATE INDEX idx_products_cache_created_at ON products_cache (created_at DESC);
CREATE INDEX idx_products_cache_price ON products_cache (price);

-- Add indexes for array columns (using GIN for better performance)
CREATE INDEX idx_orders_product_ids ON orders USING GIN(product_ids);
CREATE INDEX idx_orders_user_sellers ON orders USING GIN(user_sellers);

-- Add index for location searches
CREATE INDEX idx_users_cache_location ON users_cache (location);

-- Add indexes for event sourcing tables
CREATE INDEX idx_events_occurred_at ON events (occurred_at DESC);
CREATE INDEX idx_events_stream_name ON events (stream_name);
CREATE INDEX idx_outbox_sent_at ON outbox (sent_at) WHERE published_at IS NULL;


-- +goose Down
DROP TABLE IF EXISTS users_cache;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products_cache;

DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
