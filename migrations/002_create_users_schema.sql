-- +goose Up
CREATE SCHEMA users;

SET
SEARCH_PATH TO users, public;

CREATE TABLE users
(
    id                text        NOT NULL,
    email             text        NOT NULL,
    password          text        NOT NULL,
    enabled           bool        NOT NULL,
    name              text,
    firstname         text,
    lastname          text,
    location          text,
    phone             text,
    profile_photo_url text,
    created_at        timestamptz NOT NULL DEFAULT NOW(),
    updated_at        timestamptz NOT NULL DEFAULT NOW(),
    last_login        timestamptz,
    PRIMARY KEY (id)
);

CREATE INDEX enabled_users_idx ON users (enabled) WHERE enabled;

CREATE TRIGGER created_at_users_trgr
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_trgr
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE products
(
    id             text          NOT NULL,
    user_seller_id text          NOT NULL,
    name           text          NOT NULL,
    description    text          NOT NULL,
    sku            text,
    price          NUMERIC(9, 2) NOT NULL CHECK (price >= 0),
    stock          INT                    DEFAULT 0 CHECK (stock >= 0),
    category       text,
    created_at     timestamptz   NOT NULL DEFAULT NOW(),
    updated_at     timestamptz   NOT NULL DEFAULT NOW(),
    deleted_at     timestamptz,
    PRIMARY KEY (id)
);

CREATE INDEX user_products_idx ON products (user_seller_id);

CREATE TRIGGER created_at_products_trgr
    BEFORE UPDATE
    ON products
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_products_trgr
    BEFORE UPDATE
    ON products
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
    ON snapshots
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

CREATE INDEX users_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;


-- +goose Down
DROP SCHEMA IF EXISTS users CASCADE;
-- SET SEARCH_PATH TO users;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS snapshots;
-- DROP TABLE IF EXISTS events;
-- DROP TABLE IF EXISTS products;
-- DROP TABLE IF EXISTS users;
