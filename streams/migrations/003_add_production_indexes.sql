-- +migrate Up

-- Performance indexes for live_streams
CREATE INDEX IF NOT EXISTS idx_live_streams_status 
    ON live_streams(status) 
    WHERE status IN ('scheduled', 'active');

CREATE INDEX IF NOT EXISTS idx_live_streams_scheduled_time 
    ON live_streams(scheduled_start_time) 
    WHERE status = 'scheduled';

CREATE INDEX IF NOT EXISTS idx_live_streams_user_seller 
    ON live_streams(user_seller_id, status);

CREATE INDEX IF NOT EXISTS idx_live_streams_category 
    ON live_streams(category_id, status);

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_live_streams_active_category_time 
    ON live_streams(category_id, scheduled_start_time) 
    WHERE status = 'active';

-- Viewer session indexes
CREATE INDEX IF NOT EXISTS idx_viewer_sessions_stream_user 
    ON viewer_sessions(stream_id, user_id) 
    WHERE left_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_viewer_sessions_user_activity 
    ON viewer_sessions(user_id, joined_at);

-- Stream health indexes
CREATE INDEX IF NOT EXISTS idx_stream_health_recent 
    ON stream_health(stream_id, updated_at DESC);

-- CDN endpoint indexes
CREATE INDEX IF NOT EXISTS idx_cdn_endpoints_stream_provider 
    ON cdn_endpoints(stream_id, provider);

-- DRM config indexes
CREATE INDEX IF NOT EXISTS idx_drm_configs_stream 
    ON drm_configs(stream_id);

-- Webhook delivery optimization
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_retry_pending 
    ON webhook_deliveries(next_retry_at) 
    WHERE status IN ('pending', 'retrying') AND next_retry_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_subscription_status 
    ON webhook_deliveries(subscription_id, status, created_at DESC);

-- Text search indexes (if using PostgreSQL full text search)
CREATE INDEX IF NOT EXISTS idx_live_streams_title_search 
    ON live_streams USING gin(to_tsvector('english', title));

CREATE INDEX IF NOT EXISTS idx_live_streams_description_search 
    ON live_streams USING gin(to_tsvector('english', description));

-- +migrate Down

DROP INDEX IF EXISTS idx_live_streams_status;
DROP INDEX IF EXISTS idx_live_streams_scheduled_time;
DROP INDEX IF EXISTS idx_live_streams_user_seller;
DROP INDEX IF EXISTS idx_live_streams_category;
DROP INDEX IF EXISTS idx_live_streams_active_category_time;
DROP INDEX IF EXISTS idx_viewer_sessions_stream_user;
DROP INDEX IF EXISTS idx_viewer_sessions_user_activity;
DROP INDEX IF EXISTS idx_stream_health_recent;
DROP INDEX IF EXISTS idx_cdn_endpoints_stream_provider;
DROP INDEX IF EXISTS idx_drm_configs_stream;
DROP INDEX IF EXISTS idx_webhook_deliveries_retry_pending;
DROP INDEX IF EXISTS idx_webhook_deliveries_subscription_status;
DROP INDEX IF EXISTS idx_live_streams_title_search;
DROP INDEX IF EXISTS idx_live_streams_description_search;