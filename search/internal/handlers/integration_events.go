package handlers

import (
	"context"
	"fmt"
	"log"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/services/servicespb"

	"middleman/ordering/orderingpb"
	"middleman/posts/postspb"
	"middleman/products/productspb"

	"middleman/search/internal/application"
	"middleman/search/internal/models"

	"middleman/users/userspb"

	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	orders   application.OrderRepository
	products application.ProductCacheRepository
	users    application.UserCacheRepository
	posts    application.PostCacheRepository
	services application.ServiceCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, orders application.OrderRepository,
	users application.UserCacheRepository, products application.ProductCacheRepository, posts application.PostCacheRepository,
	services application.ServiceCacheRepository,
	mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		orders:   orders,
		users:    users,
		products: products,
		posts:    posts,
		services: services,
	}, zerolog.Logger{}, mws...)
}
func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {

	if _, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
		orderingpb.OrderReadiedEvent,
		orderingpb.OrderCanceledEvent,
		orderingpb.OrderCompletedEvent,
	}, am.GroupName("notification-orders")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductUpdatedEvent,
		productspb.ProductRemovedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductLeasedEvent,
		productspb.ProductSoldEvent,
		productspb.ProductPawnedEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductNegotiableToggledEvent,
		productspb.ProductArchivedEvent,
		productspb.ProductThumbnailAddedEvent,
		productspb.ProductThumbnailUpdatedEvent,
	}, am.GroupName("search-products")); err != nil {
		return
	}

	if _, err = subscriber.Subscribe(postspb.PostAggregateChannel, handlers, am.MessageFilter{
		postspb.PostAddedEvent,
		postspb.PostRemovedEvent,
		postspb.PostUpdatedEvent,
		postspb.PostArchivedEvent,
		postspb.PostThumbnailAddedEvent,
		postspb.PostThumbnailUpdatedEvent,
	}, am.GroupName("search-posts")); err != nil {
		return
	}
	if _, err = subscriber.Subscribe(servicespb.ServiceAggregateChannel, handlers, am.MessageFilter{
		servicespb.ServiceAddedEvent,
		servicespb.ServiceUpdatedEvent,
		servicespb.ServiceRemovedEvent,
	}, am.GroupName("search-services")); err != nil {
		return
	}
	if _, err = subscriber.Subscribe(userspb.UserAggregateChannel, handlers, am.MessageFilter{
		userspb.UserCreatedEvent,
		userspb.UserRenamedEvent,
	}, am.GroupName("search-users")); err != nil {
		return
	}

	return
}
func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case userspb.UserCreatedEvent:
		return h.onUserCreated(ctx, event)
	case userspb.UserRenamedEvent:
		return h.onUserRenamed(ctx, event)
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductUpdatedEvent:
		return h.onProductUpdated(ctx, event)
	case productspb.ProductRebrandedEvent:
		return h.onProductRebranded(ctx, event)
	case productspb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	case productspb.ProductThumbnailAddedEvent:
		return h.onProductThumbnailAdded(ctx, event)
	case productspb.ProductThumbnailUpdatedEvent:
		return h.onProductThumbnailUpdated(ctx, event)
	case productspb.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	case productspb.ProductPriceDecreasedEvent:
		return h.onProductPriceDecreased(ctx, event)
	case productspb.ProductLeasedEvent:
		return h.onProductLeased(ctx, event)
	case productspb.ProductSoldEvent:
		return h.onProductSold(ctx, event)
	case productspb.ProductPawnedEvent:
		return h.onProductPawned(ctx, event)
	case productspb.ProductStockAdjustedEvent:
		return h.onProductStockAdjusted(ctx, event)
	case productspb.ProductNegotiableToggledEvent:
		return h.onProductNegotiableToggled(ctx, event)
	case productspb.ProductArchivedEvent:
		return h.onProductArchived(ctx, event)
	case postspb.PostAddedEvent:
		return h.onPostAdded(ctx, event)
	case postspb.PostRemovedEvent:
		return h.onPostRemoved(ctx, event)
	case postspb.PostThumbnailAddedEvent:
		return h.onPostThumbnailAdded(ctx, event)
	case postspb.PostThumbnailUpdatedEvent:
		return h.onPostThumbnailUpdated(ctx, event)
	case servicespb.ServiceAddedEvent:
		return h.onServiceAdded(ctx, event)
	case servicespb.ServiceUpdatedEvent:
		return h.onServiceUpdated(ctx, event)
	case servicespb.ServiceRemovedEvent:
		return h.onServiceRemoved(ctx, event)
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case orderingpb.OrderReadiedEvent:
		return h.onOrderReadied(ctx, event)
	case orderingpb.OrderCanceledEvent:
		return h.onOrderCanceled(ctx, event)
	case orderingpb.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	case postspb.PostUpdatedEvent:
		return h.onPostUpdated(ctx, event)

	}

	return nil
}

