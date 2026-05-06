-- +goose Up

CREATE TABLE activity
(
    id         TEXT        NOT NULL,
    user_id    TEXT        NOT NULL,
    enabled    BOOL        NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Partial index for quick lookups of enabled activities
CREATE INDEX enabled_activity_idx
    ON activity (enabled) WHERE enabled;

-- Timestamps triggers for "activity"
CREATE TRIGGER created_at_activity_trgr
    BEFORE UPDATE
    ON activity
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_activity_trgr
    BEFORE UPDATE
    ON activity
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();


CREATE TABLE interactions
(
    id          TEXT        NOT NULL,
    activity_id TEXT        NOT NULL,
    item_id     TEXT        NOT NULL,
    item_type   TEXT        NOT NULL, -- e.g. "product", "comment", "profile", etc.
    action_type TEXT        NOT NULL, -- e.g. "like", "dislike"
    enabled     BOOL        NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index on activity_id for quick lookups
CREATE INDEX activity_interactions_idx
    ON interactions (activity_id);

-- Optional partial index if you frequently query only enabled interactions
CREATE INDEX enabled_interactions_idx
    ON interactions (activity_id) WHERE enabled;

-- Timestamps triggers for "interactions"
CREATE TRIGGER created_at_interactions_trgr
    BEFORE UPDATE
    ON interactions
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_interactions_trgr
    BEFORE UPDATE
    ON interactions
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();


CREATE TABLE item_interaction_counts
(
    item_id       TEXT   NOT NULL,
    item_type     TEXT   NOT NULL DEFAULT '', -- e.g. "product", "comment", etc.
    total_count   BIGINT NOT NULL DEFAULT 0,
    like_count    BIGINT NOT NULL DEFAULT 0,
    dislike_count BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (item_id, item_type)
);


-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION increment_interaction_count()
RETURNS TRIGGER AS
$$
BEGIN
    /*
      For a new row in "interactions":

      1) Increase total_count by 1 always
      2) If action_type = 'like', increment like_count
      3) If action_type = 'dislike', increment dislike_count
      4) If not found, insert a new row in item_interaction_counts
    */

    IF
NEW.action_type = 'like' THEN
UPDATE item_interaction_counts
SET like_count  = like_count + 1,
    total_count = total_count + 1
WHERE item_id = NEW.item_id
  AND item_type = NEW.item_type;

IF
NOT FOUND THEN
            INSERT INTO item_interaction_counts (item_id, item_type, like_count, dislike_count, total_count)
            VALUES (NEW.item_id, NEW.item_type, 1, 0, 1);
END IF;

    ELSIF
NEW.action_type = 'dislike' THEN
UPDATE item_interaction_counts
SET dislike_count = dislike_count + 1,
    total_count   = total_count + 1
WHERE item_id = NEW.item_id
  AND item_type = NEW.item_type;

IF
NOT FOUND THEN
            INSERT INTO item_interaction_counts (item_id, item_type, like_count, dislike_count, total_count)
            VALUES (NEW.item_id, NEW.item_type, 0, 1, 1);
END IF;

ELSE
        -- If you have other action_types, handle them or just increment total_count
UPDATE item_interaction_counts
SET total_count = total_count + 1
WHERE item_id = NEW.item_id
  AND item_type = NEW.item_type;

IF
NOT FOUND THEN
            INSERT INTO item_interaction_counts (item_id, item_type, total_count)
            VALUES (NEW.item_id, NEW.item_type, 1);
END IF;
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION decrement_interaction_count()
RETURNS TRIGGER AS
$$
BEGIN
    /*
      For a deleted row in "interactions":

      1) Decrease total_count by 1
      2) If action_type = 'like', decrement like_count
      3) If action_type = 'dislike', decrement dislike_count
      4) Then clamp any counts below zero back up to 0
    */

    IF
OLD.action_type = 'like' THEN
UPDATE item_interaction_counts
SET like_count  = like_count - 1,
    total_count = total_count - 1
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type;

ELSIF
OLD.action_type = 'dislike' THEN
UPDATE item_interaction_counts
SET dislike_count = dislike_count - 1,
    total_count   = total_count - 1
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type;

ELSE
UPDATE item_interaction_counts
SET total_count = total_count - 1
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type;
END IF;

    -- Prevent negative counts
UPDATE item_interaction_counts
SET like_count    = GREATEST(like_count, 0),
    dislike_count = GREATEST(dislike_count, 0),
    total_count   = GREATEST(total_count, 0)
WHERE item_id = OLD.item_id
  AND item_type = OLD.item_type;

RETURN OLD;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION update_interaction_count()
RETURNS TRIGGER AS
$$
BEGIN
    /*
      For an updated row in "interactions":

      If (item_id, item_type, or action_type) changed,
      then call decrement_interaction_count() on the OLD row
      and increment_interaction_count() on the NEW row.
      This effectively moves the aggregator from old to new.
    */

    IF
(OLD.item_id <> NEW.item_id)
       OR (OLD.item_type <> NEW.item_type)
       OR (OLD.action_type <> NEW.action_type)
    THEN
        -- 1) Decrement old
        PERFORM decrement_interaction_count();

        -- 2) Increment new
        PERFORM
increment_interaction_count();
END IF;

RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd


------------------------------------------------------------------------------
-- 2c) Bind these aggregator triggers to the "interactions" table
------------------------------------------------------------------------------
DROP TRIGGER IF EXISTS interactions_insert_trigger ON interactions;
CREATE TRIGGER interactions_insert_trigger
    AFTER INSERT
    ON interactions
    FOR EACH ROW
    EXECUTE PROCEDURE increment_interaction_count();

DROP TRIGGER IF EXISTS interactions_delete_trigger ON interactions;
CREATE TRIGGER interactions_delete_trigger
    AFTER DELETE
    ON interactions
    FOR EACH ROW
    EXECUTE PROCEDURE decrement_interaction_count();

DROP TRIGGER IF EXISTS interactions_update_trigger ON interactions;
CREATE TRIGGER interactions_update_trigger
    AFTER UPDATE
    ON interactions
    FOR EACH ROW
    EXECUTE PROCEDURE update_interaction_count();


------------------------------------------------------------------------------
-- 3) Create the other tables: users_cache, products_cache, events, snapshots, etc.
------------------------------------------------------------------------------
CREATE TABLE users_cache
(
    id         TEXT        NOT NULL,
    email      TEXT,
    username   TEXT        NOT NULL,
    location   TEXT,
    enabled    BOOL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
    id             TEXT        NOT NULL,
    name           TEXT        NOT NULL,
    description    TEXT        NOT NULL,
    price          INT         NOT NULL,
    user_seller_id TEXT        NOT NULL,
    stock          INT         NOT NULL,
    sku            TEXT        NOT NULL,
    category_id    TEXT        NOT NULL,
    active         TEXT        NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
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
DROP TABLE IF EXISTS interactions;
DROP TABLE IF EXISTS activity;
DROP TABLE IF EXISTS users_cache;
DROP TABLE IF EXISTS products_cache;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS item_interaction_counts;
