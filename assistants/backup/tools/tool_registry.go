package tools

import (
	"context"
	"sync"

	"middleman/assistants/internal/domain"
	ai2 "middleman/internal/ai"
)

// ToolRegistry manages all tool definitions and executions
type ToolRegistry struct {
	// Repositories
	activityRepo     domain.ActivityRepository
	basketRepo       domain.BasketRepository
	categoryRepo     domain.CategoryRepository
	commentRepo      domain.CommentRepository
	geocodingRepo    domain.GeocodingRepository
	mailerRepo       domain.MailerRepository
	mediaRepo        domain.MiddlemanMediaRepository
	messageRepo      domain.MessagesRepository
	metricRepo       domain.MetricRepository
	newsletterRepo   domain.NewsletterRepository
	notificationRepo domain.NotificationRepository
	offerRepo        domain.OfferRepository
	orderRepo        domain.OrderRepository
	paymentRepo      domain.PaymentRepository
	postRepo         domain.PostRepository
	productRepo      domain.ProductRepository
	reviewRepo       domain.ReviewRepository
	serviceRepo      domain.ServiceRepository
	shippingRepo     domain.ShippingRepository
	supportRepo      domain.SupportRepository
	userRepo         domain.UserRepository
	variantRepo      domain.VariantRepository
	vectorRepo       domain.VectorRepository
	wishlistRepo     domain.WishlistRepository
	followingRepo    domain.FollowingRepository

	// Tool definitions cache
	toolDefs []ai2.Tool
	mu       sync.RWMutex
}

// NewToolRegistry creates a new tool registry with all repositories
func NewToolRegistry(
	activityRepo domain.ActivityRepository,
	basketRepo domain.BasketRepository,
	categoryRepo domain.CategoryRepository,
	commentRepo domain.CommentRepository,
	geocodingRepo domain.GeocodingRepository,
	mailerRepo domain.MailerRepository,
	mediaRepo domain.MiddlemanMediaRepository,
	messageRepo domain.MessagesRepository,
	metricRepo domain.MetricRepository,
	newsletterRepo domain.NewsletterRepository,
	notificationRepo domain.NotificationRepository,
	offerRepo domain.OfferRepository,
	orderRepo domain.OrderRepository,
	paymentRepo domain.PaymentRepository,
	postRepo domain.PostRepository,
	productRepo domain.ProductRepository,
	reviewRepo domain.ReviewRepository,
	serviceRepo domain.ServiceRepository,
	shippingRepo domain.ShippingRepository,
	supportRepo domain.SupportRepository,
	userRepo domain.UserRepository,
	variantRepo domain.VariantRepository,
	vectorRepo domain.VectorRepository,
	wishlistRepo domain.WishlistRepository,
	followingRepo domain.FollowingRepository,
) *ToolRegistry {
	registry := &ToolRegistry{
		activityRepo:     activityRepo,
		basketRepo:       basketRepo,
		categoryRepo:     categoryRepo,
		commentRepo:      commentRepo,
		geocodingRepo:    geocodingRepo,
		mailerRepo:       mailerRepo,
		mediaRepo:        mediaRepo,
		messageRepo:      messageRepo,
		metricRepo:       metricRepo,
		newsletterRepo:   newsletterRepo,
		notificationRepo: notificationRepo,
		offerRepo:        offerRepo,
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		postRepo:         postRepo,
		productRepo:      productRepo,
		reviewRepo:       reviewRepo,
		serviceRepo:      serviceRepo,
		shippingRepo:     shippingRepo,
		supportRepo:      supportRepo,
		userRepo:         userRepo,
		variantRepo:      variantRepo,
		vectorRepo:       vectorRepo,
		wishlistRepo:     wishlistRepo,
		followingRepo:    followingRepo,
	}

	// Initialize tool definitions with OpenAI-compliant format
	registry.toolDefs = CreateOpenAICompliantTools()

	return registry
}

// GetToolDefinitions returns all available tool definitions
func (r *ToolRegistry) GetToolDefinitions() []ai2.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// OpenAI has a limit of 128 tools per request
	// We've already reduced tools to essential ones
	if len(r.toolDefs) > 128 {
		// Return first 128 tools to comply with OpenAI limit
		return r.toolDefs[:128]
	}
	
	return r.toolDefs
}

// ExecuteTool executes a tool by name with the given parameters
func (r *ToolRegistry) ExecuteTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	// Use comprehensive registry for all tool executions
	comprehensiveRegistry := NewComprehensiveToolRegistry(r)
	return comprehensiveRegistry.MainExecute(ctx, toolName, params)
}
