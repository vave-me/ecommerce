-- +goose Up
---------------------------------------------------------------------
-- 1) Create assistants table with all necessary columns
---------------------------------------------------------------------
CREATE TABLE assistants
(
    id             TEXT        NOT NULL,
    assistant_name TEXT,
    description    TEXT,
    user_id        TEXT,
    type           TEXT                 DEFAULT 'standard',
    capabilities   TEXT[],
    enabled        BOOL        NOT NULL DEFAULT false,
    temperature    REAL,
    max_tokens     INTEGER,
    system_prompt  TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

-- Indexes for assistants table
CREATE INDEX enabled_assistants_idx ON assistants (enabled) WHERE enabled;
CREATE INDEX assistants_capabilities_idx ON assistants USING GIN (capabilities);
CREATE INDEX idx_assistants_type ON assistants (type);
CREATE INDEX idx_assistants_user_id ON assistants (user_id);
CREATE INDEX idx_assistants_enabled_user ON assistants (enabled, user_id) WHERE enabled;

---------------------------------------------------------------------
-- 2) Create conversations table
---------------------------------------------------------------------
CREATE TABLE conversations
(
    id                   TEXT        NOT NULL,
    user_id              TEXT        NOT NULL,
    assistant_id         TEXT        NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active               BOOLEAN     NOT NULL DEFAULT true,
    context              JSONB,
    message_count        INTEGER     NOT NULL DEFAULT 0,
    last_message_at      TIMESTAMPTZ,
    last_message_role    TEXT,
    last_message_preview TEXT,
    assistant_type       TEXT                 DEFAULT 'standard',

    PRIMARY KEY (id)
);

-- Index for user conversation lookups
CREATE INDEX conversations_user_id_idx ON conversations (user_id, active) WHERE active;

---------------------------------------------------------------------
-- 3) Create conversation_messages table
---------------------------------------------------------------------
CREATE TABLE conversation_messages
(
    id              TEXT        NOT NULL,
    conversation_id TEXT        NOT NULL,
    assistant_id    TEXT,
    role            TEXT        NOT NULL,
    content         TEXT        NOT NULL,
    timestamp       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata        JSONB,
    actions_taken   JSONB,

    PRIMARY KEY (id),
    FOREIGN KEY (conversation_id) REFERENCES conversations (id) ON DELETE CASCADE
);

-- Performance indexes for message queries
CREATE INDEX conversation_messages_conversation_timestamp_idx ON conversation_messages (conversation_id, timestamp ASC);
CREATE INDEX conversation_messages_user_role_timestamp_idx ON conversation_messages (conversation_id, role, timestamp) WHERE role = 'user';
CREATE INDEX conversation_messages_metadata_idx ON conversation_messages USING GIN (metadata) WHERE metadata IS NOT NULL;
CREATE INDEX conversation_messages_assistant_id_idx ON conversation_messages (assistant_id) WHERE assistant_id IS NOT NULL;

---------------------------------------------------------------------
-- 4) Create event-sourcing tables
---------------------------------------------------------------------
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

-- Partial index for unpublished outbox
CREATE INDEX idx_outbox_unpublished
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


-- +goose Down
---------------------------------------------------------------------
-- Drop all tables and functions in reverse order
---------------------------------------------------------------------
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS conversation_messages;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS assistants;