package grpc

import (
	"context"
	"time"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/scheduler/internal/application"
	"middleman/scheduler/internal/application/commands"
	"middleman/scheduler/internal/application/queries"
	"middleman/scheduler/internal/domain"
	"middleman/scheduler/schedulerspb"
)

type server struct {
	app application.App
	schedulerspb.UnimplementedSchedulerServiceServer
}

var _ schedulerspb.SchedulerServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	schedulerspb.RegisterSchedulerServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateScheduler(ctx context.Context, request *schedulerspb.CreateSchedulerRequest) (*schedulerspb.CreateSchedulerResponse, error) {
	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	schedulerID := uuid.New().String()

	span.SetAttributes(
		attribute.String("SchedulerID", schedulerID),
	)

	err := s.app.CreateScheduler(ctx, commands.CreateScheduler{
		ID:     schedulerID,
		UserID: userID,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &schedulerspb.CreateSchedulerResponse{
		Id:     schedulerID,
		UserId: userID,
	}, nil
}
func (s server) AddAction(ctx context.Context, request *schedulerspb.AddActionRequest) (*schedulerspb.AddActionResponse, error) {
	span := trace.SpanFromContext(ctx)
	id := uuid.New().String()

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	// Input validation
	schedulerID := request.GetSchedulerId()
	if schedulerID == "" {
		return nil, status.Error(grpc_code.InvalidArgument, "scheduler_id is required")
	}
	
	task := request.GetNaturalLanguageTask()
	if task == "" {
		return nil, status.Error(grpc_code.InvalidArgument, "natural_language_task is required")
	}
	
	if request.GetExecutionTime() == nil {
		return nil, status.Error(grpc_code.InvalidArgument, "execution_time is required")
	}
	executionTime := request.GetExecutionTime().AsTime()

	// Authorization: verify user owns the scheduler
	scheduler, err := s.app.GetScheduler(ctx, queries.GetScheduler{UserID: userID})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		return nil, status.Error(grpc_code.NotFound, "scheduler not found")
	}
	if scheduler.ID != schedulerID {
		return nil, status.Error(grpc_code.PermissionDenied, "access denied")
	}

	span.SetAttributes(
		attribute.String("SchedulerID", schedulerID),
		attribute.String("Task", task),
		attribute.String("ExecutionTime", executionTime.String()),
		attribute.String("RequestID", id),
		attribute.String("UserID", userID),
	)
	err = s.app.AddAction(ctx, commands.AddAction{
		ID:                  id,
		SchedulerID:         schedulerID,
		NaturalLanguageTask: task,
		ExecutionTime:       executionTime,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &schedulerspb.AddActionResponse{
		Id: id,
	}, nil
}

func (s server) GetScheduler(ctx context.Context, request *schedulerspb.GetSchedulerRequest) (*schedulerspb.GetSchedulerResponse, error) {
	span := trace.SpanFromContext(ctx)
	claims, ok := auth.ClaimsFromContext(ctx)

	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	userID := claims.Subject

	span.SetAttributes(
		attribute.String("UserID", userID),
	)
	// TODO add check with claims if the userID is the same from response and request
	scheduler, err := s.app.GetScheduler(ctx, queries.GetScheduler{UserID: userID})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &schedulerspb.GetSchedulerResponse{SchedulerId: scheduler.ID, UserId: scheduler.UserID}, nil
}

func (s server) GetAction(ctx context.Context, request *schedulerspb.GetActionRequest) (*schedulerspb.GetActionResponse, error) {
	span := trace.SpanFromContext(ctx)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span.SetAttributes(
		attribute.String("ActionID", request.GetId()),
	)

	interaction, err := s.app.GetAction(ctx, queries.GetAction{ID: request.GetId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &schedulerspb.GetActionResponse{Action: s.interactionFromDomain(interaction)}, nil
}

func (s server) RemoveAction(ctx context.Context, request *schedulerspb.RemoveActionRequest) (*schedulerspb.RemoveActionResponse, error) {
	span := trace.SpanFromContext(ctx)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span.SetAttributes(
		attribute.String("ActionID", request.GetId()),
	)

	err := s.app.RemoveAction(ctx, commands.RemoveAction{
		ID: request.GetId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
	}

	return &schedulerspb.RemoveActionResponse{}, err
}

func (s server) GetActions(ctx context.Context, request *schedulerspb.GetActionsRequest) (*schedulerspb.GetActionsResponse, error) {
	span := trace.SpanFromContext(ctx)

	_, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}

	span.SetAttributes(
		attribute.String("SchedulerID", request.GetSchedulerId()),
	)

	interactions, err := s.app.GetActions(ctx, queries.GetActions{SchedulerID: request.GetSchedulerId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	protoActions := make([]*schedulerspb.Action, len(interactions))
	for i, interaction := range interactions {
		protoActions[i] = s.interactionFromDomain(interaction)
	}

	return &schedulerspb.GetActionsResponse{
		Actions: protoActions,
	}, nil
}

func (s server) interactionFromDomain(action *domain.MiddlemanAction) *schedulerspb.Action {
	result := &schedulerspb.Action{
		Id:                  action.ID,
		SchedulerId:         action.SchedulerID,
		NaturalLanguageTask: action.NaturalLanguageTask,
		ExecutionTime:       timestamppb.New(action.ExecutionTime),
		Status:              action.Status,
		CreatedAt:           timestamppb.New(action.CreatedAt),
		Result:              action.Result,
		ErrorMessage:        action.ErrorMessage,
	}
	
	if action.ExecutedAt != nil {
		result.ExecutedAt = timestamppb.New(*action.ExecutedAt)
	}
	
	return result
}

func (s server) schedulerFromDomain(scheduler *domain.MiddlemanScheduler) *schedulerspb.Scheduler {
	return &schedulerspb.Scheduler{
		Id:     scheduler.ID,
		UserId: scheduler.UserID,
	}
}

// ScheduleTask schedules a new task for a manager
func (s server) ScheduleTask(ctx context.Context, request *schedulerspb.ScheduleTaskRequest) (*schedulerspb.ScheduleTaskResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	// Generate task ID
	taskID := uuid.New().String()
	
	span.SetAttributes(
		attribute.String("TaskID", taskID),
		attribute.String("ManagerID", request.GetManagerId()),
		attribute.String("TaskType", request.GetTaskType()),
	)
	
	// Convert Unix timestamp to time.Time
	scheduledAt := time.Unix(request.GetScheduledAt(), 0)
	
	// Schedule the task
	err := s.app.ScheduleTask(ctx, commands.ScheduleTask{
		ID:          taskID,
		ManagerID:   request.GetManagerId(),
		TaskType:    request.GetTaskType(),
		ScheduledAt: scheduledAt,
		Payload:     request.GetPayload(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	return &schedulerspb.ScheduleTaskResponse{
		TaskId: taskID,
	}, nil
}

// CancelTask cancels a scheduled task
func (s server) CancelTask(ctx context.Context, request *schedulerspb.CancelTaskRequest) (*schedulerspb.CancelTaskResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("TaskID", request.GetTaskId()),
	)
	
	// Cancel the task
	err := s.app.CancelTask(ctx, commands.CancelTask{
		TaskID: request.GetTaskId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return &schedulerspb.CancelTaskResponse{Success: false}, err
	}
	
	return &schedulerspb.CancelTaskResponse{Success: true}, nil
}

// GetScheduledTasks retrieves tasks for a manager
func (s server) GetScheduledTasks(ctx context.Context, request *schedulerspb.GetScheduledTasksRequest) (*schedulerspb.GetScheduledTasksResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("ManagerID", request.GetManagerId()),
	)
	
	// Build filter
	filter := domain.TaskFilter{
		Limit: int(request.GetLimit()),
	}
	
	// Handle status filter
	if request.GetStatus() != schedulerspb.TaskStatus_TASK_STATUS_UNSPECIFIED {
		status := s.taskStatusToString(request.GetStatus())
		filter.Status = &status
	}
	
	// Handle time filters
	if request.GetScheduledAfter() > 0 {
		t := time.Unix(request.GetScheduledAfter(), 0)
		filter.ScheduledAfter = &t
	}
	if request.GetScheduledBefore() > 0 {
		t := time.Unix(request.GetScheduledBefore(), 0)
		filter.ScheduledBefore = &t
	}
	
	// Get tasks
	tasks, err := s.app.GetTasks(ctx, queries.GetTasks{
		ManagerID: request.GetManagerId(),
		Filter:    filter,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	// Convert to proto
	protoTasks := make([]*schedulerspb.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = s.taskFromDomain(task)
	}
	
	// Get total count
	count, err := s.app.CountTasksByManagerID(ctx, queries.CountTasksByManagerID{
		ManagerID: request.GetManagerId(),
		Filter:    filter,
	})
	if err != nil {
		// Non-critical error, log but don't fail
		span.RecordError(err)
		count = len(tasks)
	}
	
	return &schedulerspb.GetScheduledTasksResponse{
		Tasks: protoTasks,
		Total: int32(count),
	}, nil
}

// UpdateTask updates a scheduled task
func (s server) UpdateTask(ctx context.Context, request *schedulerspb.UpdateTaskRequest) (*schedulerspb.UpdateTaskResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	span.SetAttributes(
		attribute.String("TaskID", request.GetTaskId()),
	)
	
	// Build update command based on what's provided
	cmd := commands.UpdateTask{
		TaskID: request.GetTaskId(),
	}
	
	// Update scheduled time if provided
	if request.GetScheduledAt() > 0 {
		t := time.Unix(request.GetScheduledAt(), 0)
		cmd.ScheduledAt = &t
	}
	
	// Update payload if provided
	if len(request.GetPayload()) > 0 {
		cmd.Payload = request.GetPayload()
	}
	
	// Execute update
	err := s.app.UpdateTask(ctx, cmd)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	// Handle status updates separately (these would be done through different commands)
	if request.GetStatus() != schedulerspb.TaskStatus_TASK_STATUS_UNSPECIFIED {
		switch request.GetStatus() {
		case schedulerspb.TaskStatus_RUNNING:
			err = s.app.ExecuteTask(ctx, commands.ExecuteTask{TaskID: request.GetTaskId()})
		case schedulerspb.TaskStatus_COMPLETED:
			err = s.app.CompleteTask(ctx, commands.CompleteTask{
				TaskID: request.GetTaskId(),
				Result: request.GetResult(),
			})
		case schedulerspb.TaskStatus_FAILED:
			err = s.app.FailTask(ctx, commands.FailTask{
				TaskID:       request.GetTaskId(),
				ErrorMessage: request.GetError(),
			})
		}
		
		if err != nil {
			span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}
	
	// Get updated task
	task, err := s.app.GetTask(ctx, queries.GetTask{TaskID: request.GetTaskId()})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	return &schedulerspb.UpdateTaskResponse{
		Task: s.taskFromDomain(task),
	}, nil
}

// Helper to convert domain task to proto
func (s server) taskFromDomain(task *domain.CatalogTask) *schedulerspb.Task {
	protoTask := &schedulerspb.Task{
		Id:          task.ID,
		ManagerId:   task.ManagerID,
		TaskType:    task.TaskType,
		ScheduledAt: task.ScheduledAt.Unix(),
		Payload:     task.Payload,
		Status:      s.stringToTaskStatus(task.Status),
		CreatedAt:   task.CreatedAt.Unix(),
		UpdatedAt:   task.UpdatedAt.Unix(),
		Result:      task.Result,
		Error:       task.ErrorMessage,
	}
	
	if task.ExecutedAt != nil {
		protoTask.ExecutedAt = task.ExecutedAt.Unix()
	}
	
	return protoTask
}

// Helper to convert string status to proto enum
func (s server) stringToTaskStatus(status string) schedulerspb.TaskStatus {
	switch status {
	case string(domain.TaskStatusPending):
		return schedulerspb.TaskStatus_PENDING
	case string(domain.TaskStatusRunning):
		return schedulerspb.TaskStatus_RUNNING
	case string(domain.TaskStatusCompleted):
		return schedulerspb.TaskStatus_COMPLETED
	case string(domain.TaskStatusFailed):
		return schedulerspb.TaskStatus_FAILED
	case string(domain.TaskStatusCancelled):
		return schedulerspb.TaskStatus_CANCELLED
	default:
		return schedulerspb.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

// Helper to convert proto enum to string status
func (s server) taskStatusToString(status schedulerspb.TaskStatus) string {
	switch status {
	case schedulerspb.TaskStatus_PENDING:
		return string(domain.TaskStatusPending)
	case schedulerspb.TaskStatus_RUNNING:
		return string(domain.TaskStatusRunning)
	case schedulerspb.TaskStatus_COMPLETED:
		return string(domain.TaskStatusCompleted)
	case schedulerspb.TaskStatus_FAILED:
		return string(domain.TaskStatusFailed)
	case schedulerspb.TaskStatus_CANCELLED:
		return string(domain.TaskStatusCancelled)
	default:
		return ""
	}
}
