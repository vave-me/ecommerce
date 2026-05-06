package domain

import "context"

type MiddlemanMedia struct {
	ID         string
	ItemID     string
	ItemType   ItemType
	UserID     string
	Status     MediaStatus
	MediaOrder map[int]MediaOrderItem
}

type MiddlemanMediaRepository interface {
	AddMedia(ctx context.Context, id, itemID string, itemType ItemType, userID string, status MediaStatus) error
	UpdateMedia(ctx context.Context, id, itemID string, itemType ItemType, userID string, status MediaStatus) error
	GetMedia(ctx context.Context, id string) (*MiddlemanMedia, error)
	AddMediaItemOrder(ctx context.Context, id, mediaItemId, mediaItemUrl string, displayOrder int) error
	GetMediaByItem(ctx context.Context, itemID string) (*MiddlemanMedia, error)
	GetAllUserMedia(ctx context.Context, userID string) ([]*MiddlemanMedia, error)
	GetItemMedia(ctx context.Context, id string) (*MiddlemanMedia, error)
	RemoveMedia(ctx context.Context, id string) error
}
