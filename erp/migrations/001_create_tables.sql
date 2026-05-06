-- +goose Up
-----------------------------------------------------------------------
-- Event-sourcing / snapshots / sagas / outbox / inbox tables
-----------------------------------------------------------------------
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

-- Index for unpublished outbox
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

-----------------------------------------------------------------------
-- ERP-specific tables
-----------------------------------------------------------------------

-- Webhook Events table
CREATE TABLE webhook_events
(
    id VARCHAR(36) PRIMARY KEY,
    erp_type VARCHAR(50) NOT NULL,
    event_id VARCHAR(100),
    event_type VARCHAR(100),
    source VARCHAR(50) NOT NULL,
    signature VARCHAR(255),
    payload JSONB NOT NULL,
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'received',
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    connector_id VARCHAR(100),
    headers JSONB
);

-- Indexes for webhook_events
CREATE INDEX idx_webhook_events_status ON webhook_events(status);
CREATE INDEX idx_webhook_events_erp_type ON webhook_events(erp_type);
CREATE INDEX idx_webhook_events_event_type ON webhook_events(event_type);
CREATE INDEX idx_webhook_events_received_at ON webhook_events(received_at);
CREATE INDEX idx_webhook_events_connector_id ON webhook_events(connector_id);

-- Sync Status table
CREATE TABLE sync_status
(
    id VARCHAR(36) PRIMARY KEY,
    erp_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL DEFAULT '',
    last_synced_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (erp_type, entity_type, entity_id)
);

-- Indexes for sync_status
CREATE INDEX idx_sync_status_erp_entity ON sync_status(erp_type, entity_type);
CREATE INDEX idx_sync_status_status ON sync_status(status);
CREATE INDEX idx_sync_status_last_synced ON sync_status(last_synced_at);

-- Sync Logs table
CREATE TABLE sync_logs
(
    id VARCHAR(36) PRIMARY KEY,
    erp_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    duration INTERVAL,
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    records_processed INTEGER DEFAULT 0,
    records_success INTEGER,
    records_failed INTEGER,
    records_total INTEGER DEFAULT 0,
    last_sync_time TIMESTAMP,
    error_message TEXT,
    error TEXT,
    metadata JSONB,
    connector_id VARCHAR(100)
);

-- Indexes for sync_logs
CREATE INDEX idx_sync_logs_erp_type ON sync_logs(erp_type);
CREATE INDEX idx_sync_logs_entity_type ON sync_logs(entity_type);
CREATE INDEX idx_sync_logs_started_at ON sync_logs(started_at);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_connector_id ON sync_logs(connector_id);

-- Sync Configurations table
CREATE TABLE sync_configurations
(
    id VARCHAR(36) PRIMARY KEY,
    erp_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT true,
    sync_interval INTEGER NOT NULL DEFAULT 300, -- seconds
    batch_size INTEGER NOT NULL DEFAULT 100,
    retry_attempts INTEGER NOT NULL DEFAULT 3,
    retry_delay INTEGER NOT NULL DEFAULT 60, -- seconds
    filters JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (erp_type, entity_type)
);

-- Index for sync_configurations
CREATE INDEX idx_sync_config_enabled ON sync_configurations(enabled);

