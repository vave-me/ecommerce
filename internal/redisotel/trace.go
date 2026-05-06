// File: redisotel/redisotel.go
package redisotel

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RedisDB defines the subset of Redis commands to be traced.
// Extend this interface with additional methods as needed.
type RedisDB interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Ping(ctx context.Context) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	// Add other Redis commands you wish to trace here
}

// tracedRedisClient wraps the Redis client to add OpenTelemetry tracing.
type tracedRedisClient struct {
	db redis.Client
}

// Ensure tracedRedisClient implements RedisDB interface.
var _ RedisDB = (*tracedRedisClient)(nil)

// Trace wraps an existing redis.Client with tracing capabilities.
func Trace(client redis.Client) RedisDB {
	return tracedRedisClient{db: client}
}

// Get wraps the Redis GET command with tracing.
func (t tracedRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.AddEvent("Get", trace.WithAttributes(
		attribute.String("Redis.Command", "GET"),
		attribute.String("Redis.Key", key),
	))

	cmd := t.db.Get(ctx, key)

	start := time.Now()
	_, err := cmd.Result()
	duration := time.Since(start).Seconds()

	span.AddEvent("GetResult", trace.WithAttributes(
		attribute.Float64("Redis.Took", duration),
	))

	t.recordError(span, err)

	return cmd
}

// Set wraps the Redis SET command with tracing.
func (t tracedRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.AddEvent("Set", trace.WithAttributes(
		attribute.String("Redis.Command", "SET"),
		attribute.String("Redis.Key", key),
		attribute.Float64("Redis.ExpirationSeconds", expiration.Seconds()),
	))

	cmd := t.db.Set(ctx, key, value, expiration)

	start := time.Now()
	_, err := cmd.Result()
	duration := time.Since(start).Seconds()

	span.AddEvent("SetResult", trace.WithAttributes(
		attribute.Float64("Redis.Took", duration),
	))

	t.recordError(span, err)

	return cmd
}

// Ping wraps the Redis PING command with tracing.
func (t tracedRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.AddEvent("Ping", trace.WithAttributes(
		attribute.String("Redis.Command", "PING"),
	))

	cmd := t.db.Ping(ctx)

	start := time.Now()
	_, err := cmd.Result()
	duration := time.Since(start).Seconds()

	span.AddEvent("PingResult", trace.WithAttributes(
		attribute.Float64("Redis.Took", duration),
	))

	t.recordError(span, err)

	return cmd
}

// Del wraps the Redis DEL command with tracing.
func (t tracedRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.AddEvent("Del", trace.WithAttributes(
		attribute.String("Redis.Command", "DEL"),
		attribute.StringSlice("Redis.Keys", keys),
	))

	cmd := t.db.Del(ctx, keys...)

	start := time.Now()
	_, err := cmd.Result()
	duration := time.Since(start).Seconds()

	span.AddEvent("DelResult", trace.WithAttributes(
		attribute.Float64("Redis.Took", duration),
	))

	t.recordError(span, err)

	return cmd
}

// Incr wraps the Redis INCR command with tracing.
func (t tracedRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.AddEvent("Incr", trace.WithAttributes(
		attribute.String("Redis.Command", "INCR"),
		attribute.String("Redis.Key", key),
	))

	cmd := t.db.Incr(ctx, key)

	start := time.Now()
	_, err := cmd.Result()
	duration := time.Since(start).Seconds()

	span.AddEvent("IncrResult", trace.WithAttributes(
		attribute.Float64("Redis.Took", duration),
	))

	t.recordError(span, err)

	return cmd
}

// recordError records errors in the span with relevant attributes.
func (t tracedRedisClient) recordError(span trace.Span, err error) {
	if err != nil {
		var redisErr *redis.Error
		if errors.As(err, &redisErr) {
			span.AddEvent("Redis Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
				attribute.String("ErrorType", "Nil"),
				attribute.String("ErrorMessage", "Key does not exist"),
			))
		} else {
			span.AddEvent("Redis Error", trace.WithAttributes(
				attribute.String("Error", err.Error()),
			))
		}
	}
}
