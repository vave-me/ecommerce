package queries

import (
	"context"
	"log"
	"middleman/managers/internal/domain"
)

type GetManagers struct {
	UserID string
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetManagersHandler struct {
	readModel domain.CatalogRepository
}

func NewGetManagersHandler(readModel domain.CatalogRepository) GetManagersHandler {
	return GetManagersHandler{
		readModel: readModel,
	}
}

func (h GetManagersHandler) GetManagers(ctx context.Context, query GetManagers) ([]*domain.CatalogManager, error) {
	// Log the query parameters
	log.Printf("[GetManagersHandler] Query for managers - UserID: %s, Limit: %d, Offset: %d", 
		query.UserID, query.Limit, query.Offset)
	
	// Call the read model to get all managers
	managers, err := h.readModel.FindAll(ctx, query.UserID)
	if err != nil {
		log.Printf("[GetManagersHandler] Error finding managers: %v", err)
		return nil, err
	}

	log.Printf("[GetManagersHandler] Found %d managers for user %s", len(managers), query.UserID)
	for i, mgr := range managers {
		log.Printf("[GetManagersHandler] Manager[%d]: ID=%s, Name=%s, Type=%s, Active=%t", 
			i, mgr.ID, mgr.Name, mgr.Type, mgr.Active)
	}

	return managers, nil
}
