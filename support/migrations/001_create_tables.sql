-- +goose Up
-- Support Channels
CREATE TABLE support_channels
(
    id               text        NOT NULL,
    user_id          text        NOT NULL,
    business_id      text        NOT NULL,
    channel_type     text        NOT NULL,
    active           bool        NOT NULL DEFAULT true,
    settings         jsonb       NOT NULL DEFAULT '{}',
    open_tickets     int         NOT NULL DEFAULT 0,
    total_tickets    int         NOT NULL DEFAULT 0,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    closed_at        timestamptz,
    PRIMARY KEY (id)
);

CREATE INDEX idx_support_channels_user_id ON support_channels (user_id);
CREATE INDEX idx_support_channels_business_id ON support_channels (business_id);
CREATE INDEX idx_support_channels_active ON support_channels (active) WHERE active;

-- Tickets
CREATE TABLE tickets
(
    id                   text        NOT NULL,
    channel_id           text        NOT NULL REFERENCES support_channels(id),
    title                text        NOT NULL,
    description          text        NOT NULL,
    status               text        NOT NULL DEFAULT 'SUBMITTED',
    priority             text        NOT NULL DEFAULT 'MEDIUM',
    category             text        NOT NULL,
    tags                 text[]      NOT NULL DEFAULT '{}',
    metadata             jsonb       NOT NULL DEFAULT '{}',
    assignee_id          text,
    assignee_type        text,
    created_by           text        NOT NULL,
    current_tier         text        NOT NULL DEFAULT 'TIER_1',
    response_count       int         NOT NULL DEFAULT 0,
    reopen_count         int         NOT NULL DEFAULT 0,
    satisfaction_rating  text,
    linked_ticket_ids    text[]      NOT NULL DEFAULT '{}',
    merged_ticket_ids    text[]      NOT NULL DEFAULT '{}',
    created_at           timestamptz NOT NULL DEFAULT NOW(),
    updated_at           timestamptz NOT NULL DEFAULT NOW(),
    resolved_at          timestamptz,
    closed_at            timestamptz,
    first_response_at    timestamptz,
    PRIMARY KEY (id)
);

CREATE INDEX idx_tickets_channel_id ON tickets (channel_id);
CREATE INDEX idx_tickets_status ON tickets (status);
CREATE INDEX idx_tickets_priority ON tickets (priority);
CREATE INDEX idx_tickets_assignee ON tickets (assignee_id) WHERE assignee_id IS NOT NULL;
CREATE INDEX idx_tickets_created_at ON tickets (created_at);

