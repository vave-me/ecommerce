package application

import (
	"context"
	"middleman/internal/ddd"
	"middleman/reviews/internal/application/commands"
	"middleman/reviews/internal/application/queries"
	"middleman/reviews/internal/domain"
)

type (
	App interface {
		Commands
		Queries
	}
	Commands interface {
		AddReview(ctx context.Context, cmd commands.AddReview) error
		ApproveReview(ctx context.Context, cmd commands.ApproveReview) error
		EditReview(ctx context.Context, cmd commands.EditReview) error
		FlagReview(ctx context.Context, cmd commands.FlagReview) error
		RejectReview(ctx context.Context, cmd commands.RejectReview) error
		RemoveReview(ctx context.Context, cmd commands.RemoveReview) error
	}
	Queries interface {
		GetReview(ctx context.Context, query queries.GetReview) (*domain.MiddlemanReview, error)
		GetReviews(ctx context.Context, query queries.GetReviews) ([]*domain.MiddlemanReview, error)
		GetReviewsBySender(ctx context.Context, query queries.GetReviewsBySender) ([]*domain.MiddlemanReview, error)
		GetMostReviewed(ctx context.Context, query queries.GetMostReviewed) ([]*domain.ItemReviewCount, error)
		GetMostReviewedByCategory(ctx context.Context, query queries.GetMostReviewedByCategory) ([]*domain.ItemReviewCount, error)
		GetApprovedReviews(ctx context.Context, query queries.GetApprovedReviews) ([]*domain.MiddlemanReview, error)
	}

	Application struct {
		appCommands
		appQueries
	}

	appCommands struct {
		commands.AddReviewHandler
		commands.ApproveReviewHandler
		commands.EditReviewHandler
		commands.FlagReviewHandler
		commands.RejectReviewHandler
		commands.RemoveReviewHandler
	}

	appQueries struct {
		queries.GetReviewHandler
		queries.GetReviewsHandler
		queries.GetReviewsBySenderHandler
		queries.GetMostReviewedHandler
		queries.GetMostReviewedByCategoryHandler
		queries.GetApprovedReviewsHandler
	}
)

var _ App = (*Application)(nil)

func New(reviews domain.ReviewRepository, middlemanReviews domain.MiddlemanCacheRepository,
	publisher ddd.EventPublisher[ddd.Event]) *Application {
	return &Application{
		appCommands: appCommands{
			AddReviewHandler:     commands.NewAddReviewHandler(reviews, publisher),
			ApproveReviewHandler: commands.NewApproveReviewHandler(reviews, publisher),
			EditReviewHandler:    commands.NewEditReviewHandler(reviews, publisher),
			FlagReviewHandler:    commands.NewFlagReviewHandler(reviews, publisher),
			RejectReviewHandler:  commands.NewRejectReviewHandler(reviews, publisher),
			RemoveReviewHandler:  commands.NewRemoveReviewHandler(reviews, publisher),
		},
		appQueries: appQueries{
			GetReviewHandler:                 queries.NewGetReviewHandler(middlemanReviews),
			GetReviewsHandler:                queries.NewGetReviewsHandler(middlemanReviews),
			GetReviewsBySenderHandler:        queries.NewGetReviewsBySenderHandler(middlemanReviews),
			GetMostReviewedHandler:           queries.NewGetMostReviewedHandler(middlemanReviews),
			GetMostReviewedByCategoryHandler: queries.NewGetMostReviewedByCategoryHandler(middlemanReviews),
			GetApprovedReviewsHandler:        queries.NewGetApprovedReviewsHandler(middlemanReviews),
		},
	}
}
