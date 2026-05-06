package queries

import (
	"context"
	"middleman/media/internal/domain"
	"time"
)

type (
	GetImportStatus struct {
		SessionID string
	}

	GetImportStatusHandler struct {
		imports domain.ImportSessionRepository
		items   domain.ImportItemRepository
	}
)

func NewGetImportStatusHandler(imports domain.ImportSessionRepository, items domain.ImportItemRepository) GetImportStatusHandler {
	return GetImportStatusHandler{
		imports: imports,
		items:   items,
	}
}

func (h GetImportStatusHandler) GetImportStatus(ctx context.Context, query GetImportStatus) (*ImportStatus, error) {
	session, err := h.imports.Get(ctx, query.SessionID)
	if err != nil {
		return nil, err
	}

	// Get recent errors
	failedItems, err := h.items.GetBySession(ctx, query.SessionID, domain.ImportItemStatusFailed)
	if err != nil {
		return nil, err
	}

	var recentErrors []*ImportError
	for i, item := range failedItems {
		if i >= 10 { // Limit to 10 most recent errors
			break
		}
		recentErrors = append(recentErrors, &ImportError{
			ExternalID:   item.ExternalID,
			ErrorCode:    "IMPORT_FAILED",
			ErrorMessage: item.ErrorMessage,
		})
	}

	// Calculate progress and ETA
	progressPercentage := float32(0)
	if session.TotalImages > 0 {
		progressPercentage = float32(session.ProcessedImages) / float32(session.TotalImages) * 100
	}

	var estimatedCompletionTime int64
	if session.ProcessedImages > 0 && progressPercentage < 100 {
		elapsed := time.Since(session.StartedAt)
		avgTimePerImage := elapsed / time.Duration(session.ProcessedImages)
		remainingImages := session.TotalImages - session.ProcessedImages
		remainingTime := avgTimePerImage * time.Duration(remainingImages)
		estimatedCompletionTime = time.Now().Add(remainingTime).Unix()
	}

	return &ImportStatus{
		SessionID:               session.ID,
		Status:                  string(session.Status),
		TotalImages:             session.TotalImages,
		ProcessedImages:         session.ProcessedImages,
		FailedImages:            session.FailedImages,
		RecentErrors:            recentErrors,
		ProgressPercentage:      progressPercentage,
		EstimatedCompletionTime: estimatedCompletionTime,
	}, nil
}

type ImportStatus struct {
	SessionID               string
	Status                  string
	TotalImages             int32
	ProcessedImages         int32
	FailedImages            int32
	RecentErrors            []*ImportError
	ProgressPercentage      float32
	EstimatedCompletionTime int64
}

type ImportError struct {
	ExternalID   string
	ErrorCode    string
	ErrorMessage string
}