package services

import (
	"context"
	"middleman/managers/internal/domain"
)

// RequestAnalyzerAdapter adapts UnifiedRequestAnalyzer to implement domain.RequestAnalyzer
type RequestAnalyzerAdapter struct {
	analyzer *UnifiedRequestAnalyzer
}

// NewRequestAnalyzerAdapter creates a new adapter
func NewRequestAnalyzerAdapter() *RequestAnalyzerAdapter {
	return &RequestAnalyzerAdapter{
		analyzer: NewUnifiedRequestAnalyzer(),
	}
}

// AnalyzeRequest implements domain.RequestAnalyzer
func (r *RequestAnalyzerAdapter) AnalyzeRequest(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error) {
	// Convert domain request to services request
	servicesReq := AnalysisRequest{
		RequestID:    req.RequestID,
		UserID:       req.UserID,
		Message:      req.Message,
		Context:      req.Context,
		Timestamp:    req.Timestamp,
		RequestType:  req.RequestType,
		Capabilities: req.Capabilities,
		SecurityCtx:  convertSecurityContext(req.SecurityCtx),
	}

	// Call the services analyzer
	result, err := r.analyzer.AnalyzeRequest(ctx, servicesReq)
	if err != nil {
		return nil, err
	}

	// Convert services result to domain result
	return &domain.AnalysisResult{
		Actions:        result.Actions,
		Intent:         result.Intent,
		Confidence:     result.Confidence,
		EntityType:     result.EntityType,
		SecurityStatus: result.SecurityStatus,
		Reasoning:      result.Reasoning,
	}, nil
}

// convertCapabilities is no longer needed since we use domain types directly
// Kept for backward compatibility but just returns the input
func convertCapabilities(caps []domain.ManagerCapability) []domain.ManagerCapability {
	return caps
}

// convertSecurityContext is no longer needed since we use domain types directly
// Kept for backward compatibility but just returns the input
func convertSecurityContext(ctx *domain.SecurityContext) *domain.SecurityContext {
	if ctx == nil {
		return nil
	}
	return ctx
}

// convertActions is no longer needed since we use domain types directly
// Kept for backward compatibility but just returns the input
func convertActions(actions []domain.ManagerAction) []domain.ManagerAction {
	return actions
}
