package grpc

import (
	"context"
	"middleman/assistants/assistantspb"
	"middleman/scheduler/internal/domain"
	"middleman/internal/rpc"

	"google.golang.org/grpc"
)

type AssistantRepository struct {
	endpoint    string
	assistantID string
}

var _ domain.AssistantRepository = (*AssistantRepository)(nil)

func NewAssistantRepository(endpoint string, assistantID string) AssistantRepository {
	return AssistantRepository{
		endpoint:    endpoint,
		assistantID: assistantID,
	}
}

func (r AssistantRepository) ProcessUserInput(ctx context.Context, task string, taskContext map[string]string) (string, error) {
	var conn *grpc.ClientConn
	conn, err := r.dial(ctx)
	if err != nil {
		return "", err
	}

	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	// Add scheduler context to the task context
	if taskContext == nil {
		taskContext = make(map[string]string)
	}
	taskContext["source"] = "scheduler"
	taskContext["request_type"] = "scheduled_task"

	resp, err := assistantspb.NewAssistantsServiceClient(conn).ProcessUserInput(ctx, &assistantspb.ProcessUserInputRequest{
		AssistantId: r.assistantID,
		Message:     task,
		Context:     taskContext,
		RequestType: "scheduled_task",
	})
	if err != nil {
		return "", err
	}

	// Return the assistant's response message
	return resp.GetMessage(), nil
}

func (r AssistantRepository) dial(ctx context.Context) (*grpc.ClientConn, error) {
	return rpc.Dial(ctx, r.endpoint)
}