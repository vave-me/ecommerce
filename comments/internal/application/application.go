package application

import (
	"context"
	"middleman/comments/internal/application/commands"
	"middleman/comments/internal/application/queries"
	"middleman/comments/internal/domain"
	"middleman/internal/ddd"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		AddComment(ctx context.Context, cmd commands.AddComment) error
		ApproveComment(ctx context.Context, cmd commands.ApproveComment) error
		EditComment(ctx context.Context, cmd commands.EditComment) error
		FlagComment(ctx context.Context, cmd commands.FlagComment) error
		RejectComment(ctx context.Context, cmd commands.RejectComment) error
		RemoveComment(ctx context.Context, cmd commands.RemoveComment) error
	}
	Queries interface {
		GetComment(ctx context.Context, query queries.GetComment) (*domain.MiddlemanComment, error)
		GetComments(ctx context.Context, query queries.GetComments) ([]*domain.MiddlemanComment, error)
		GetCommentsBySender(ctx context.Context, query queries.GetCommentsBySender) ([]*domain.MiddlemanComment, error)
		GetMostCommented(ctx context.Context, query queries.GetMostCommented) ([]*domain.ItemCommentCount, error)
		GetMostCommentedByCategory(ctx context.Context, query queries.GetMostCommentedByCategory) ([]*domain.ItemCommentCount, error)
		GetApprovedComments(ctx context.Context, query queries.GetApprovedComments) ([]*domain.MiddlemanComment, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddCommentHandler
		commands.ApproveCommentHandler
		commands.EditCommentHandler
		commands.FlagCommentHandler
		commands.RejectCommentHandler
		commands.RemoveCommentHandler
	}

	appQueries struct {
		queries.GetCommentHandler
		queries.GetCommentsHandler
		queries.GetCommentsBySenderHandler
		queries.GetMostCommentedHandler
		queries.GetMostCommentedByCategoryHandler
		queries.GetApprovedCommentsHandler
	}
)

var _ App = (*Application)(nil)

func New(comments domain.CommentRepository, middlemanComments domain.MiddlemanCacheRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddCommentHandler:     commands.NewAddCommentHandler(comments, publisher),
			ApproveCommentHandler: commands.NewApproveCommentHandler(comments, publisher),
			EditCommentHandler:    commands.NewEditCommentHandler(comments, publisher),
			FlagCommentHandler:    commands.NewFlagCommentHandler(comments, publisher),
			RejectCommentHandler:  commands.NewRejectCommentHandler(comments, publisher),
			RemoveCommentHandler:  commands.NewRemoveCommentHandler(comments, publisher),
		},
		appQueries: appQueries{
			GetCommentHandler:                 queries.NewGetCommentHandler(middlemanComments),
			GetCommentsHandler:                queries.NewGetCommentsHandler(middlemanComments),
			GetCommentsBySenderHandler:        queries.NewGetCommentsBySenderHandler(middlemanComments),
			GetMostCommentedHandler:           queries.NewGetMostCommentedHandler(middlemanComments),
			GetMostCommentedByCategoryHandler: queries.NewGetMostCommentedByCategoryHandler(middlemanComments),
			GetApprovedCommentsHandler:        queries.NewGetApprovedCommentsHandler(middlemanComments),
		},
	}
}
