package tools

import (
	"context"
	"fmt"
	"log"
	"time"

	"middleman/managers/internal/domain"
)

// UserToolService handles all user-related operations with streaming execution
type UserToolService struct {
	userRepository domain.UserRepository
}

// NewUserToolService creates a new user tool service instance
func NewUserToolService(userRepository domain.UserRepository) *UserToolService {
	return &UserToolService{
		userRepository: userRepository,
	}
}

// ExecuteOperation executes user operations with streaming progress updates
func (s *UserToolService) ExecuteOperation(
	ctx context.Context,
	operation string,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (*ToolOperationResult, error) {
	startTime := time.Now()

	// Send initial progress
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 25.0,
		Metadata: map[string]interface{}{
			"step":      "executing_user_operation",
			"operation": operation,
		},
		Timestamp: time.Now().Unix(),
	}

	var result interface{}
	var err error

	switch operation {
	case "find", "get":
		result, err = s.handleFind(ctx, parameters, streamChan, toolID)
	case "get_base_user":
		result, err = s.handleGetBaseUser(ctx, parameters, streamChan, toolID)
	case "list_users":
		result, err = s.handleListUsers(ctx, parameters, streamChan, toolID)
	case "list_participating_users":
		result, err = s.handleListParticipatingUsers(ctx, parameters, streamChan, toolID)
	case "create_user":
		result, err = s.handleCreateUser(ctx, parameters, streamChan, toolID)
	case "update_user":
		result, err = s.handleUpdateUser(ctx, parameters, streamChan, toolID)
	case "login_user":
		result, err = s.handleLoginUser(ctx, parameters, streamChan, toolID)
	default:
		err = fmt.Errorf("unsupported user operation: %s", operation)
	}

	duration := time.Since(startTime)

	if err != nil {
		// Send error update
		streamChan <- ToolExecutionStream{
			ID:       toolID,
			ToolName: "user_operation",
			Status:   "error",
			Error:    err.Error(),
			Metadata: map[string]interface{}{
				"operation": operation,
			},
			Timestamp: time.Now().Unix(),
		}
		return &ToolOperationResult{
			EntityType: "users",
			Operation:  operation,
			Success:    false,
			Error:      err.Error(),
			Duration:   duration,
		}, err
	}

	// Send completion update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "completed",
		Progress: 100,
		Result:   result,
		Metadata: map[string]interface{}{
			"operation": operation,
			"success":   true,
		},
		Timestamp: time.Now().Unix(),
	}

	return &ToolOperationResult{
		EntityType: "users",
		Operation:  operation,
		Success:    true,
		Result:     result,
		Duration:   duration,
	}, nil
}

// handleFind retrieves a user by ID
func (s *UserToolService) handleFind(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID, ok := parameters["id"].(string)
	if !ok {
		userID, ok = parameters["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("id or user_id parameter required")
		}
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "finding_user",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Finding user with ID: %s", userID)
	user, err := s.userRepository.Find(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user find failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "find",
		"result":      user,
		"user_id":     userID,
	}, nil
}

// handleGetBaseUser retrieves base user information
func (s *UserToolService) handleGetBaseUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userID, ok := parameters["id"].(string)
	if !ok {
		userID, ok = parameters["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("id or user_id parameter required")
		}
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "getting_base_user",
			"user_id": userID,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Getting base user with ID: %s", userID)
	baseUser, err := s.userRepository.GetBaseUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get base user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "get_base_user",
		"result":      baseUser,
		"user_id":     userID,
	}, nil
}

