package domain

type MessageRole string

const (
	UserRole      MessageRole = "user"
	AssistantRole MessageRole = "assistant"
	SystemRole    MessageRole = "system"
)
