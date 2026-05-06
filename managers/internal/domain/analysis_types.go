package domain

import (
	"context"
	"time"
)

// AnalysisRequest represents a request to analyze user input
type AnalysisRequest struct {
	RequestID    string
	UserID       string
	Message      string
	Context      map[string]interface{}
	Timestamp    time.Time
	RequestType  string
	Capabilities []ManagerCapability
	SecurityCtx  *SecurityContext
}

// AnalysisResult represents the result of analyzing a request
type AnalysisResult struct {
	Actions        []ManagerAction
	Intent         string
	Confidence     float64
	EntityType     string
	SecurityStatus string
	Reasoning      string
}

// RequestAnalyzer analyzes user requests
type RequestAnalyzer interface {
	AnalyzeRequest(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)
}
