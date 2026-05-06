package queries

import (
	"context"
	"middleman/offers/internal/domain"
)

type ListOffers struct {
	UserSellerID   string
	UserCustomerID string
	OfferStatus    string
	Limit          int
	Offset         int
}

type ListOffersHandler struct {
	middleman domain.MiddlemanRepository
}

func NewListOffersHandler(middleman domain.MiddlemanRepository) ListOffersHandler {
	return ListOffersHandler{middleman: middleman}
}

func (h ListOffersHandler) ListOffers(ctx context.Context, query ListOffers) ([]*domain.Offer, int64, error) {
	// If UserSellerID or UserCustomerID is provided, use the All method
	// Otherwise, we need to extend the repository with more query methods
	var middlemanOffers []*domain.MiddlemanOffer
	var err error
	
	if query.UserSellerID != "" {
		middlemanOffers, err = h.middleman.All(ctx, query.UserSellerID)
	} else if query.UserCustomerID != "" {
		middlemanOffers, err = h.middleman.All(ctx, query.UserCustomerID)
	} else {
		// For now, if no user filter is provided, return empty
		// In production, you'd want to add a method to get all offers with pagination
		return []*domain.Offer{}, 0, nil
	}
	
	if err != nil {
		return nil, 0, err
	}
	
	// Convert from read model to domain model
	offers := make([]*domain.Offer, 0, len(middlemanOffers))
	for _, mo := range middlemanOffers {
		// Apply additional filters if needed
		if query.OfferStatus != "" {
			// Skip if status doesn't match (you'd need to store status in read model)
			continue
		}
		
		offer := &domain.Offer{
			Aggregate:      domain.NewOffer(mo.ID).Aggregate,
			UserSellerID:   mo.UserSellerID,
			UserCustomerID: mo.UserCustomerID,
			ProductID:      mo.ProductID,
			Status:         domain.OfferStatusActive, // Default status
		}
		offers = append(offers, offer)
	}
	
	// Apply pagination
	total := int64(len(offers))
	start := query.Offset
	end := query.Offset + query.Limit
	
	if start > len(offers) {
		start = len(offers)
	}
	if end > len(offers) || query.Limit == 0 {
		end = len(offers)
	}
	
	paginatedOffers := offers[start:end]
	
	return paginatedOffers, total, nil
}
