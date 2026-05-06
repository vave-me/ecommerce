package consciousness

import (
	"context"
	"fmt"
	
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingManager provides comprehensive tracing for the consciousness system
type TracingManager struct {
	tracer  trace.Tracer
	logger  zerolog.Logger
	metrics *MetricsCollector
}

// NewTracingManager creates a new tracing manager
func NewTracingManager(serviceName string, logger zerolog.Logger, metrics *MetricsCollector) *TracingManager {
	tracer := otel.Tracer(serviceName)
	
	return &TracingManager{
		tracer:  tracer,
		logger:  logger,
		metrics: metrics,
	}
}

// StartEventProcessingSpan starts a span for event processing
func (tm *TracingManager) StartEventProcessingSpan(ctx context.Context, eventType, eventID string) (context.Context, trace.Span) {
	ctx, span := tm.tracer.Start(ctx, "consciousness.ProcessEvent",
		trace.WithAttributes(
			attribute.String("event.type", eventType),
			attribute.String("event.id", eventID),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	
	tm.logger.Debug().
		Str("event_type", eventType).
		Str("event_id", eventID).
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Msg("Started event processing span")
	
	return ctx, span
}

// StartPatternDetectionSpan starts a span for pattern detection
func (tm *TracingManager) StartPatternDetectionSpan(ctx context.Context, eventType string) (context.Context, trace.Span) {
	ctx, span := tm.tracer.Start(ctx, "consciousness.DetectPattern",
		trace.WithAttributes(
			attribute.String("event.type", eventType),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	
	return ctx, span
}

// StartDecisionMakingSpan starts a span for decision making
func (tm *TracingManager) StartDecisionMakingSpan(ctx context.Context, patternType string, confidence float64) (context.Context, trace.Span) {
	ctx, span := tm.tracer.Start(ctx, "consciousness.MakeDecision",
		trace.WithAttributes(
			attribute.String("pattern.type", patternType),
			attribute.Float64("pattern.confidence", confidence),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	
	return ctx, span
}

// StartActionExecutionSpan starts a span for action execution
func (tm *TracingManager) StartActionExecutionSpan(ctx context.Context, actionType string, decisionID string) (context.Context, trace.Span) {
	ctx, span := tm.tracer.Start(ctx, "consciousness.ExecuteAction",
		trace.WithAttributes(
			attribute.String("action.type", actionType),
			attribute.String("decision.id", decisionID),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	
	return ctx, span
}

// StartToolExecutionSpan starts a span for tool execution
func (tm *TracingManager) StartToolExecutionSpan(ctx context.Context, toolName string, params map[string]interface{}) (context.Context, trace.Span) {
	ctx, span := tm.tracer.Start(ctx, fmt.Sprintf("tool.%s", toolName),
		trace.WithAttributes(
			attribute.String("tool.name", toolName),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	
	// Add parameter attributes
	for key, value := range params {
		if str, ok := value.(string); ok {
			span.SetAttributes(attribute.String(fmt.Sprintf("tool.param.%s", key), str))
		}
	}
	
	return ctx, span
}

// RecordError records an error on the current span
func (tm *TracingManager) RecordError(ctx context.Context, err error, description string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.RecordError(err,
			trace.WithAttributes(
				attribute.String("error.description", description),
			),
		)
		span.SetStatus(codes.Error, err.Error())
		
		tm.logger.Error().
			Err(err).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String()).
			Str("description", description).
			Msg("Recorded error in span")
	}
}

// AddEvent adds an event to the current span
func (tm *TracingManager) AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// SetSpanAttributes sets attributes on the current span
func (tm *TracingManager) SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// EndSpanWithStatus ends a span with the given status
func (tm *TracingManager) EndSpanWithStatus(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "Success")
	}
	span.End()
}

// ExtractTraceContext extracts trace context from the current span
func (tm *TracingManager) ExtractTraceContext(ctx context.Context) TraceContext {
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()
	
	return TraceContext{
		TraceID: spanCtx.TraceID().String(),
		SpanID:  spanCtx.SpanID().String(),
		Sampled: spanCtx.IsSampled(),
	}
}

// TraceContext holds trace context information
type TraceContext struct {
	TraceID string
	SpanID  string
	Sampled bool
}

// EnhancedLogger provides structured logging with trace context
type EnhancedLogger struct {
	logger  zerolog.Logger
	tracing *TracingManager
}

// NewEnhancedLogger creates a new enhanced logger
func NewEnhancedLogger(logger zerolog.Logger, tracing *TracingManager) *EnhancedLogger {
	return &EnhancedLogger{
		logger:  logger,
		tracing: tracing,
	}
}

// WithContext returns a logger with trace context
func (el *EnhancedLogger) WithContext(ctx context.Context) zerolog.Logger {
	traceCtx := el.tracing.ExtractTraceContext(ctx)
	
	return el.logger.With().
		Str("trace_id", traceCtx.TraceID).
		Str("span_id", traceCtx.SpanID).
		Bool("sampled", traceCtx.Sampled).
		Logger()
}

// Debug logs a debug message with trace context
func (el *EnhancedLogger) Debug(ctx context.Context) *zerolog.Event {
	return el.WithContext(ctx).Debug()
}

// Info logs an info message with trace context
func (el *EnhancedLogger) Info(ctx context.Context) *zerolog.Event {
	return el.WithContext(ctx).Info()
}

// Warn logs a warning message with trace context
func (el *EnhancedLogger) Warn(ctx context.Context) *zerolog.Event {
	return el.WithContext(ctx).Warn()
}

// Error logs an error message with trace context
func (el *EnhancedLogger) Error(ctx context.Context) *zerolog.Event {
	return el.WithContext(ctx).Error()
}

// LogEventProcessing logs event processing details
func (el *EnhancedLogger) LogEventProcessing(ctx context.Context, event interface{}, stage string) {
	logger := el.WithContext(ctx)
	
	logger.Info().
		Str("stage", stage).
		Interface("event", event).
		Msg("Processing event")
		
	// Also add to span
	el.tracing.AddEvent(ctx, fmt.Sprintf("event.%s", stage),
		attribute.String("stage", stage),
	)
}

// LogDecisionMaking logs decision making details
func (el *EnhancedLogger) LogDecisionMaking(ctx context.Context, pattern *Pattern, decision *Decision) {
	logger := el.WithContext(ctx)
	
	if decision != nil {
		logger.Info().
			Str("pattern_type", pattern.Type).
			Float64("confidence", pattern.Confidence).
			Str("decision_id", decision.ID).
			Str("decision_type", decision.Type).
			Int("action_count", len(decision.Actions)).
			Msg("Decision made")
			
		// Add to span
		el.tracing.SetSpanAttributes(ctx,
			attribute.String("decision.id", decision.ID),
			attribute.String("decision.type", decision.Type),
			attribute.Int("decision.action_count", len(decision.Actions)),
		)
	} else {
		logger.Debug().
			Str("pattern_type", pattern.Type).
			Float64("confidence", pattern.Confidence).
			Msg("No decision made")
	}
}

// LogActionExecution logs action execution details
func (el *EnhancedLogger) LogActionExecution(ctx context.Context, action Action, result interface{}, err error) {
	logger := el.WithContext(ctx)
	
	if err != nil {
		logger.Error().
			Err(err).
			Str("action_type", action.Type).
			Interface("parameters", action.Parameters).
			Msg("Action execution failed")
			
		el.tracing.RecordError(ctx, err, "Action execution failed")
	} else {
		logger.Info().
			Str("action_type", action.Type).
			Interface("result", result).
			Msg("Action executed successfully")
			
		el.tracing.AddEvent(ctx, "action.success",
			attribute.String("action.type", action.Type),
		)
	}
}

// LogToolExecution logs tool execution details
func (el *EnhancedLogger) LogToolExecution(ctx context.Context, toolName string, duration float64, success bool, err error) {
	logger := el.WithContext(ctx)
	
	if success {
		logger.Debug().
			Str("tool_name", toolName).
			Float64("duration_ms", duration).
			Msg("Tool executed successfully")
	} else {
		logger.Error().
			Err(err).
			Str("tool_name", toolName).
			Float64("duration_ms", duration).
			Msg("Tool execution failed")
			
		el.tracing.RecordError(ctx, err, fmt.Sprintf("Tool %s failed", toolName))
	}
	
	// Add metric
	el.tracing.SetSpanAttributes(ctx,
		attribute.Float64("tool.duration_ms", duration),
		attribute.Bool("tool.success", success),
	)
}