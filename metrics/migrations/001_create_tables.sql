-- +goose Up
CREATE
EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users_metrics_cache
(
    id                      text        NOT NULL,
    entity_type             text        NOT NULL,
    likes_count             int                  default 0,
    dislikes_count          int                  default 0,
    comments_count          int                  default 0,
    messages_count          int                  default 0,
    shared_count            int                  default 0,
    added_to_wishlist_count int                  default 0,
    added_to_basket_count   int                  default 0,
    visited_count           int                  default 0,
    reported_count          int                  default 0,
    follower_count          int                  default 0,
    reviews_count           int                  default 0,
    rating_count            int                  default 0,
    videos_count            int                  default 0,
    images_count            int                  default 0,
    rating                  int                  default 0,
    review                  int                  default 0,
    category                text,
    category_id             text,
    category_slug           text,
    media_added_count       int                  default 0,
    comment_added_count     int                  default 0,
    liked_added_count       int                  default 0,
    products_added_count    int                  default 0,
    videos_added_count      int                  default 0,
    services_added_count    int                  default 0,
    jobs_added_count        int                  default 0,
    posts_added_count       int                  default 0,
    vehicles_added_count    int                  default 0,
    properties_added_count  int                  default 0,
    location                geography(Point, 4326),
    lat                     DOUBLE PRECISION     DEFAULT 0.00,
    lng                     DOUBLE PRECISION     DEFAULT 0.00,
    created_at              timestamptz NOT NULL DEFAULT NOW(),
    updated_at              timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_users_metrics_cache_trgr
    BEFORE UPDATE
    ON users_metrics_cache
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_users_metrics_cache_trgr
    BEFORE UPDATE
    ON users_metrics_cache
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE items_metrics_cache
(
    id                      text        NOT NULL,
    entity_type             text        NOT NULL,
    likes_count             int                  default 0,
    dislikes_count          int                  default 0,
    comments_count          int                  default 0,
    messages_count          int                  default 0,
    shared_count            int                  default 0,
    added_to_wishlist_count int                  default 0,
    added_to_basket_count   int                  default 0,
    visited_count           int                  default 0,
    reported_count          int                  default 0,
    follower_count          int                  default 0,
    reviews_count           int                  default 0,
    rating_count            int                  default 0,
    videos_count            int                  default 0,
    images_count            int                  default 0,
    rating                  int                  default 0,
    review                  int                  default 0,
    category                text,
    category_id             text,
    category_slug           text,
    price              INT         NOT NULL CHECK (price >= 0),
    location                geography(Point, 4326),
    lat                     DOUBLE PRECISION     DEFAULT 0.00,
    lng                     DOUBLE PRECISION     DEFAULT 0.00,
    created_at              timestamptz NOT NULL DEFAULT NOW(),
    updated_at              timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_items_metrics_cache_trgr
    BEFORE UPDATE
    ON items_metrics_cache
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_item_metrics_cache_trgr
    BEFORE UPDATE
    ON items_metrics_cache
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

-- +goose Down
DROP TABLE IF EXISTS user_metrics_cache;
DROP TABLE IF EXISTS item_metrics_cache;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
