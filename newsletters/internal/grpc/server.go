package grpc

import (
	"context"
	"time"
	
	"middleman/internal/auth"
	"middleman/internal/errorsotel"
	"middleman/newsletters/internal/application"
	"middleman/newsletters/internal/application/commands"
	"middleman/newsletters/internal/application/queries"
	"middleman/newsletters/internal/domain"
	"middleman/newsletters/newsletterspb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_code "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	app application.App
	newsletterspb.UnimplementedNewslettersServiceServer
}

var _ newsletterspb.NewslettersServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	newsletterspb.RegisterNewslettersServiceServer(registrar, server{app: app})
	return nil
}

// Newsletter Management
func (s server) CreateNewsletter(ctx context.Context, request *newsletterspb.CreateNewsletterRequest) (*newsletterspb.CreateNewsletterResponse, error) {
	span := trace.SpanFromContext(ctx)
	
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	
	newsletterID, err := s.app.CreateNewsletter(ctx, commands.CreateNewsletter{
		UserID:      claims.Subject,
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Frequency:   request.GetFrequency(),
		Category:    request.GetCategory(),
		TemplateID:  request.GetTemplateId(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("newsletterID", newsletterID))

	// Get the created newsletter
	newsletter, err := s.app.GetNewsletter(ctx, queries.GetNewsletter{ID: newsletterID})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.CreateNewsletterResponse{
		Newsletter: s.newsletterFromDomain(newsletter),
	}, nil
}

func (s server) UpdateNewsletter(ctx context.Context, request *newsletterspb.UpdateNewsletterRequest) (*newsletterspb.UpdateNewsletterResponse, error) {
	span := trace.SpanFromContext(ctx)

	err := s.app.UpdateNewsletter(ctx, commands.UpdateNewsletter{
		ID:          request.GetId(),
		Name:        request.GetName(),
		Description: request.GetDescription(),
		Frequency:   request.GetFrequency(),
		Category:    request.GetCategory(),
		TemplateID:  request.GetTemplateId(),
		IsActive:    request.GetIsActive(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Get the updated newsletter
	newsletter, err := s.app.GetNewsletter(ctx, queries.GetNewsletter{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.UpdateNewsletterResponse{
		Newsletter: s.newsletterFromDomain(newsletter),
	}, nil
}

func (s server) GetNewsletter(ctx context.Context, request *newsletterspb.GetNewsletterRequest) (*newsletterspb.GetNewsletterResponse, error) {
	newsletter, err := s.app.GetNewsletter(ctx, queries.GetNewsletter{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetNewsletterResponse{
		Newsletter: s.newsletterFromDomain(newsletter),
	}, nil
}

func (s server) ListNewsletters(ctx context.Context, request *newsletterspb.ListNewslettersRequest) (*newsletterspb.ListNewslettersResponse, error) {
	newsletters, total, err := s.app.ListNewsletters(ctx, queries.ListNewsletters{
		UserID:     request.GetUserId(),
		Category:   request.GetCategory(),
		ActiveOnly: request.GetActiveOnly(),
		Page:       int(request.GetPage()),
		Limit:      int(request.GetLimit()),
	})
	if err != nil {
		return nil, err
	}

	pbNewsletters := make([]*newsletterspb.Newsletter, len(newsletters))
	for i, n := range newsletters {
		pbNewsletters[i] = s.newsletterFromDomain(n)
	}

	return &newsletterspb.ListNewslettersResponse{
		Newsletters: pbNewsletters,
		Total:       int32(total),
		Page:        request.GetPage(),
		Limit:       request.GetLimit(),
	}, nil
}

func (s server) DeleteNewsletter(ctx context.Context, request *newsletterspb.DeleteNewsletterRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteNewsletter(ctx, commands.DeleteNewsletter{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Subscription Management
func (s server) Subscribe(ctx context.Context, request *newsletterspb.SubscribeRequest) (*newsletterspb.SubscribeResponse, error) {
	span := trace.SpanFromContext(ctx)

	var freqOverride *domain.NewsletterFrequency
	if request.GetPreferences() != nil && request.GetPreferences().GetFrequencyOverride() != "" {
		freq := domain.ToNewsletterFrequency(request.GetPreferences().GetFrequencyOverride())
		freqOverride = &freq
	}

	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, status.Error(grpc_code.Unauthenticated, "authentication required")
	}
	
	subscriptionID, err := s.app.Subscribe(ctx, commands.Subscribe{
		UserID:       claims.Subject,
		NewsletterID: request.GetNewsletterId(),
		Preferences: domain.SubscriptionPreferences{
			FrequencyOverride: freqOverride,
			Topics:            request.GetPreferences().GetTopics(),
			Format:            domain.ToContentFormat(request.GetPreferences().GetFormat()),
		},
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("subscriptionID", subscriptionID))

	// Get the created subscription
	subscription, err := s.app.GetSubscription(ctx, queries.GetSubscription{ID: subscriptionID})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.SubscribeResponse{
		Subscription: s.subscriptionFromDomain(subscription),
	}, nil
}

func (s server) Unsubscribe(ctx context.Context, request *newsletterspb.UnsubscribeRequest) (*emptypb.Empty, error) {
	err := s.app.Unsubscribe(ctx, commands.Unsubscribe{
		SubscriptionID: request.GetSubscriptionId(),
		Reason:         request.GetReason(),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s server) UpdateSubscription(ctx context.Context, request *newsletterspb.UpdateSubscriptionRequest) (*newsletterspb.UpdateSubscriptionResponse, error) {
	var freqOverride *domain.NewsletterFrequency
	if request.GetPreferences() != nil && request.GetPreferences().GetFrequencyOverride() != "" {
		freq := domain.ToNewsletterFrequency(request.GetPreferences().GetFrequencyOverride())
		freqOverride = &freq
	}

	err := s.app.UpdateSubscription(ctx, commands.UpdateSubscription{
		ID:     request.GetId(),
		Status: request.GetStatus(),
		Preferences: domain.SubscriptionPreferences{
			FrequencyOverride: freqOverride,
			Topics:            request.GetPreferences().GetTopics(),
			Format:            domain.ToContentFormat(request.GetPreferences().GetFormat()),
		},
	})
	if err != nil {
		return nil, err
	}

	// Get the updated subscription
	subscription, err := s.app.GetSubscription(ctx, queries.GetSubscription{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.UpdateSubscriptionResponse{
		Subscription: s.subscriptionFromDomain(subscription),
	}, nil
}

func (s server) GetSubscription(ctx context.Context, request *newsletterspb.GetSubscriptionRequest) (*newsletterspb.GetSubscriptionResponse, error) {
	subscription, err := s.app.GetSubscription(ctx, queries.GetSubscription{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetSubscriptionResponse{
		Subscription: s.subscriptionFromDomain(subscription),
	}, nil
}

func (s server) ListSubscriptions(ctx context.Context, request *newsletterspb.ListSubscriptionsRequest) (*newsletterspb.ListSubscriptionsResponse, error) {
	subscriptions, total, err := s.app.ListSubscriptions(ctx, queries.ListSubscriptions{
		UserID:       request.GetUserId(),
		NewsletterID: request.GetNewsletterId(),
		Status:       request.GetStatus(),
		Page:         int(request.GetPage()),
		Limit:        int(request.GetLimit()),
	})
	if err != nil {
		return nil, err
	}

	pbSubscriptions := make([]*newsletterspb.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		pbSubscriptions[i] = s.subscriptionFromDomain(sub)
	}

	return &newsletterspb.ListSubscriptionsResponse{
		Subscriptions: pbSubscriptions,
		Total:         int32(total),
		Page:          request.GetPage(),
		Limit:         request.GetLimit(),
	}, nil
}

// Edition Management
func (s server) CreateEdition(ctx context.Context, request *newsletterspb.CreateEditionRequest) (*newsletterspb.CreateEditionResponse, error) {
	span := trace.SpanFromContext(ctx)

	var scheduledAt *time.Time
	if request.GetScheduledAt() != nil {
		t := request.GetScheduledAt().AsTime()
		scheduledAt = &t
	}

	editionID, err := s.app.CreateEdition(ctx, commands.CreateEdition{
		NewsletterID: request.GetNewsletterId(),
		Subject:      request.GetSubject(),
		ContentHTML:  request.GetContentHtml(),
		ContentText:  request.GetContentText(),
		TemplateData: request.GetTemplateData(),
		ScheduledAt:  scheduledAt,
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("editionID", editionID))

	// Get the created edition
	edition, err := s.app.GetEdition(ctx, queries.GetEdition{ID: editionID})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.CreateEditionResponse{
		Edition: s.editionFromDomain(edition),
	}, nil
}

func (s server) UpdateEdition(ctx context.Context, request *newsletterspb.UpdateEditionRequest) (*newsletterspb.UpdateEditionResponse, error) {
	var scheduledAt *time.Time
	if request.GetScheduledAt() != nil {
		t := request.GetScheduledAt().AsTime()
		scheduledAt = &t
	}

	err := s.app.UpdateEdition(ctx, commands.UpdateEdition{
		ID:           request.GetId(),
		Subject:      request.GetSubject(),
		ContentHTML:  request.GetContentHtml(),
		ContentText:  request.GetContentText(),
		TemplateData: request.GetTemplateData(),
		ScheduledAt:  scheduledAt,
	})
	if err != nil {
		return nil, err
	}

	// Get the updated edition
	edition, err := s.app.GetEdition(ctx, queries.GetEdition{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.UpdateEditionResponse{
		Edition: s.editionFromDomain(edition),
	}, nil
}

func (s server) GetEdition(ctx context.Context, request *newsletterspb.GetEditionRequest) (*newsletterspb.GetEditionResponse, error) {
	edition, err := s.app.GetEdition(ctx, queries.GetEdition{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetEditionResponse{
		Edition: s.editionFromDomain(edition),
	}, nil
}

func (s server) ListEditions(ctx context.Context, request *newsletterspb.ListEditionsRequest) (*newsletterspb.ListEditionsResponse, error) {
	editions, total, err := s.app.ListEditions(ctx, queries.ListEditions{
		NewsletterID: request.GetNewsletterId(),
		Status:       request.GetStatus(),
		Page:         int(request.GetPage()),
		Limit:        int(request.GetLimit()),
	})
	if err != nil {
		return nil, err
	}

	pbEditions := make([]*newsletterspb.NewsletterEdition, len(editions))
	for i, e := range editions {
		pbEditions[i] = s.editionFromDomain(e)
	}

	return &newsletterspb.ListEditionsResponse{
		Editions: pbEditions,
		Total:    int32(total),
		Page:     request.GetPage(),
		Limit:    request.GetLimit(),
	}, nil
}

func (s server) ScheduleEdition(ctx context.Context, request *newsletterspb.ScheduleEditionRequest) (*newsletterspb.ScheduleEditionResponse, error) {
	err := s.app.ScheduleEdition(ctx, commands.ScheduleEdition{
		ID:          request.GetId(),
		ScheduledAt: request.GetScheduledAt().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	// Get the updated edition
	edition, err := s.app.GetEdition(ctx, queries.GetEdition{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.ScheduleEditionResponse{
		Edition: s.editionFromDomain(edition),
	}, nil
}

func (s server) SendEdition(ctx context.Context, request *newsletterspb.SendEditionRequest) (*newsletterspb.SendEditionResponse, error) {
	recipientCount, err := s.app.SendEdition(ctx, commands.SendEdition{
		ID:       request.GetId(),
		TestMode: request.GetTestMode(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.SendEditionResponse{
		Id:               request.GetId(),
		RecipientsQueued: int32(recipientCount),
	}, nil
}

func (s server) PreviewEdition(ctx context.Context, request *newsletterspb.PreviewEditionRequest) (*newsletterspb.PreviewEditionResponse, error) {
	// Get the edition
	edition, err := s.app.GetEdition(ctx, queries.GetEdition{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	// TODO: Apply template rendering with template data

	return &newsletterspb.PreviewEditionResponse{
		PreviewHtml: edition.ContentHTML,
		PreviewText: edition.ContentText,
	}, nil
}

// Template Management
func (s server) CreateTemplate(ctx context.Context, request *newsletterspb.CreateTemplateRequest) (*newsletterspb.CreateTemplateResponse, error) {
	span := trace.SpanFromContext(ctx)

	templateID, err := s.app.CreateTemplate(ctx, commands.CreateTemplate{
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		HTMLTemplate: request.GetHtmlTemplate(),
		TextTemplate: request.GetTextTemplate(),
		Variables:    request.GetVariables(),
		PreviewData:  request.GetPreviewData(),
		IsPublic:     request.GetIsPublic(),
	})
	if err != nil {
		span.RecordError(err, trace.WithAttributes(errorsotel.ErrAttrs(err)...))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("templateID", templateID))

	// Get the created template
	template, err := s.app.GetTemplate(ctx, queries.GetTemplate{ID: templateID})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.CreateTemplateResponse{
		Template: s.templateFromDomain(template),
	}, nil
}

func (s server) UpdateTemplate(ctx context.Context, request *newsletterspb.UpdateTemplateRequest) (*newsletterspb.UpdateTemplateResponse, error) {
	err := s.app.UpdateTemplate(ctx, commands.UpdateTemplate{
		ID:           request.GetId(),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		HTMLTemplate: request.GetHtmlTemplate(),
		TextTemplate: request.GetTextTemplate(),
		Variables:    request.GetVariables(),
		PreviewData:  request.GetPreviewData(),
		IsPublic:     request.GetIsPublic(),
	})
	if err != nil {
		return nil, err
	}

	// Get the updated template
	template, err := s.app.GetTemplate(ctx, queries.GetTemplate{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.UpdateTemplateResponse{
		Template: s.templateFromDomain(template),
	}, nil
}

func (s server) GetTemplate(ctx context.Context, request *newsletterspb.GetTemplateRequest) (*newsletterspb.GetTemplateResponse, error) {
	template, err := s.app.GetTemplate(ctx, queries.GetTemplate{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetTemplateResponse{
		Template: s.templateFromDomain(template),
	}, nil
}

func (s server) ListTemplates(ctx context.Context, request *newsletterspb.ListTemplatesRequest) (*newsletterspb.ListTemplatesResponse, error) {
	templates, total, err := s.app.ListTemplates(ctx, queries.ListTemplates{
		UserID:     request.GetUserId(),
		PublicOnly: request.GetPublicOnly(),
		Page:       int(request.GetPage()),
		Limit:      int(request.GetLimit()),
	})
	if err != nil {
		return nil, err
	}

	pbTemplates := make([]*newsletterspb.NewsletterTemplate, len(templates))
	for i, t := range templates {
		pbTemplates[i] = s.templateFromDomain(t)
	}

	return &newsletterspb.ListTemplatesResponse{
		Templates: pbTemplates,
		Total:     int32(total),
		Page:      request.GetPage(),
		Limit:     request.GetLimit(),
	}, nil
}

func (s server) DeleteTemplate(ctx context.Context, request *newsletterspb.DeleteTemplateRequest) (*emptypb.Empty, error) {
	err := s.app.DeleteTemplate(ctx, commands.DeleteTemplate{
		ID: request.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Analytics
func (s server) GetNewsletterStats(ctx context.Context, request *newsletterspb.GetNewsletterStatsRequest) (*newsletterspb.GetNewsletterStatsResponse, error) {
	stats, err := s.app.GetNewsletterStats(ctx, queries.GetNewsletterStats{
		NewsletterID: request.GetNewsletterId(),
		StartDate:    request.GetStartDate().AsTime(),
		EndDate:      request.GetEndDate().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetNewsletterStatsResponse{
		TotalSubscribers:  int32(stats.TotalSubscribers),
		NewSubscribers:    int32(stats.NewSubscribers),
		Unsubscribes:      int32(stats.Unsubscribes),
		EditionsSent:      int32(stats.EditionsSent),
		AverageOpenRate:   stats.AverageOpenRate,
		AverageClickRate:  stats.AverageClickRate,
	}, nil
}

func (s server) GetEditionStats(ctx context.Context, request *newsletterspb.GetEditionStatsRequest) (*newsletterspb.GetEditionStatsResponse, error) {
	stats, err := s.app.GetEditionStats(ctx, queries.GetEditionStats{
		EditionID: request.GetEditionId(),
	})
	if err != nil {
		return nil, err
	}

	return &newsletterspb.GetEditionStatsResponse{
		Recipients: int32(stats.Recipients),
		Delivered:  int32(stats.Delivered),
		Opened:     int32(stats.Opened),
		Clicked:    int32(stats.Clicked),
		Bounced:    int32(stats.Bounced),
		Complaints: int32(stats.Complaints),
		OpenRate:   stats.OpenRate,
		ClickRate:  stats.ClickRate,
	}, nil
}

// Helper functions to convert domain objects to protobuf
func (s server) newsletterFromDomain(n *domain.CatalogNewsletter) *newsletterspb.Newsletter {
	return &newsletterspb.Newsletter{
		Id:              n.ID,
		UserId:          n.UserID,
		Name:            n.Name,
		Description:     n.Description,
		Frequency:       n.Frequency,
		Category:        n.Category,
		TemplateId:      n.TemplateID,
		IsActive:        n.IsActive,
		CreatedAt:       timestamppb.New(n.CreatedAt),
		UpdatedAt:       timestamppb.New(n.UpdatedAt),
		SubscriberCount: int32(n.SubscriberCount),
	}
}

func (s server) subscriptionFromDomain(sub *domain.CatalogSubscription) *newsletterspb.Subscription {
	resp := &newsletterspb.Subscription{
		Id:            sub.ID,
		UserId:        sub.UserID,
		NewsletterId:  sub.NewsletterID,
		Status:        sub.Status,
		SubscribedAt:  timestamppb.New(sub.SubscribedAt),
		Preferences: &newsletterspb.SubscriptionPreferences{
			FrequencyOverride: sub.FrequencyOverride,
			Topics:            sub.Topics,
			Format:            sub.Format,
		},
	}

	if sub.UnsubscribedAt != nil {
		resp.UnsubscribedAt = timestamppb.New(*sub.UnsubscribedAt)
	}

	if sub.Newsletter != nil {
		resp.Newsletter = s.newsletterFromDomain(sub.Newsletter)
	}

	return resp
}

func (s server) editionFromDomain(e *domain.CatalogEdition) *newsletterspb.NewsletterEdition {
	resp := &newsletterspb.NewsletterEdition{
		Id:             e.ID,
		NewsletterId:   e.NewsletterID,
		Subject:        e.Subject,
		ContentHtml:    e.ContentHTML,
		ContentText:    e.ContentText,
		TemplateData:   e.TemplateData,
		Status:         e.Status,
		CreatedBy:      e.CreatedBy,
		CreatedAt:      timestamppb.New(e.CreatedAt),
		UpdatedAt:      timestamppb.New(e.UpdatedAt),
		RecipientCount: int32(e.RecipientCount),
	}

	if e.ScheduledAt != nil {
		resp.ScheduledAt = timestamppb.New(*e.ScheduledAt)
	}

	if e.SentAt != nil {
		resp.SentAt = timestamppb.New(*e.SentAt)
	}

	return resp
}

func (s server) templateFromDomain(t *domain.CatalogTemplate) *newsletterspb.NewsletterTemplate {
	return &newsletterspb.NewsletterTemplate{
		Id:           t.ID,
		UserId:       t.UserID,
		Name:         t.Name,
		Description:  t.Description,
		HtmlTemplate: t.HTMLTemplate,
		TextTemplate: t.TextTemplate,
		Variables:    t.Variables,
		PreviewData:  t.PreviewData,
		IsPublic:     t.IsPublic,
		CreatedAt:    timestamppb.New(t.CreatedAt),
		UpdatedAt:    timestamppb.New(t.UpdatedAt),
	}
}