func (h integrationHandlers[T]) onUserCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*userspb.UserCreated)
	return h.users.Add(ctx, payload.GetId(), payload.GetEmail(), payload.GetUserName(), payload.GetFirstName(), payload.GetLastName(), payload.GetLocation(), payload.GetEnabled())
}
func (h integrationHandlers[T]) onUserRenamed(ctx context.Context, event T) error {
	payload := event.Payload().(*userspb.UserRenamed)
	return h.users.Rename(ctx, payload.GetId(), payload.GetName())

}
func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductAdded)

	// Convert proto options to domain options
	domainOpts := make([]models.Option, len(payload.GetOptions()))
	for i, opt := range payload.GetOptions() {
		domainOpts[i] = models.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()),
		}
	}

	domainAttrs := make([]models.Attribute, len(payload.GetAttributes()))
	for i, opt := range payload.GetAttributes() {
		domainAttrs[i] = models.Attribute{
			Key:   opt.GetKey(),
			Value: opt.GetValue(),
		}
	}

	// Convert proto tags from repeated string => []string for domain
	domainTags := payload.GetTags()

	return h.products.Add(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetBasePrice(),
		payload.GetUserSellerId(),
		payload.GetCategoryId(),
		payload.CategorySlug,
		payload.GetBrand(),
		payload.GetCondition(),
		payload.GetModel(),
		domainTags,
		payload.GetManageStocks(),
		payload.GetStock(),
		payload.GetSku(),
		domainAttrs,
		payload.GetWeight(),
		payload.GetHeight(),
		payload.GetWidth(),
		payload.GetDepth(),
		payload.GetStatus(),
		payload.GetNegotiable(),
		payload.GetUserType(),
		payload.GetMiddlemanService(),
		payload.GetShippingCost(),
		payload.GetHasVariants(),
		domainOpts,
		float64(payload.Lat),
		float64(payload.Lng),
		payload.GetThumbnail(),
		models.ProductType, // now passing a slice of Option
	)
}

func (h integrationHandlers[T]) onProductRebranded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRebranded)

	// SAFETY CHECK: Fetch existing product first to preserve non-empty fields
	existing, err := h.products.Find(ctx, payload.GetId())
	if err != nil {
		log.Printf("❌ [onProductRebranded] Could not find existing product %s for rebrand event - skipping: %v", payload.GetId(), err)
		// CRITICAL: Cannot safely rebrand non-existent product with partial data - skip update
		return nil
	}

	// Merge update with existing data - only update non-empty fields
	name := existing.Name
	if payload.GetName() != "" {
		name = payload.GetName()
	}

	description := existing.Description
	if payload.GetDescription() != "" {
		description = payload.GetDescription()
	}

	basePrice := existing.BasePrice
	if payload.GetBasePrice() > 0 {
		basePrice = payload.GetBasePrice()
	}

	stock := existing.Stock
	if payload.GetStock() > 0 {
		stock = payload.GetStock()
	}

	sku := existing.SKU
	if payload.GetSku() != "" {
		sku = payload.GetSku()
	}

	categoryID := existing.CategoryID
	if payload.GetCategoryId() != "" {
		categoryID = payload.GetCategoryId()
	}

	log.Printf("✅ [onProductRebranded] Safely rebranding product %s with merged data", payload.GetId())

	return h.products.Rebrand(ctx, payload.GetId(), name, description, basePrice, stock, sku, categoryID, payload.GetActive())
}
func (h integrationHandlers[T]) onProductUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductUpdated)

	domainOpts := make([]models.Option, len(payload.GetOptions()))
	for i, opt := range payload.GetOptions() {
		domainOpts[i] = models.Option{
			Name:  opt.GetName(),
			Value: opt.GetValue(),
			Price: float64(opt.GetPrice()),
		}
	}

	domainAttrs := make([]models.Attribute, len(payload.GetAttributes()))
	for i, opt := range payload.GetAttributes() {
		domainAttrs[i] = models.Attribute{
			Key:   opt.GetKey(),
			Value: opt.GetValue(),
		}
	}

	return h.products.Update(ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetBasePrice(),
		payload.GetCategoryId(),
		payload.GetCategorySlug(),
		payload.GetBrand(),
		payload.GetCondition(),
		payload.GetModel(),
		payload.GetTags(),
		payload.GetManageStocks(),
		payload.GetStock(),
		payload.GetSku(),
		domainAttrs,
		payload.GetWeight(),
		payload.GetHeight(),
		payload.GetWidth(),
		payload.GetDepth(),
		payload.GetStatus(),
		payload.GetNegotiable(),
		payload.GetMiddlemanService(),
		payload.GetUserType(),
		payload.GetShippingCost(),
		payload.GetHasVariants(),
		domainOpts,
		payload.GetThumbnail(),
		float64(payload.GetLat()),
		float64(payload.GetLng()),
	)
}

