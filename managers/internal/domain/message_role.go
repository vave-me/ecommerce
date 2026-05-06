package domain

type MessageRole string

const (
	UserRole    MessageRole = "user"
	ManagerRole MessageRole = "manager"
	SystemRole  MessageRole = "system"
)