-- Order Syncs table
CREATE TABLE order_syncs
(
    id VARCHAR(36) PRIMARY KEY,
    connector_id VARCHAR(100) NOT NULL,
    order_id VARCHAR(100) NOT NULL,
    direction VARCHAR(20) NOT NULL, -- 'inbound' or 'outbound'
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempted_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    error TEXT,
    payload JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for order_syncs
CREATE INDEX idx_order_syncs_connector_id ON order_syncs(connector_id);
CREATE INDEX idx_order_syncs_order_id ON order_syncs(order_id);
CREATE INDEX idx_order_syncs_status ON order_syncs(status);
CREATE INDEX idx_order_syncs_attempted_at ON order_syncs(attempted_at);

-- Invoice Syncs table
CREATE TABLE invoice_syncs
(
    id VARCHAR(36) PRIMARY KEY,
    connector_id VARCHAR(100) NOT NULL,
    invoice_id VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'create', 'update', 'approve', 'void', 'send', 'payment'
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempted_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    error TEXT,
    payload JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for invoice_syncs
CREATE INDEX idx_invoice_syncs_connector_id ON invoice_syncs(connector_id);
CREATE INDEX idx_invoice_syncs_invoice_id ON invoice_syncs(invoice_id);
CREATE INDEX idx_invoice_syncs_status ON invoice_syncs(status);
CREATE INDEX idx_invoice_syncs_action ON invoice_syncs(action);
CREATE INDEX idx_invoice_syncs_attempted_at ON invoice_syncs(attempted_at);

-----------------------------------------------------------------------
-- Connectors table for storing ERP connector configurations
-----------------------------------------------------------------------
CREATE TABLE connectors
(
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('odoo', 'dynamics365', 'netsuite', 'sap', 'erpnext', 'frappe')),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive' CHECK (status IN ('active', 'inactive', 'error', 'maintenance')),
    
    -- Authentication configuration (encrypted)
    auth_config_encrypted BYTEA NOT NULL, -- Encrypted JSON containing credentials
    auth_config_salt VARCHAR(64) NOT NULL, -- Salt for encryption
    
    -- Base configuration
    base_url VARCHAR(500) NOT NULL,
    environment VARCHAR(20) NOT NULL DEFAULT 'production' CHECK (environment IN ('production', 'staging', 'development', 'sandbox')),
    
    -- Webhook configuration
    webhook_enabled BOOLEAN NOT NULL DEFAULT false,
    webhook_secret_encrypted BYTEA, -- Encrypted webhook secret
    webhook_url VARCHAR(500), -- Our webhook URL for this connector
    webhook_events TEXT[], -- Array of event types to subscribe to
    
    -- Sync configuration
    sync_enabled BOOLEAN NOT NULL DEFAULT true,
    sync_interval_seconds INTEGER NOT NULL DEFAULT 300, -- 5 minutes default
    batch_size INTEGER NOT NULL DEFAULT 100,
    
    -- Rate limiting configuration
    rate_limit_requests_per_second INTEGER NOT NULL DEFAULT 10,
    rate_limit_burst INTEGER NOT NULL DEFAULT 20,
    
    -- Retry configuration
    retry_max_attempts INTEGER NOT NULL DEFAULT 3,
    retry_initial_delay_ms INTEGER NOT NULL DEFAULT 1000,
    retry_max_delay_ms INTEGER NOT NULL DEFAULT 60000,
    retry_multiplier NUMERIC(3,2) NOT NULL DEFAULT 2.0,
    
    -- Additional settings
    custom_headers JSONB, -- Custom headers to send with requests
    timeout_seconds INTEGER NOT NULL DEFAULT 30,
    
    -- Health check
    last_health_check_at TIMESTAMP,
    last_health_check_status VARCHAR(20),
    last_health_check_error TEXT,
    
    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL,
    updated_by VARCHAR(100) NOT NULL,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1
);

-- Indexes
CREATE INDEX idx_connectors_type ON connectors(type);
CREATE INDEX idx_connectors_status ON connectors(status);
CREATE INDEX idx_connectors_created_at ON connectors(created_at);
CREATE UNIQUE INDEX idx_connectors_name ON connectors(name);

-- Note: Triggers removed to keep migrations simple
-- Application will handle updated_at and version updates

-- Connector sync entity configurations
-- This table stores which entities to sync for each connector
CREATE TABLE connector_sync_entities
(
    id VARCHAR(36) PRIMARY KEY,
    connector_id VARCHAR(100) NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('product', 'stock', 'price', 'customer', 'order', 'invoice', 'return')),
    enabled BOOLEAN NOT NULL DEFAULT true,
    sync_direction VARCHAR(20) NOT NULL DEFAULT 'bidirectional' CHECK (sync_direction IN ('inbound', 'outbound', 'bidirectional')),
    last_sync_at TIMESTAMP,
    filters JSONB, -- Entity-specific filters
    field_mapping JSONB, -- Field mapping configuration
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (connector_id, entity_type)
);

CREATE INDEX idx_connector_sync_entities_connector ON connector_sync_entities(connector_id);
CREATE INDEX idx_connector_sync_entities_enabled ON connector_sync_entities(enabled);

-- Connector API keys table
-- For connectors that support multiple API keys (e.g., for different operations)
CREATE TABLE connector_api_keys
(
    id VARCHAR(36) PRIMARY KEY,
    connector_id VARCHAR(100) NOT NULL REFERENCES connectors(id) ON DELETE CASCADE,
    key_name VARCHAR(100) NOT NULL,
    key_value_encrypted BYTEA NOT NULL, -- Encrypted API key
    key_salt VARCHAR(64) NOT NULL,
    key_type VARCHAR(50) NOT NULL DEFAULT 'general', -- 'general', 'read_only', 'write_only', etc.
    expires_at TIMESTAMP,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (connector_id, key_name)
);

CREATE INDEX idx_connector_api_keys_connector ON connector_api_keys(connector_id);
CREATE INDEX idx_connector_api_keys_expires ON connector_api_keys(expires_at);

-- Audit log for connector changes
CREATE TABLE connector_audit_log
(
    id VARCHAR(36) PRIMARY KEY,
    connector_id VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL CHECK (action IN ('created', 'updated', 'deleted', 'activated', 'deactivated', 'credentials_updated')),
    changed_by VARCHAR(100) NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT
);

CREATE INDEX idx_connector_audit_connector ON connector_audit_log(connector_id);
CREATE INDEX idx_connector_audit_changed_at ON connector_audit_log(changed_at);
CREATE INDEX idx_connector_audit_action ON connector_audit_log(action);

-- +goose Down
-----------------------------------------------------------------------
-- Drop in reverse order
-----------------------------------------------------------------------
DROP TABLE IF EXISTS connector_audit_log CASCADE;
DROP TABLE IF EXISTS connector_api_keys CASCADE;
DROP TABLE IF EXISTS connector_sync_entities CASCADE;
DROP TABLE IF EXISTS connectors CASCADE;
DROP TABLE IF EXISTS invoice_syncs CASCADE;
DROP TABLE IF EXISTS order_syncs CASCADE;
DROP TABLE IF EXISTS sync_configurations CASCADE;
DROP TABLE IF EXISTS sync_logs CASCADE;
DROP TABLE IF EXISTS sync_status CASCADE;
DROP TABLE IF EXISTS webhook_events CASCADE;
DROP TABLE IF EXISTS sagas CASCADE;
DROP TABLE IF EXISTS outbox CASCADE;
DROP TABLE IF EXISTS inbox CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;
DROP TABLE IF EXISTS events CASCADE;