func (h integrationHandlers[T]) onProductThumbnailAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductThumbnailAdded)
	return h.products.UpdateThumbnail(ctx, payload.GetId(), payload.GetThumbnail())
}
func (h integrationHandlers[T]) onProductThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductThumbnailUpdated)

	// SAFETY CHECK: Only update if thumbnail is not empty
	if payload.GetThumbnail() == "" {
		log.Printf("⚠️ [onProductThumbnailUpdated] Empty thumbnail in update event for product %s - skipping update to preserve existing", payload.GetId())
		return nil
	}

	log.Printf("✅ [onProductThumbnailUpdated] Safely updating thumbnail for product %s", payload.GetId())
	return h.products.UpdateThumbnail(ctx, payload.GetId(), payload.GetThumbnail())
}

func (h integrationHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceIncreased)

	// SAFETY CHECK: Only update if new price is provided
	if payload.GetNewPrice() <= 0 {
		log.Printf("⚠️ [onProductPriceIncreased] Invalid new price in event for product %s - skipping update", payload.GetId())
		return nil
	}

	log.Printf("✅ [onProductPriceIncreased] Using dedicated IncreasePrice method for product %s", payload.GetId())
	return h.products.IncreasePrice(ctx, payload.GetId(), payload.GetNewPrice())
}

func (h integrationHandlers[T]) onProductPriceDecreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPriceDecreased)

	// SAFETY CHECK: Only update if new price is provided
	if payload.GetNewPrice() <= 0 {
		log.Printf("⚠️ [onProductPriceDecreased] Invalid new price in event for product %s - skipping update", payload.GetId())
		return nil
	}

	log.Printf("✅ [onProductPriceDecreased] Using dedicated DecreasePrice method for product %s", payload.GetId())
	return h.products.DecreasePrice(ctx, payload.GetId(), payload.GetNewPrice())
}

func (h integrationHandlers[T]) onProductLeased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductLeased)
	log.Printf("✅ [onProductLeased] Using dedicated MarkAsLeased method for product %s", payload.GetId())
	return h.products.MarkAsLeased(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onProductSold(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductSold)
	log.Printf("✅ [onProductSold] Using dedicated MarkAsSold method for product %s", payload.GetId())
	return h.products.MarkAsSold(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onProductPawned(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductPawned)
	log.Printf("✅ [onProductPawned] Using dedicated MarkAsPawned method for product %s", payload.GetId())
	return h.products.MarkAsPawned(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onProductStockAdjusted(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductStockAdjusted)

	// SAFETY CHECK: Only update if new stock is valid (>= 0)
	if payload.GetNewStock() < 0 {
		log.Printf("⚠️ [onProductStockAdjusted] Invalid new stock in event for product %s - skipping update", payload.GetId())
		return nil
	}

	log.Printf("✅ [onProductStockAdjusted] Using dedicated AdjustStock method for product %s", payload.GetId())
	return h.products.AdjustStock(ctx, payload.GetId(), payload.GetNewStock())
}

func (h integrationHandlers[T]) onProductNegotiableToggled(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductNegotiableToggled)

	// Fetch current product to determine new negotiable state
	existing, err := h.products.Find(ctx, payload.GetId())
	if err != nil {
		log.Printf("⚠️ [onProductNegotiableToggled] Product %s not found - skipping toggle", payload.GetId())
		return nil
	}

	// Toggle the current negotiable state
	newNegotiable := !existing.Negotiable
	log.Printf("✅ [onProductNegotiableToggled] Using dedicated ToggleNegotiable method for product %s (new state: %v)", payload.GetId(), newNegotiable)
	return h.products.ToggleNegotiable(ctx, payload.GetId(), newNegotiable)
}

func (h integrationHandlers[T]) onProductArchived(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductArchived)
	log.Printf("✅ [onProductArchived] Using dedicated ArchiveProduct method for product %s", payload.GetId())
	return h.products.ArchiveProduct(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onPostThumbnailAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostThumbnailAdded)
	return h.posts.UpdateThumbnail(ctx, payload.GetId(), payload.GetThumbnail())
}
func (h integrationHandlers[T]) onPostThumbnailUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostThumbnailUpdated)

	// SAFETY CHECK: Only update if thumbnail is not empty
	if payload.GetThumbnail() == "" {
		log.Printf("⚠️ [onPostThumbnailUpdated] Empty thumbnail in update event for post %s - skipping update to preserve existing", payload.GetId())
		return nil
	}

	log.Printf("✅ [onPostThumbnailUpdated] Safely updating thumbnail for post %s", payload.GetId())
	return h.posts.UpdateThumbnail(ctx, payload.GetId(), payload.GetThumbnail())
}
func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*productspb.ProductRemoved)
	return h.products.Remove(ctx, payload.GetId())
}
func (h integrationHandlers[T]) onOrderReadied(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderReadied)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Ready For Pickup")
}

