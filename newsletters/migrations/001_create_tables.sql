-- +goose Up
-- Newsletter definitions created by users
CREATE TABLE newsletters
(
    id               text        NOT NULL,
    user_id          text        NOT NULL, -- User who created the newsletter
    name             text        NOT NULL,
    description      text,
    frequency        text        NOT NULL, -- daily, weekly, monthly
    category         text,
    template_id      text,
    is_active        bool        NOT NULL DEFAULT true,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX newsletters_user_idx ON newsletters (user_id);
CREATE INDEX active_newsletters_idx ON newsletters (is_active) WHERE is_active;

-- User subscriptions to newsletters
CREATE TABLE newsletter_subscriptions
(
    id               text        NOT NULL,
    user_id          text        NOT NULL, -- Subscriber
    newsletter_id    text        NOT NULL,
    status           text        NOT NULL DEFAULT 'active', -- active, paused, unsubscribed
    preferences      jsonb,      -- frequency override, topics, format preferences
    subscribed_at    timestamptz NOT NULL DEFAULT NOW(),
    unsubscribed_at  timestamptz,
    PRIMARY KEY (id),
    UNIQUE(user_id, newsletter_id)
);

CREATE INDEX subscriptions_user_idx ON newsletter_subscriptions (user_id);
CREATE INDEX subscriptions_newsletter_idx ON newsletter_subscriptions (newsletter_id);
CREATE INDEX active_subscriptions_idx ON newsletter_subscriptions (status) WHERE status = 'active';

-- Newsletter content/editions
CREATE TABLE newsletter_editions
(
    id               text        NOT NULL,
    newsletter_id    text        NOT NULL,
    subject          text        NOT NULL,
    content_html     text        NOT NULL,
    content_text     text,
    template_data    jsonb,      -- Variables for template rendering
    scheduled_at     timestamptz,
    sent_at          timestamptz,
    status           text        NOT NULL DEFAULT 'draft', -- draft, scheduled, sending, sent
    created_by       text        NOT NULL, -- User who created the edition
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX editions_newsletter_idx ON newsletter_editions (newsletter_id);
CREATE INDEX editions_status_idx ON newsletter_editions (status);
CREATE INDEX editions_scheduled_idx ON newsletter_editions (scheduled_at) WHERE status = 'scheduled';

-- Newsletter templates
CREATE TABLE newsletter_templates
(
    id               text        NOT NULL,
    user_id          text,       -- NULL for system templates, user_id for custom templates
    name             text        NOT NULL,
    description      text,
    html_template    text        NOT NULL,
    text_template    text,
    variables        jsonb,      -- Expected variables schema
    preview_data     jsonb,      -- Sample data for preview
    is_public        bool        NOT NULL DEFAULT false,
    created_at       timestamptz NOT NULL DEFAULT NOW(),
    updated_at       timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX templates_user_idx ON newsletter_templates (user_id);
CREATE INDEX public_templates_idx ON newsletter_templates (is_public) WHERE is_public;

-- Sending logs for tracking
CREATE TABLE newsletter_send_logs
(
    id               text        NOT NULL,
    edition_id       text        NOT NULL,
    user_id          text        NOT NULL, -- Recipient
    email            text        NOT NULL,
    status           text        NOT NULL, -- queued, sent, failed, bounced, opened, clicked
    sent_at          timestamptz,
    opened_at        timestamptz,
    clicked_at       timestamptz,
    error_message    text,
    metadata         jsonb,
    PRIMARY KEY (id)
);

CREATE INDEX send_logs_edition_idx ON newsletter_send_logs (edition_id);
CREATE INDEX send_logs_user_idx ON newsletter_send_logs (user_id);
CREATE INDEX send_logs_status_idx ON newsletter_send_logs (status);

-- Event sourcing tables
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

-- +goose Down
DROP TABLE IF EXISTS newsletter_send_logs;
DROP TABLE IF EXISTS newsletter_editions;
DROP TABLE IF EXISTS newsletter_subscriptions;
DROP TABLE IF EXISTS newsletter_templates;
DROP TABLE IF EXISTS newsletters;

DROP TABLE IF EXISTS sagas;
DROP TABLE IF EXISTS outbox;
DROP TABLE IF EXISTS inbox;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS events;