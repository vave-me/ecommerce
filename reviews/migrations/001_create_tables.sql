-------------------------------------------------------------------------------
-- +goose Up
-------------------------------------------------------------------------------

-- 1) Create table `reviews`, triggers, indexes
CREATE TABLE reviews
(
    id          TEXT        NOT NULL,
    sender_id   TEXT        NOT NULL,
    item_id     TEXT        NOT NULL,
    item_type   TEXT        NOT NULL,
    content     TEXT        NOT NULL,
    category_id TEXT        NOT NULL DEFAULT '',
    parent_id   TEXT        NOT NULL DEFAULT '',
    approved    BOOL        NOT NULL DEFAULT TRUE,
    flagged     BOOL        NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX approved_reviews_idx
    ON reviews (approved) WHERE approved;

CREATE TRIGGER created_at_reviews_trgr
    BEFORE UPDATE
    ON reviews
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_reviews_trgr
    BEFORE UPDATE
    ON reviews
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();


-------------------------------------------------------------------------------
-- 2) Create table `item_review_counts` for aggregator
-------------------------------------------------------------------------------
CREATE TABLE item_review_counts
(
    item_id        TEXT   NOT NULL,
    item_type      TEXT   NOT NULL,
    category_id    TEXT   NOT NULL DEFAULT '',
    total_reviews BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (item_id, item_type, category_id)
);

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION increment_review_count()
RETURNS TRIGGER AS
$$
BEGIN
    -- Simple example: increment if a new review is inserted
UPDATE item_review_counts
SET total_reviews = total_reviews + 1
WHERE item_id = NEW.item_id
  AND item_type = NEW.item_type
  AND category_id = COALESCE(NEW.category_id, '');

IF
NOT FOUND THEN
        INSERT INTO item_review_counts (item_id, item_type, category_id, total_reviews)
        VALUES (NEW.item_id, NEW.item_type, COALESCE(NEW.category_id, ''), 1);
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION decrement_review_count()
RETURNS TRIGGER AS
$$
BEGIN
UPDATE item_review_counts
SET total_reviews = total_reviews - 1
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type
  AND category_id = COALESCE(OLD.category_id, '');

-- Optionally reset negatives to zero:
UPDATE item_review_counts
SET total_reviews = 0
WHERE total_reviews < 0;

RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION update_review_count()
RETURNS TRIGGER AS
$$
BEGIN
    IF
(OLD.item_id <> NEW.item_id)
       OR (OLD.category_id <> NEW.category_id)
       OR (OLD.item_type <> NEW.item_type)
    THEN
        -- Decrement the old item/category/type combo
UPDATE item_review_counts
SET total_reviews = total_reviews - 1
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type
  AND category_id = COALESCE(OLD.category_id, '');

UPDATE item_review_counts
SET total_reviews = 0
WHERE total_reviews < 0;

-- Increment the new item/category/type combo
UPDATE item_review_counts
SET total_reviews = total_reviews + 1
WHERE item_id = NEW.item_id
  AND item_type = NEW.item_type
  AND category_id = COALESCE(NEW.category_id, '');

IF
NOT FOUND THEN
            INSERT INTO item_review_counts (item_id, item_type, category_id, total_reviews)
            VALUES (NEW.item_id, NEW.item_type, COALESCE(NEW.category_id, ''), 1);
END IF;
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS reviews_insert_trigger ON reviews;
CREATE TRIGGER reviews_insert_trigger
    AFTER INSERT
    ON reviews
    FOR EACH ROW
    EXECUTE PROCEDURE increment_review_count();

DROP TRIGGER IF EXISTS reviews_delete_trigger ON reviews;
CREATE TRIGGER reviews_delete_trigger
    AFTER DELETE
    ON reviews
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_review_count();

DROP TRIGGER IF EXISTS reviews_update_trigger ON reviews;
CREATE TRIGGER reviews_update_trigger
    AFTER UPDATE
    ON reviews
    FOR EACH ROW
    EXECUTE PROCEDURE update_review_count();


-------------------------------------------------------------------------------
-- 3) Create tables for events, snapshots, inbox, outbox, sagas
-------------------------------------------------------------------------------
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

CREATE INDEX unpublished_idx
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

-- Always drop triggers before dropping the table that hosts them
DROP TRIGGER IF EXISTS reviews_insert_trigger ON reviews;
DROP TRIGGER IF EXISTS reviews_delete_trigger ON reviews;
DROP TRIGGER IF EXISTS reviews_update_trigger ON reviews;

DROP FUNCTION IF EXISTS increment_review_count();
DROP FUNCTION IF EXISTS decrement_review_count();
DROP FUNCTION IF EXISTS update_review_count();

DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS item_review_counts;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