func (h integrationHandlers[T]) onOrderCanceled(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCanceled)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Canceled")
}

func (h integrationHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCompleted)
	return h.orders.UpdateStatus(ctx, payload.GetId(), "Completed")
}
func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	// Batch fetch: collect all product IDs first
	productIDs := make([]string, len(payload.GetItems()))
	for i, item := range payload.GetItems() {
		productIDs[i] = item.GetProductId()
	}

	// Fetch all products at once to avoid N+1 queries
	var products map[string]*models.Product
	var productErr error

	// Check if repository supports batch fetching
	if batchRepo, ok := h.products.(application.ProductBatchRepository); ok {
		products, productErr = batchRepo.FindBatch(ctx, productIDs)
		if productErr != nil {
			return fmt.Errorf("failed to batch fetch products: %w", productErr)
		}
	} else {
		// Fallback to individual fetches if batch not supported
		products = make(map[string]*models.Product)
		for _, id := range productIDs {
			product, err := h.products.Find(ctx, id)
			if err != nil {
				return fmt.Errorf("failed to fetch product %s: %w", id, err)
			}
			products[id] = product
		}
	}

	// Extract unique seller IDs
	sellerIDs := make([]string, 0, len(products))
	seenSellerIDs := make(map[string]bool)
	for _, product := range products {
		if product != nil && product.UserSellerID != "" && !seenSellerIDs[product.UserSellerID] {
			sellerIDs = append(sellerIDs, product.UserSellerID)
			seenSellerIDs[product.UserSellerID] = true
		}
	}

	// Batch fetch all users (customer + sellers)
	allUserIDs := append([]string{payload.UserCustomerId}, sellerIDs...)
	var users map[string]*models.User
	var userErr error

	// Check if repository supports batch fetching
	if batchRepo, ok := h.users.(application.UserBatchRepository); ok {
		users, userErr = batchRepo.FindBatch(ctx, allUserIDs)
		if userErr != nil {
			return fmt.Errorf("failed to batch fetch users: %w", userErr)
		}
	} else {
		// Fallback to individual fetches if batch not supported
		users = make(map[string]*models.User)
		for _, id := range allUserIDs {
			user, err := h.users.Find(ctx, id)
			if err != nil {
				return fmt.Errorf("failed to fetch user %s: %w", id, err)
			}
			users[id] = user
		}
	}

	// Get customer user
	userCustomer, found := users[payload.UserCustomerId]
	if !found || userCustomer == nil {
		return fmt.Errorf("customer user %s not found", payload.UserCustomerId)
	}

	// Build order items
	var total int64
	items := make([]models.Item, len(payload.GetItems()))
	for i, item := range payload.GetItems() {
		product, found := products[item.GetProductId()]
		if !found || product == nil {
			return fmt.Errorf("product %s not found", item.GetProductId())
		}

		userSeller, found := users[product.UserSellerID]
		if !found || userSeller == nil {
			return fmt.Errorf("seller user %s not found", product.UserSellerID)
		}

		items[i] = models.Item{
			ProductID:      product.ProductID,
			UserSellerID:   userSeller.ID,
			ProductName:    product.Name,
			UserSellerName: userSeller.FirstName,
			Price:          item.Price,
			Quantity:       item.Quantity,
		}
		total += item.Quantity * item.Price
	}

	order := &models.Order{
		OrderID:          payload.GetId(),
		UserCustomerID:   userCustomer.ID,
		UserCustomerName: userCustomer.FirstName,
		Items:            items,
		Total:            total,
		Status:           "New",
	}
	return h.orders.Add(ctx, order)
}

