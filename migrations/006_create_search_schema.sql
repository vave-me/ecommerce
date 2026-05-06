-- +goose Up
CREATE SCHEMA search;

SET
SEARCH_PATH TO search, public;

CREATE TABLE users_cache
(
    id         text        NOT NULL,
    first_name text        NOT NULL,
    last_name  text        NOT NULL,
    email      text        NOT NULL,
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
    id             text        NOT NULL,
    user_seller_id text        NOT NULL,
    name           text        NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT NOW(),
    updated_at     timestamptz NOT NULL DEFAULT NOW(),
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

CREATE TABLE orders
(
    order_id           text        NOT NULL,
    user_customer_id   text        NOT NULL,
    user_customer_name text        NOT NULL,
    items              bytea       NOT NULL,
    status             text        NOT NULL,
    product_ids        text ARRAY NOT NULL,
    user_seller_ids    text ARRAY NOT NULL,
    created_at         timestamptz NOT NULL DEFAULT NOW(),
    updated_at         timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (order_id)
);

CREATE TRIGGER updated_at_sorders_trgr
    BEFORE UPDATE
    ON orders
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

CREATE INDEX search_unpublished_idx ON outbox (published_at) WHERE published_at IS NULL;


-- +goose Down
DROP SCHEMA IF EXISTS search CASCADE;
-- SET SEARCH_PATH TO search;
--
-- DROP TABLE IF EXISTS outbox;
-- DROP TABLE IF EXISTS inbox;
-- DROP TABLE IF EXISTS orders;
-- DROP TABLE IF EXISTS products_cache;
-- DROP TABLE IF EXISTS users_cache;
