package models

import "time"

// Core Media Structures
type Media struct {
	ID         string           `json:"id"`
	ItemID     string           `json:"item_id"`
	ItemType   string           `json:"item_type"`
	UserID     string           `json:"user_id"`
	FileType   string           `json:"file_type"`
	MediaOrder []MediaOrderItem `json:"media_order"`
	// Legacy field for backward compatibility
	Status string `json:"status,omitempty"`
}

type MediaOrderItem struct {
	MediaItemID string `json:"media_item_id"`
	URL         string `json:"url"`
}

type Image struct {
	ID           string    `json:"id"`
	MediaID      string    `json:"media_id"`
	DisplayOrder int32     `json:"display_order"`
	IsMain       bool      `json:"is_main"`
	URL          string    `json:"url"`
	Metadata     string    `json:"metadata"`
	FileType     string    `json:"file_type"`
	Thumbnail    string    `json:"thumbnail"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type Video struct {
	ID           string    `json:"id"`
	MediaID      string    `json:"media_id"`
	DisplayOrder int32     `json:"display_order"`
	IsMain       bool      `json:"is_main"`
	URL          string    `json:"url"`
	Metadata     string    `json:"metadata"`
	FileType     string    `json:"file_type"`
	Thumbnail    string    `json:"thumbnail"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Response Types

// Media Management Responses
type CreateMediaResponse struct {
	ID string `json:"id"`
}

type UpdateMediaResponse struct {
	ID string `json:"id"`
}

type RemoveMediaResponse struct {
	ID string `json:"id"`
}

// Image Management Responses
type AddImageResponse struct {
	URL     string `json:"url"`
	ViewURL string `json:"view_url"`
}

type RemoveImageResponse struct {
	ID string `json:"id"`
}

type GetAllItemImagesResponse struct {
	Images []Image `json:"images"`
}

type GetAllMediaImagesResponse struct {
	Images []Image `json:"images"`
}

type GetAllImagesByAuthorResponse struct {
	Images      []Image `json:"images"`
	TotalCount  int64   `json:"total_count"`
	CurrentPage int64   `json:"current_page"`
	TotalPages  int64   `json:"total_pages"`
}

// Video Management Responses
type AddVideoResponse struct {
	URL     string `json:"url"`
	ViewURL string `json:"view_url"`
}

type RemoveVideoResponse struct {
	ID string `json:"id"`
}

type GetAllItemVideosResponse struct {
	Videos []Video `json:"videos"`
}

type GetAllMediaVideosResponse struct {
	Videos []Video `json:"videos"`
}

type GetAllVideosResponse struct {
	Videos      []Video `json:"videos"`
	TotalCount  int64   `json:"total_count"`
	CurrentPage int64   `json:"current_page"`
	TotalPages  int64   `json:"total_pages"`
}

type GetAllVideosByAuthorResponse struct {
	Videos      []Video `json:"videos"`
	TotalCount  int64   `json:"total_count"`
	CurrentPage int64   `json:"current_page"`
	TotalPages  int64   `json:"total_pages"`
}