// handleListUsers retrieves multiple users by IDs
func (s *UserToolService) handleListUsers(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	userIDs, ok := parameters["user_ids"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("user_ids parameter required for list_users operation")
	}

	// Convert interface slice to string slice
	strUserIDs := make([]string, len(userIDs))
	for i, id := range userIDs {
		if strID, ok := id.(string); ok {
			strUserIDs[i] = strID
		} else {
			return nil, fmt.Errorf("invalid user_id at index %d", i)
		}
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":     "listing_users",
			"user_ids": strUserIDs,
			"count":    len(strUserIDs),
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Listing users with IDs: %v", strUserIDs)
	users, err := s.userRepository.ListUsers(ctx, strUserIDs)
	if err != nil {
		return nil, fmt.Errorf("list users failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "list_users",
		"results":     users,
		"count":       len(users),
		"user_ids":    strUserIDs,
	}, nil
}

// handleListParticipatingUsers retrieves all participating users
func (s *UserToolService) handleListParticipatingUsers(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step": "listing_participating_users",
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Listing participating users")
	users, err := s.userRepository.ListParticipatingUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list participating users failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "list_participating_users",
		"results":     users,
		"count":       len(users),
	}, nil
}

// handleCreateUser creates a new user
func (s *UserToolService) handleCreateUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	email := getUserStringParam(parameters, "email", "")
	password := getUserStringParam(parameters, "password", "")
	username := getUserStringParam(parameters, "username", "")
	firstName := getUserStringParam(parameters, "first_name", "")
	lastName := getUserStringParam(parameters, "last_name", "")
	location := getUserStringParam(parameters, "location", "")
	lat := getUserFloat32Param(parameters, "lat", 0.0)
	lng := getUserFloat32Param(parameters, "lng", 0.0)
	thumbnail := getUserStringParam(parameters, "thumbnail", "")
	language := getUserStringParam(parameters, "language", "")

	if email == "" || password == "" || username == "" {
		return nil, fmt.Errorf("email, password, and username are required for create_user operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":     "creating_user",
			"email":    email,
			"username": username,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Creating user with email: %s, username: %s", email, username)
	userID, err := s.userRepository.CreateUser(ctx, email, password, username, firstName, lastName, location, lat, lng, thumbnail, language)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "create_user",
		"result":      map[string]interface{}{"id": userID},
		"user_id":     userID,
	}, nil
}

// handleUpdateUser updates an existing user
func (s *UserToolService) handleUpdateUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	id := getUserStringParam(parameters, "id", "")
	username := getUserStringParam(parameters, "username", "")
	firstName := getUserStringParam(parameters, "first_name", "")
	lastName := getUserStringParam(parameters, "last_name", "")
	bio := getUserStringParam(parameters, "bio", "")
	privacy := getUserStringParam(parameters, "privacy", "")
	background := getUserStringParam(parameters, "background", "")
	location := getUserStringParam(parameters, "location", "")
	lat := getUserFloat32Param(parameters, "lat", 0.0)
	lng := getUserFloat32Param(parameters, "lng", 0.0)
	thumbnail := getUserStringParam(parameters, "thumbnail", "")

	if id == "" {
		return nil, fmt.Errorf("id parameter required for update_user operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":    "updating_user",
			"user_id": id,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Updating user with ID: %s", id)
	updatedUserID, err := s.userRepository.UpdateUser(ctx, id, username, firstName, lastName, bio, privacy, background, location, lat, lng, thumbnail)
	if err != nil {
		return nil, fmt.Errorf("update user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "update_user",
		"result":      map[string]interface{}{"user_id": updatedUserID},
		"user_id":     updatedUserID,
	}, nil
}

// handleLoginUser authenticates a user
func (s *UserToolService) handleLoginUser(
	ctx context.Context,
	parameters map[string]interface{},
	streamChan chan<- ToolExecutionStream,
	toolID string,
) (interface{}, error) {
	email := getUserStringParam(parameters, "email", "")
	password := getUserStringParam(parameters, "password", "")

	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required for login_user operation")
	}

	// Send progress update
	streamChan <- ToolExecutionStream{
		ID:       toolID,
		ToolName: "user_operation",
		Status:   "progress",
		Progress: 75.0,
		Metadata: map[string]interface{}{
			"step":  "logging_in_user",
			"email": email,
		},
		Timestamp: time.Now().Unix(),
	}

	log.Printf("UserToolService: Logging in user with email: %s", email)
	loginResponse, err := s.userRepository.LoginUser(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("login user failed: %w", err)
	}

	return map[string]interface{}{
		"entity_type": "users",
		"operation":   "login_user",
		"result":      loginResponse,
		"email":       email,
	}, nil
}

// Helper functions for parameter extraction
func getUserStringParam(params map[string]interface{}, key, defaultValue string) string {
	if val, ok := params[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getUserFloat32Param(params map[string]interface{}, key string, defaultValue float32) float32 {
	if val, ok := params[key]; ok {
		switch v := val.(type) {
		case float64:
			return float32(v)
		case float32:
			return v
		case int:
			return float32(v)
		case int64:
			return float32(v)
		}
	}
	return defaultValue
}
