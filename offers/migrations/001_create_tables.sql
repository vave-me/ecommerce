-- +goose Up

--------------------------------------------------------------------------
-- Table: offers
--------------------------------------------------------------------------
CREATE TABLE offers
(
    id               text        NOT NULL,
    user_seller_id   text        NOT NULL,
    user_customer_id text        NOT NULL,
    product_id       text        NOT NULL,
    price            float       NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index specific to "offers" (unique name)
CREATE INDEX offers_active_idx ON offers (active) WHERE active;

-- Triggers named specifically for the offers table
CREATE TRIGGER created_at_offers_trgr
    BEFORE UPDATE
    ON offers
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_offers_trgr
    BEFORE UPDATE
    ON offers
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

--------------------------------------------------------------------------
-- Table: lease
--------------------------------------------------------------------------
CREATE TABLE lease
(
    id               text        NOT NULL,
    user_seller_id   text        NOT NULL,
    user_customer_id text        NOT NULL,
    product_id       text        NOT NULL,
    price            float       NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index unique to "lease" table
CREATE INDEX lease_active_idx ON lease (active) WHERE active;

-- Triggers for "lease" table
CREATE TRIGGER created_at_lease_trgr
    BEFORE UPDATE
    ON lease
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_lease_trgr
    BEFORE UPDATE
    ON lease
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

--------------------------------------------------------------------------
-- Table: buy_back
--------------------------------------------------------------------------
CREATE TABLE buy_back
(
    id               text        NOT NULL,
    user_seller_id   text        NOT NULL,
    user_customer_id text        NOT NULL,
    product_id       text        NOT NULL,
    price            float       NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index unique to "buy_back" table
CREATE INDEX buy_back_active_idx ON buy_back (active) WHERE active;

-- Triggers for "buy_back" table
CREATE TRIGGER created_at_buy_back_trgr
    BEFORE UPDATE
    ON buy_back
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_buy_back_trgr
    BEFORE UPDATE
    ON buy_back
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

--------------------------------------------------------------------------
-- Table: buy_now
--------------------------------------------------------------------------
CREATE TABLE buy_now
(
    id               text        NOT NULL,
    user_seller_id   text        NOT NULL,
    user_customer_id text        NOT NULL,
    product_id       text        NOT NULL,
    price            float       NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index unique to "buy_now" table
CREATE INDEX buy_now_active_idx ON buy_now (active) WHERE active;

-- Triggers for "buy_now" table
CREATE TRIGGER created_at_buy_now_trgr
    BEFORE UPDATE
    ON buy_now
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_buy_now_trgr
    BEFORE UPDATE
    ON buy_now
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();


--------------------------------------------------------------------------
-- Table: buy_now
--------------------------------------------------------------------------
CREATE TABLE reservations
(
    id               text        NOT NULL,
    user_seller_id   text        NOT NULL,
    user_customer_id text        NOT NULL,
    product_id       text        NOT NULL,
    price            float       NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index unique to "buy_now" table
CREATE INDEX reservations_active_idx ON reservations (active) WHERE active;

-- Triggers for "buy_now" table
CREATE TRIGGER created_at_reservations_trgr
    BEFORE UPDATE
    ON reservations
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_reservations_trgr
    BEFORE UPDATE
    ON reservations
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();
--------------------------------------------------------------------------
-- Table: events
--------------------------------------------------------------------------
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

--------------------------------------------------------------------------
-- Table: snapshots
--------------------------------------------------------------------------
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

--------------------------------------------------------------------------
-- Table: inbox
--------------------------------------------------------------------------
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

--------------------------------------------------------------------------
-- Table: outbox
--------------------------------------------------------------------------
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

CREATE INDEX unpublished_outbox_idx ON outbox (published_at) WHERE published_at IS NULL;

--------------------------------------------------------------------------
-- Table: sagas
--------------------------------------------------------------------------
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
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();


-- +goose Down
DROP TABLE IF EXISTS offers CASCADE;
DROP TABLE IF EXISTS lease CASCADE;
DROP TABLE IF EXISTS buy_back CASCADE;
DROP TABLE IF EXISTS buy_now CASCADE;

DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;


