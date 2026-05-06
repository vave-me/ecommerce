package services

import (
	"middleman/managers/internal/models"
)

// OperationType represents the type of operation to perform
type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationRead   OperationType = "read"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
	OperationSearch OperationType = "search"
	OperationList   OperationType = "list"
)

// RepositoryQuery represents a query to execute against a repository
type RepositoryQuery struct {
	EntityType models.EntityType
	Operation  OperationType
	Parameters map[string]interface{}
}

// RepositoryResponse represents a response from a repository operation
type RepositoryResponse struct {
	Data     interface{}
	Metadata ResponseMetadata
}

// ResponseMetadata contains metadata about the response
type ResponseMetadata struct {
	EntityType models.EntityType
	Operation  OperationType
}

// QueryParameters represents common query parameters
type QueryParameters struct {
	ID          string
	UserID      string
	ItemID      string
	SearchTerm  string
	Name        string
	Description string
	CategoryID  string
	MinPrice    int64
	MaxPrice    int64
	Brand       string
	Condition   string
	Model       string
	Status      string
	Page        int
	PageSize    int
	SortBy      string
	SortOrder   string
}