package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"middleman/internal/auth"
	"middleman/internal/rpc"
	"middleman/managers/internal/domain"
	schedulerpb "middleman/scheduler/schedulerspb"
)

type schedulerRepository struct {
	endpoint string
	auth     *auth.Auth // Optional auth for authenticated calls
}

var _ domain.SchedulerRepository = (*schedulerRepository)(nil)

// NewSchedulerRepository creates a new scheduler repository with proper gRPC client
func NewSchedulerRepository(endpoint string, authInstance *auth.Auth) domain.SchedulerRepository {
	return &schedulerRepository{
		endpoint: endpoint,
		auth:     authInstance,
	}
}

// dial creates a connection to the scheduler service
func (r *schedulerRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}

// dialWithAuth creates an authenticated connection to the scheduler service
func (r *schedulerRepository) dialWithAuth(ctx context.Context) (*grpc.ClientConn, error) {
	// Use authenticated dial if auth is available, otherwise use regular dial
	return rpc.DialWithAuth(ctx, r.endpoint, r.auth)
}

// ScheduleAction schedules a new action to be executed at a specific time
func (r *schedulerRepository) ScheduleAction(ctx context.Context, action *domain.ScheduledAction) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("scheduler service unavailable, action not scheduled")
		return nil // Graceful degradation
	}
	defer conn.Close()

	client := schedulerpb.NewSchedulerServiceClient(conn)

	// Convert domain model to protobuf
	req := &schedulerpb.ScheduleTaskRequest{
		ManagerId:   action.EntityID,
		TaskType:    action.Action,
		ScheduledAt: action.ScheduledAt.Unix(),
		Payload:     convertParametersToPayload(action.Parameters),
	}

	_, err = client.ScheduleTask(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			log.Warn().Err(err).Msg("scheduler service unavailable")
			return nil // Graceful degradation
		}
		return fmt.Errorf("failed to schedule action: %w", err)
	}

	return nil
}

// CancelScheduledAction cancels a previously scheduled action
func (r *schedulerRepository) CancelScheduledAction(ctx context.Context, actionID string) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("scheduler service unavailable, cannot cancel action")
		return nil // Graceful degradation
	}
	defer conn.Close()

	client := schedulerpb.NewSchedulerServiceClient(conn)

	req := &schedulerpb.CancelTaskRequest{
		TaskId: actionID,
	}

	_, err = client.CancelTask(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			log.Warn().Err(err).Msg("scheduler service unavailable")
			return nil // Graceful degradation
		}
		return fmt.Errorf("failed to cancel scheduled action: %w", err)
	}

	return nil
}

// GetScheduledActions retrieves scheduled actions for a specific entity
func (r *schedulerRepository) GetScheduledActions(ctx context.Context, entityID string, entityType string) ([]*domain.ScheduledAction, error) {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("scheduler service unavailable, returning empty results")
		return []*domain.ScheduledAction{}, nil // Graceful degradation
	}
	defer conn.Close()

	client := schedulerpb.NewSchedulerServiceClient(conn)

	req := &schedulerpb.GetScheduledTasksRequest{
		ManagerId: entityID,
	}

	resp, err := client.GetScheduledTasks(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			log.Warn().Err(err).Msg("scheduler service unavailable")
			return []*domain.ScheduledAction{}, nil // Graceful degradation
		}
		return nil, fmt.Errorf("failed to get scheduled actions: %w", err)
	}

	// Convert protobuf response to domain models
	actions := make([]*domain.ScheduledAction, 0, len(resp.Tasks))
	for _, task := range resp.Tasks {
		actions = append(actions, &domain.ScheduledAction{
			ID:          task.Id,
			Name:        task.TaskType,
			EntityID:    task.ManagerId,
			EntityType:  entityType,
			Action:      task.TaskType,
			Parameters:  convertPayloadToParameters(task.Payload),
			ScheduledAt: time.Unix(task.ScheduledAt, 0),
			Status:      mapTaskStatus(task.Status),
			CreatedAt:   time.Unix(task.CreatedAt, 0),
			UpdatedAt:   time.Unix(task.UpdatedAt, 0),
		})
	}

	return actions, nil
}

// UpdateScheduledAction updates an existing scheduled action
func (r *schedulerRepository) UpdateScheduledAction(ctx context.Context, actionID string, updates *domain.ScheduledActionUpdate) error {
	conn, err := r.dialWithAuth(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("scheduler service unavailable, cannot update action")
		return nil // Graceful degradation
	}
	defer conn.Close()

	client := schedulerpb.NewSchedulerServiceClient(conn)

	req := &schedulerpb.UpdateTaskRequest{
		TaskId: actionID,
	}

	// Apply updates
	if updates.ScheduledAt != nil {
		req.ScheduledAt = updates.ScheduledAt.Unix()
	}
	if updates.Parameters != nil {
		req.Payload = convertParametersToPayload(updates.Parameters)
	}

	_, err = client.UpdateTask(ctx, req)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			log.Warn().Err(err).Msg("scheduler service unavailable")
			return nil // Graceful degradation
		}
		return fmt.Errorf("failed to update scheduled action: %w", err)
	}

	return nil
}

// Health checks if scheduler service is available
func (r *schedulerRepository) Health(ctx context.Context) bool {
	conn, err := r.dial(ctx)
	if err != nil {
		return false
	}
	defer conn.Close()

	client := schedulerpb.NewSchedulerServiceClient(conn)

	// Try a simple operation to check health
	_, err = client.GetScheduledTasks(ctx, &schedulerpb.GetScheduledTasksRequest{
		ManagerId: "health-check",
	})

	return err == nil || status.Code(err) != codes.Unavailable
}

// Helper functions

func convertParametersToPayload(params map[string]interface{}) map[string]string {
	payload := make(map[string]string)
	for k, v := range params {
		payload[k] = fmt.Sprintf("%v", v)
	}
	return payload
}

func convertPayloadToParameters(payload map[string]string) map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range payload {
		params[k] = v
	}
	return params
}

func mapTaskStatus(pbStatus schedulerpb.TaskStatus) string {
	switch pbStatus {
	case schedulerpb.TaskStatus_PENDING:
		return domain.SchedulerActionStatusPending
	case schedulerpb.TaskStatus_RUNNING:
		return domain.SchedulerActionStatusExecuting
	case schedulerpb.TaskStatus_COMPLETED:
		return domain.SchedulerActionStatusCompleted
	case schedulerpb.TaskStatus_FAILED:
		return domain.SchedulerActionStatusFailed
	case schedulerpb.TaskStatus_CANCELLED:
		return domain.SchedulerActionStatusCancelled
	default:
		return domain.SchedulerActionStatusPending
	}
}
