-- +goose Up
CREATE TABLE shipping
(
    id                    text        NOT NULL,
    product_id            text        NOT NULL,
    order_id              text,
    basket_id             text,
    tracking_number       text        NOT NULL,
    label_url             text,
    sender_name           text        NOT NULL,
    sender_address        text        NOT NULL,
    receiver_name         text        NOT NULL,
    receiver_address      text        NOT NULL,
    weight                text        NOT NULL,
    dimensions            text        NOT NULL,
    service_type          text        NOT NULL,
    status                text        NOT NULL DEFAULT 'created',
    carrier_id            text,
    carrier_name          text,
    pickup_scheduled_at   timestamptz,
    delivered_at          timestamptz,
    cancelled_at          timestamptz,
    created_at            timestamptz NOT NULL DEFAULT NOW(),
    updated_at            timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_shipping_tracking_number ON shipping(tracking_number);
CREATE INDEX idx_shipping_product_id ON shipping(product_id);
CREATE INDEX idx_shipping_order_id ON shipping(order_id);
CREATE INDEX idx_shipping_status ON shipping(status);
CREATE INDEX idx_shipping_carrier_id ON shipping(carrier_id);

CREATE TRIGGER created_at_shipping_trgr
    BEFORE UPDATE
    ON shipping
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_shipping_trgr
    BEFORE UPDATE
    ON shipping
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE shipment_events
(
    id               text        NOT NULL,
    shipment_id      text        NOT NULL,
    event_type       text        NOT NULL,
    status           text        NOT NULL,
    location         text,
    description      text,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (shipment_id) REFERENCES shipping(id) ON DELETE CASCADE
);

CREATE INDEX idx_shipment_events_shipment_id ON shipment_events(shipment_id);
CREATE INDEX idx_shipment_events_created_at ON shipment_events(created_at);




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
DROP TABLE IF EXISTS shipping;


DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;