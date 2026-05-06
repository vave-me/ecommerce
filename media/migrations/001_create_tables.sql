-- +goose Up
CREATE TABLE media
(
    id          text        NOT NULL,
    item_id     text        NOT NULL,
    item_type   text        NOT NULL,
    user_id     text        NOT NULL,
    status      text        NOT NULL,
    media_order text,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TRIGGER created_at_media_trgr
    BEFORE UPDATE
    ON media
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_media_trgr
    BEFORE UPDATE
    ON media
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TABLE images
(
    id            text        NOT NULL,
    media_id      text        NOT NULL,
    display_order int         NOT NULL,
    is_main       bool        NOT NULL,
    url           text        NOT NULL,
    metadata      text        NOT NULL,
    thumbnail     text,
    user_id       text,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX media_images_idx ON images (media_id);

CREATE TRIGGER created_at_images_trgr
    BEFORE UPDATE
    ON images
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_images_trgr
    BEFORE UPDATE
    ON images
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();


CREATE TABLE videos
(
    id            text        NOT NULL,
    media_id      text        NOT NULL,
    display_order int         NOT NULL,
    is_main       bool        NOT NULL,
    url           text        NOT NULL,
    metadata      text        NOT NULL,
    thumbnail     text,
    user_id       text,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX media_videos_idx ON videos (media_id);

CREATE TRIGGER created_at_videos_trgr
    BEFORE UPDATE
    ON videos
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();
CREATE TRIGGER updated_at_videos_trgr
    BEFORE UPDATE
    ON videos
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();



CREATE TABLE import_sessions
(
    id                   UUID PRIMARY KEY,
    external_system_id   VARCHAR(50) NOT NULL,
    external_system_type VARCHAR(20) NOT NULL,
    total_images         INT         NOT NULL,
    processed_images     INT DEFAULT 0,
    failed_images        INT DEFAULT 0,
    status               VARCHAR(20) NOT NULL,
    started_at           TIMESTAMP   NOT NULL,
    completed_at         TIMESTAMP,
    metadata             JSONB
);

-- Import items tracking
CREATE TABLE import_items
(
    id            UUID PRIMARY KEY,
    session_id    UUID REFERENCES import_sessions (id),
    external_id   VARCHAR(255) NOT NULL,
    sku           VARCHAR(100) NOT NULL,
    image_url     TEXT         NOT NULL,
    status        VARCHAR(20)  NOT NULL,
    error_message TEXT,
    retry_count   INT DEFAULT 0,
    media_id      UUID REFERENCES media (id),
    image_id      UUID REFERENCES images (id)
);

-- SKU mapping cache
CREATE TABLE sku_mappings
(
    sku       VARCHAR(100) PRIMARY KEY,
    item_id   UUID        NOT NULL,
    item_type VARCHAR(20) NOT NULL,
    last_used TIMESTAMP DEFAULT NOW()
);
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
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS videos;
DROP TABLE IF EXISTS media;
DROP TABLE IF EXISTS import_sessions;
DROP TABLE IF EXISTS import_items;
DROP TABLE IF EXISTS sku_mappings;


DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;