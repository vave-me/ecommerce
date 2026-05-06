-------------------------------------------------------------------------------
-- +goose Up
-------------------------------------------------------------------------------

-- 1) Create table `following`, triggers, indexes
CREATE TABLE following
(
    id          TEXT        NOT NULL,
    user_id   TEXT        NOT NULL,
    followed_user_id     TEXT        NOT NULL,
    followed_user_type   TEXT        NOT NULL,
    content     TEXT        NOT NULL,
    category_id TEXT        NOT NULL DEFAULT '',
    parent_id   TEXT        NOT NULL DEFAULT '',
    approved    BOOL        NOT NULL DEFAULT TRUE,
    flagged     BOOL        NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX approved_following_idx
    ON following (approved) WHERE approved;

CREATE TRIGGER created_at_following_trgr
    BEFORE UPDATE
    ON following
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_following_trgr
    BEFORE UPDATE
    ON following
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();


-------------------------------------------------------------------------------
-- 2) Create table `item_follow_counts` for aggregator
-------------------------------------------------------------------------------
CREATE TABLE item_follow_counts
(
    followed_user_id        TEXT   NOT NULL,
    followed_user_type      TEXT   NOT NULL,
    category_id    TEXT   NOT NULL DEFAULT '',
    total_following BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (followed_user_id, followed_user_type, category_id)
);

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION increment_follow_count()
RETURNS TRIGGER AS
$$
BEGIN
    -- Simple example: increment if a new follow is inserted
UPDATE item_follow_counts
SET total_following = total_following + 1
WHERE followed_user_id = NEW.followed_user_id
  AND followed_user_type = NEW.followed_user_type
  AND category_id = COALESCE(NEW.category_id, '');

IF
NOT FOUND THEN
        INSERT INTO item_follow_counts (followed_user_id, followed_user_type, category_id, total_following)
        VALUES (NEW.followed_user_id, NEW.followed_user_type, COALESCE(NEW.category_id, ''), 1);
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION decrement_follow_count()
RETURNS TRIGGER AS
$$
BEGIN
UPDATE item_follow_counts
SET total_following = total_following - 1
WHERE followed_user_id = OLD.followed_user_id
  AND followed_user_type = OLD.followed_user_type
  AND category_id = COALESCE(OLD.category_id, '');

-- Optionally reset negatives to zero:
UPDATE item_follow_counts
SET total_following = 0
WHERE total_following < 0;

RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION update_follow_count()
RETURNS TRIGGER AS
$$
BEGIN
    IF
(OLD.followed_user_id <> NEW.followed_user_id)
       OR (OLD.category_id <> NEW.category_id)
       OR (OLD.followed_user_type <> NEW.followed_user_type)
    THEN
        -- Decrement the old item/category/type combo
UPDATE item_follow_counts
SET total_following = total_following - 1
WHERE followed_user_id = OLD.followed_user_id
  AND followed_user_type = OLD.followed_user_type
  AND category_id = COALESCE(OLD.category_id, '');

UPDATE item_follow_counts
SET total_following = 0
WHERE total_following < 0;

-- Increment the new item/category/type combo
UPDATE item_follow_counts
SET total_following = total_following + 1
WHERE followed_user_id = NEW.followed_user_id
  AND followed_user_type = NEW.followed_user_type
  AND category_id = COALESCE(NEW.category_id, '');

IF
NOT FOUND THEN
            INSERT INTO item_follow_counts (followed_user_id, followed_user_type, category_id, total_following)
            VALUES (NEW.followed_user_id, NEW.followed_user_type, COALESCE(NEW.category_id, ''), 1);
END IF;
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS following_insert_trigger ON following;
CREATE TRIGGER following_insert_trigger
    AFTER INSERT
    ON following
    FOR EACH ROW
    EXECUTE PROCEDURE increment_follow_count();

DROP TRIGGER IF EXISTS following_delete_trigger ON following;
CREATE TRIGGER following_delete_trigger
    AFTER DELETE
    ON following
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_follow_count();

DROP TRIGGER IF EXISTS following_update_trigger ON following;
CREATE TRIGGER following_update_trigger
    AFTER UPDATE
    ON following
    FOR EACH ROW
    EXECUTE PROCEDURE update_follow_count();


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
DROP TRIGGER IF EXISTS following_insert_trigger ON following;
DROP TRIGGER IF EXISTS following_delete_trigger ON following;
DROP TRIGGER IF EXISTS following_update_trigger ON following;

DROP FUNCTION IF EXISTS increment_follow_count();
DROP FUNCTION IF EXISTS decrement_follow_count();
DROP FUNCTION IF EXISTS update_follow_count();

DROP TABLE IF EXISTS following;
DROP TABLE IF EXISTS item_follow_counts;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