-- Communications (replies and notes)
CREATE TABLE communications
(
    id               text        NOT NULL,
    ticket_id        text        NOT NULL REFERENCES tickets(id),
    author_id        text        NOT NULL,
    author_type      text        NOT NULL,
    content          text        NOT NULL,
    is_public        bool        NOT NULL DEFAULT true,
    mentioned_users  text[]      NOT NULL DEFAULT '{}',
    metadata         jsonb       NOT NULL DEFAULT '{}',
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_communications_ticket_id ON communications (ticket_id);
CREATE INDEX idx_communications_author_id ON communications (author_id);
CREATE INDEX idx_communications_created_at ON communications (created_at);

-- Attachments
CREATE TABLE attachments
(
    id              text        NOT NULL,
    entity_id       text        NOT NULL, -- Can be ticket_id or communication_id
    entity_type     text        NOT NULL, -- 'ticket' or 'communication'
    filename        text        NOT NULL,
    content_type    text        NOT NULL,
    size_bytes      bigint      NOT NULL,
    url             text        NOT NULL,
    uploaded_at     timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_attachments_entity ON attachments (entity_id, entity_type);

-- Knowledge Articles
CREATE TABLE knowledge_articles
(
    id                   text        NOT NULL,
    title                text        NOT NULL,
    content              text        NOT NULL,
    categories           text[]      NOT NULL DEFAULT '{}',
    tags                 text[]      NOT NULL DEFAULT '{}',
    public               bool        NOT NULL DEFAULT true,
    view_count           int         NOT NULL DEFAULT 0,
    average_rating       float       NOT NULL DEFAULT 0,
    rating_count         int         NOT NULL DEFAULT 0,
    created_by           text        NOT NULL,
    related_ticket_ids   text[]      NOT NULL DEFAULT '{}',
    created_at           timestamptz NOT NULL DEFAULT NOW(),
    updated_at           timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_knowledge_articles_categories ON knowledge_articles USING GIN (categories);
CREATE INDEX idx_knowledge_articles_tags ON knowledge_articles USING GIN (tags);
CREATE INDEX idx_knowledge_articles_public ON knowledge_articles (public) WHERE public;

-- Article Ratings
CREATE TABLE article_ratings
(
    id              text        NOT NULL,
    article_id      text        NOT NULL REFERENCES knowledge_articles(id),
    rated_by        text        NOT NULL,
    rating          int         NOT NULL CHECK (rating >= 1 AND rating <= 5),
    feedback        text,
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (article_id, rated_by)
);

CREATE INDEX idx_article_ratings_article_id ON article_ratings (article_id);

-- AI Configurations
CREATE TABLE ai_configurations
(
    id              text        NOT NULL,
    channel_id      text        NOT NULL REFERENCES support_channels(id),
    assistant_id    text        NOT NULL,
    configuration   jsonb       NOT NULL,
    active          bool        NOT NULL DEFAULT true,
    created_at      timestamptz NOT NULL DEFAULT NOW(),
    updated_at      timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (channel_id)
);

CREATE INDEX idx_ai_configurations_channel_id ON ai_configurations (channel_id);

-- Catalog view for support channels
CREATE TABLE support_channels_catalog
(
    id               text        NOT NULL,
    user_id          text        NOT NULL,
    business_id      text        NOT NULL,
    channel_type     text        NOT NULL,
    active           bool        NOT NULL,
    open_tickets     int         NOT NULL,
    total_tickets    int         NOT NULL,
    created_at       timestamptz NOT NULL,
    updated_at       timestamptz NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_support_channels_catalog_user_id ON support_channels_catalog (user_id);
CREATE INDEX idx_support_channels_catalog_business_id ON support_channels_catalog (business_id);

-- Catalog view for tickets
CREATE TABLE tickets_catalog
(
    id               text        NOT NULL,
    channel_id       text        NOT NULL,
    title            text        NOT NULL,
    status           text        NOT NULL,
    priority         text        NOT NULL,
    category         text        NOT NULL,
    assignee_id      text,
    assignee_type    text,
    created_by       text        NOT NULL,
    created_at       timestamptz NOT NULL,
    updated_at       timestamptz NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_tickets_catalog_channel_id ON tickets_catalog (channel_id);
CREATE INDEX idx_tickets_catalog_status ON tickets_catalog (status);
CREATE INDEX idx_tickets_catalog_assignee ON tickets_catalog (assignee_id) WHERE assignee_id IS NOT NULL;

-- Catalog view for knowledge articles
CREATE TABLE knowledge_articles_catalog
(
    id               text        NOT NULL,
    title            text        NOT NULL,
    categories       text[]      NOT NULL,
    public           bool        NOT NULL,
    view_count       int         NOT NULL,
    average_rating   float       NOT NULL,
    created_at       timestamptz NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX idx_knowledge_articles_catalog_categories ON knowledge_articles_catalog USING GIN (categories);
CREATE INDEX idx_knowledge_articles_catalog_public ON knowledge_articles_catalog (public) WHERE public;

-- Triggers for updated_at
CREATE TRIGGER created_at_support_channels_trgr
    BEFORE UPDATE ON support_channels
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_support_channels_trgr
    BEFORE UPDATE ON support_channels
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TRIGGER created_at_tickets_trgr
    BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_tickets_trgr
    BEFORE UPDATE ON tickets
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

CREATE TRIGGER created_at_knowledge_articles_trgr
    BEFORE UPDATE ON knowledge_articles
    FOR EACH ROW EXECUTE PROCEDURE created_at_trigger();

CREATE TRIGGER updated_at_knowledge_articles_trgr
    BEFORE UPDATE ON knowledge_articles
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- Event Sourcing Tables
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
    BEFORE UPDATE ON snapshots
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- Messaging Tables
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
    BEFORE UPDATE ON sagas
    FOR EACH ROW EXECUTE PROCEDURE updated_at_trigger();

-- +goose Down
DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;

DROP TABLE IF EXISTS knowledge_articles_catalog;
DROP TABLE IF EXISTS tickets_catalog;
DROP TABLE IF EXISTS support_channels_catalog;

DROP TABLE IF EXISTS ai_configurations;
DROP TABLE IF EXISTS article_ratings;
DROP TABLE IF EXISTS knowledge_articles;
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS communications;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS support_channels;