-- +goose Up

-- Scheduler table: main aggregate for user schedulers
CREATE TABLE scheduler
(
    id         TEXT        NOT NULL,
    user_id    TEXT        NOT NULL,
    enabled    BOOL        NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Actions table: scheduled natural language tasks
CREATE TABLE actions
(
    id                    TEXT        NOT NULL,
    scheduler_id          TEXT        NOT NULL,
    natural_language_task TEXT        NOT NULL,
    execution_time        TIMESTAMPTZ NOT NULL,
    status                TEXT        NOT NULL DEFAULT 'pending', -- pending, executing, completed, failed
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    executed_at           TIMESTAMPTZ,
    result                TEXT,
    error_message         TEXT,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

-- Index on scheduler_id for quick lookups
CREATE INDEX scheduler_actions_idx
    ON actions (scheduler_id);

-- Index on execution_time for the scheduler worker to find pending tasks
CREATE INDEX pending_actions_execution_time_idx
    ON actions (execution_time) WHERE status = 'pending';

-- Index on status for monitoring
CREATE INDEX actions_status_idx
    ON actions (status);

-- Timestamps triggers for "scheduler"
CREATE TRIGGER created_at_scheduler_trgr
    BEFORE UPDATE
    ON scheduler
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_scheduler_trgr
    BEFORE UPDATE
    ON scheduler
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- Timestamps triggers for "actions"
CREATE TRIGGER created_at_actions_trgr
    BEFORE UPDATE
    ON actions
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_actions_trgr
    BEFORE UPDATE
    ON actions
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- Event sourcing tables
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

-- Message handling tables
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

-- Saga orchestration table
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

-- Catalog tasks table: read model for manager tasks
CREATE TABLE catalog_tasks
(
    id             TEXT        NOT NULL,
    manager_id     TEXT        NOT NULL,
    task_type      TEXT        NOT NULL,
    scheduled_at   TIMESTAMPTZ NOT NULL,
    payload        JSONB       NOT NULL DEFAULT '{}',
    status         TEXT        NOT NULL DEFAULT 'pending', -- pending, executing, completed, failed, cancelled
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    executed_at    TIMESTAMPTZ,
    result         TEXT,
    error_message  TEXT,
    PRIMARY KEY (id)
);

-- Index on manager_id for quick lookups by manager
CREATE INDEX catalog_tasks_manager_id_idx
    ON catalog_tasks (manager_id);

-- Index on scheduled_at for finding tasks due for execution
CREATE INDEX catalog_tasks_scheduled_at_idx
    ON catalog_tasks (scheduled_at) WHERE status = 'pending';

-- Index on status for monitoring and filtering
CREATE INDEX catalog_tasks_status_idx
    ON catalog_tasks (status);

-- Composite index for manager queries with status filter
CREATE INDEX catalog_tasks_manager_status_idx
    ON catalog_tasks (manager_id, status);

-- Index on task_type for filtering by type
CREATE INDEX catalog_tasks_type_idx
    ON catalog_tasks (task_type);

-- Timestamps triggers for "catalog_tasks"
CREATE TRIGGER created_at_catalog_tasks_trgr
    BEFORE UPDATE
    ON catalog_tasks
    FOR EACH ROW
    EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_catalog_tasks_trgr
    BEFORE UPDATE
    ON catalog_tasks
    FOR EACH ROW
    EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down

DROP TABLE IF EXISTS catalog_tasks;
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS actions;
DROP TABLE IF EXISTS scheduler;