func (h integrationHandlers[T]) onPostAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostAdded)

	domainTags := payload.GetTags()

	return h.posts.Add(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetUserId(),
		domainTags,
		payload.GetStatus(),
		float64(payload.GetLat()),
		float64(payload.GetLng()),
		payload.GetThumbnail(),
		models.PostType,
	)
}

func (h integrationHandlers[T]) onPostRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostRemoved)
	return h.posts.Remove(ctx, payload.GetId())
}

func (h integrationHandlers[T]) onPostUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*postspb.PostUpdated)

	return h.posts.UpdatePost(ctx, payload.GetId(), payload.GetName(), payload.GetDescription(), payload.GetTags(), payload.GetStatus(), payload.GetThumbnail())
}
func (h integrationHandlers[T]) onServiceAdded(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*servicespb.ServiceAdded)

	domainTags := payload.GetTags()
	qualifications := payload.GetQualifications()
	pricing := payload.GetPricing()
	attributes := make([]string, len(payload.GetAttributes()))
	for i, attr := range payload.GetAttributes() {
		attributes[i] = attr.String()
	}
	options := make([]string, len(payload.GetOptions()))
	for i, opt := range payload.GetOptions() {
		options[i] = opt.String()
	}

	return h.services.Add(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetServiceType(),
		payload.GetBasePrice(),
		pricing,
		payload.GetAvailability(),
		payload.GetProviderName(),
		payload.GetUserId(),
		payload.GetCategoryId(),
		payload.GetCategorySlug(),
		payload.GetDescriptionShort(),
		payload.GetDescriptionLong(),
		qualifications,
		payload.GetContact(),
		payload.GetFaq(),
		domainTags,
		models.Status(payload.GetStatus()),
		models.UserType(payload.GetUserType()),
		payload.GetShippingCost(),
		payload.GetNegotiable(),
		payload.GetHasVariants(),
		payload.GetMiddlemanService(),
		attributes,
		options,
		payload.GetThumbnail(),
		float64(payload.GetLat()),
		float64(payload.GetLng()),
		models.ServiceType,
	)
}

func (h integrationHandlers[T]) onServiceUpdated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*servicespb.ServiceUpdated)

	domainTags := payload.GetTags()
	qualifications := payload.GetQualifications()
	pricing := payload.GetPricing()
	attributes := make([]string, len(payload.GetAttributes()))
	for i, attr := range payload.GetAttributes() {
		attributes[i] = attr.String()
	}
	options := make([]string, len(payload.GetOptions()))
	for i, opt := range payload.GetOptions() {
		options[i] = opt.String()
	}

	return h.services.UpdateService(
		ctx,
		payload.GetId(),
		payload.GetName(),
		payload.GetDescription(),
		payload.GetServiceType(),
		payload.GetBasePrice(),
		pricing,
		payload.GetAvailability(),
		payload.GetProviderName(),
		payload.GetUserId(),
		payload.GetCategoryId(),
		payload.GetCategorySlug(),
		payload.GetDescriptionShort(),
		payload.GetDescriptionLong(),
		qualifications,
		payload.GetContact(),
		payload.GetFaq(),
		domainTags,
		models.Status(payload.GetStatus()),
		models.UserType(payload.GetUserType()),
		payload.GetShippingCost(),
		payload.GetNegotiable(),
		payload.GetHasVariants(),
		payload.GetMiddlemanService(),
		attributes,
		options,
		payload.GetThumbnail(),
		float64(payload.GetLat()),
		float64(payload.GetLng()),
		models.ServiceType,
	)
}

func (h integrationHandlers[T]) onServiceRemoved(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*servicespb.ServiceRemoved)
	return h.services.Remove(ctx, payload.GetId())
}
