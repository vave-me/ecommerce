-- +migrate Up

-- Webhook subscriptions table
CREATE TABLE IF NOT EXISTS webhook_subscriptions (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    events JSONB NOT NULL DEFAULT '[]',
    headers JSONB NOT NULL DEFAULT '{}',
    
    -- Retry policy
    retry_max_retries INT NOT NULL DEFAULT 3,
    retry_backoff_factor FLOAT NOT NULL DEFAULT 2.0,
    retry_initial_delay BIGINT NOT NULL DEFAULT 1000, -- milliseconds
    retry_max_backoff BIGINT NOT NULL DEFAULT 300000, -- 5 minutes in milliseconds
    
    -- Status
    active BOOLEAN NOT NULL DEFAULT true,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Webhook deliveries table
CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id UUID PRIMARY KEY,
    subscription_id UUID NOT NULL REFERENCES webhook_subscriptions(id) ON DELETE CASCADE,
    event_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload BYTEA NOT NULL,
    
    -- Delivery status
    status TEXT NOT NULL DEFAULT 'pending',
    attempts INT NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    last_attempt_at TIMESTAMP,
    next_retry_at TIMESTAMP,
    
    -- Response data
    response_status INT,
    response_body TEXT,
    error TEXT,
    
    CONSTRAINT webhook_deliveries_status_check CHECK (status IN ('pending', 'retrying', 'delivered', 'failed'))
);

-- Indexes for performance
CREATE INDEX idx_webhook_subscriptions_active ON webhook_subscriptions(active);
CREATE INDEX idx_webhook_subscriptions_events ON webhook_subscriptions USING GIN(events);
CREATE INDEX idx_webhook_deliveries_subscription_id ON webhook_deliveries(subscription_id);
CREATE INDEX idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at) WHERE status IN ('pending', 'retrying');
CREATE INDEX idx_webhook_deliveries_created_at ON webhook_deliveries(created_at);

-- +migrate Down

DROP TABLE IF EXISTS webhook_deliveries;
DROP TABLE IF EXISTS webhook_subscriptions;