package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/messages/internal/application/commands"
	"middleman/messages/internal/application/queries"
	"middleman/messages/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		StartConversation(ctx context.Context, cmd commands.StartConversation) error
		DeleteConversation(ctx context.Context, cmd commands.DeleteConversation) error
		SendMessage(ctx context.Context, cmd commands.SendMessage) error
		DeleteMessage(ctx context.Context, cmd commands.DeleteMessage) error
	}
	Queries interface {
		GetConversation(ctx context.Context, query queries.GetConversation) (*domain.MiddlemanConversation, error)
		GetConversationByRecipientAndItem(ctx context.Context, query queries.GetConversationByRecipientAndItem) (*domain.MiddlemanConversation, error)
		GetConversations(ctx context.Context, query queries.GetConversations) ([]*domain.MiddlemanConversation, error)
		GetMessage(ctx context.Context, query queries.GetMessage) (*domain.MiddlemanMessage, error)
		GetMessages(ctx context.Context, query queries.GetMessages) ([]*domain.MiddlemanMessage, error)
	}

	Application struct {
		appCommands
		appQueries
	}
	appCommands struct {
		commands.StartConversationHandler
		commands.DeleteConversationHandler
		commands.SendMessageHandler
		commands.DeleteMessageHandler
	}
	appQueries struct {
		queries.GetConversationHandler
		queries.GetConversationByRecipientAndItemHandler
		queries.GetConversationsHandler
		queries.GetMessagesHandler
		queries.GetMessageHandler
	}
)

var _ App = (*Application)(nil)

func New(conversations domain.ConversationRepository, messages domain.MessageRepository, middlemanConversations domain.MiddlemanRepository, messenger domain.MessengerRepository, publisher ddd.EventPublisher[ddd.Event]) *Application {

	return &Application{
		appCommands: appCommands{
			StartConversationHandler:  commands.NewStartConversationHandler(conversations, publisher),
			DeleteConversationHandler: commands.NewDeleteConversationHandler(conversations, publisher),
			SendMessageHandler:        commands.NewSendMessageHandler(messages, publisher),
			DeleteMessageHandler:      commands.NewDeleteMessageHandler(messages, publisher),
		},
		appQueries: appQueries{
			GetConversationHandler:                   queries.NewGetConversationHandler(middlemanConversations),
			GetConversationsHandler:                  queries.NewGetConversationsHandler(middlemanConversations),
			GetConversationByRecipientAndItemHandler: queries.NewGetConversationByRecipientAndItemHandler(middlemanConversations),
			GetMessageHandler:                        queries.NewGetMessageHandler(messenger),
			GetMessagesHandler:                       queries.NewGetMessagesHandler(messenger),
		},
	}
}
