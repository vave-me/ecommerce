package handlers

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	content "google.golang.org/api/content/v2.1"
	"middleman/internal/am"
	"middleman/internal/ddd"
	"middleman/internal/errorsotel"
	"middleman/internal/registry"
	"middleman/merchant/internal/application"
	"middleman/products/productspb"
	"time"
)

type integrationHandlers[T ddd.Event] struct {
	app application.App
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(reg registry.Registry, app application.App, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app: app,
	}, zerolog.Logger{}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(productspb.ProductAggregateChannel, handlers, am.MessageFilter{
		productspb.ProductAddedEvent,
		productspb.ProductUpdatedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductThumbnailUpdatedEvent,
		productspb.ProductRemovedEvent,
	}, am.GroupName("products-merchant"))
	return err
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
	case productspb.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case productspb.ProductUpdatedEvent,
		productspb.ProductRebrandedEvent,
		productspb.ProductPriceIncreasedEvent,
		productspb.ProductPriceDecreasedEvent,
		productspb.ProductStockAdjustedEvent,
		productspb.ProductThumbnailUpdatedEvent:
		return h.onProductUpdated(ctx, event)
	case productspb.ProductRemovedEvent:
		return h.onProductRemoved(ctx, event)
	}
	return nil
}

func (h integrationHandlers[T]) onProductAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*productspb.ProductAdded)

	// REQUIRED FIELDS per Google's spec:
	//   ID [id], Title [title]/Structured title [structured_title],
	//   Description [description]/Structured description [structured_description],
	//   Link [link], Image link [image_link], Price [price], Availability [availability].

	// 1) Construct a unique product "Id" used by Merchant Center.
	//    Commonly: "channel:language:offerId" or some variant:
	gmcID := fmt.Sprintf("online:en:%s", payload.Id)

	// 2) Title. If you have a standard text-based title, set `Title`.
	//    If you have a generative-AI-based title, use StructuredTitle.
	//    For demonstration, we'll just fill the normal `Title`.
	title := payload.Name

	// 3) Description. Similarly, choose either normal `Description` or `StructuredDescription`.
	description := payload.Description

	// 4) Link [link]: A landing page on your website (if available).
	//    For demonstration, we build a dummy link. Replace with your real site’s URL.
	link := fmt.Sprintf("https://www.example.com/products/%s", payload.Id)

	// 5) Image link [image_link]: Must be a valid, crawlable URL.
	//    We'll assume payload.Product.Thumbnail is a full image URL, or build one:
	imageLink := payload.Thumbnail
	if imageLink == "" {
		imageLink = "https://www.example.com/images/default.jpg"
	}

	// 6) Price [price]:
	//    Merchant Center requires "Value" + "Currency" in ISO 4217.
	//    If base_price is stored as integer cents, convert to decimal string.
	priceValue := fmt.Sprintf("%.2f", float64(payload.BasePrice)/100.0)
	price := &content.Price{
		Value:    priceValue,
		Currency: "EUR", // Adjust to your local currency if needed
	}

	// 7) Availability [availability]:
	//    Map your product "status" & "stock" to "in_stock", "out_of_stock", "preorder", etc.
	availability := "out_of_stock"
	if (payload.Status == "ACTIVE" || payload.Status == "PUBLISHED") && payload.Stock > 0 {
		availability = "in_stock"
	}

	// OPTIONAL / RECOMMENDED FIELDS:

	// brand [brand]: Up to 70 characters
	brand := payload.Brand

	// condition [condition]: "new", "used", "refurbished"
	condition := payload.Condition
	if condition == "" {
		condition = "new" // default if unspecified
	}

	// mpn [mpn]: Provide if there's no GTIN, or if you have a known MPN:
	mpn := payload.Model
	// or if you store a separate MPN field.
	// "model" might be your MPN. Adjust as needed.

	// identifier_exists [identifier_exists]:
	// If you do NOT have brand/GTIN/MPN, set to "no". Otherwise "yes".
	identifierExists := true
	if brand == "" && mpn == "" {
		identifierExists = false
	}

	// google_product_category [google_product_category]:
	// If you have a mapped numeric category ID, set it; else the text path.
	// e.g. if "category_id" is "Apparel & Accessories > Shoes", or the numeric ID "187".
	googleProductCategory := payload.CategoryId // Placeholder for your real mapping

	// Multipack [multipack]: If you define a merchant multipack
	var multipack int64
	if payload.HasVariants {
		// Example: if "has_variants" means you have multiple items packaged
		// but let's assume 1 if you don't actually do a multi-quantity pack
		multipack = 1
	}

	// is_bundle [is_bundle]: If you create a merchant-defined custom group of different items
	isBundle := false // For now, assume no. If you do custom bundling, set to true.

	// adult [adult]: Mark yes if product is adult-oriented. We'll default to "no"
	adult := false

	// shipping [shipping]: If you want to override shipping.
	// For demonstration, if shipping_cost is known, define a single shipping rule:
	var shipping []*content.ProductShipping
	if payload.ShippingCost > 0 {
		val := fmt.Sprintf("%.2f", float64(payload.ShippingCost)/100.0)
		shipping = []*content.ProductShipping{
			{
				Country: "US",
				Price: &content.Price{
					Value:    val,
					Currency: "USD",
				},
				Service: "Standard Shipping",
				// You could also define min/max transit times if known
			},
		}
	}

	// item_group_id [item_group_id]: If item is part of a variant group
	// that differ by color, size, etc.
	var itemGroupID string
	if payload.HasVariants {
		itemGroupID = payload.Id // or any grouping key
	}

	// product_weight [product_weight], product_height [product_height],
	// product_width [product_width], product_length [product_length].
	// Convert from your int64 fields to "X unit" strings, e.g. "3.5 lb" or "20 in".
	// For simplicity, we assume you want to store them as "in" or "lb".
	// If your data is in millimeters, etc., you'll need to convert.
	var productWeight *content.ProductWeight
	if payload.Weight > 0 {
		// example: treat weight as ounces or grams

		productWeight = &content.ProductWeight{
			Value: float64(payload.Weight),
			Unit:  "g", // or "lb", "kg" etc.
		}
	}

	// Additional image link [additional_image_link]
	// If you have more images, you can store them in a slice:
	// additionalImageLinks := []string{"https://www.example.com/img2.jpg", "https://www.example.com/img3.jpg"}

	// Build final content.Product struct
	pr := &content.Product{
		// Basic product data
		Id:                    gmcID,
		OfferId:               payload.Id,   // required
		Title:                 title,        // required
		Description:           description,  // required
		Link:                  link,         // required
		ImageLink:             imageLink,    // required
		ContentLanguage:       "de",         // typically your target language
		TargetCountry:         "DE",         // typically your target country
		Price:                 price,        // required
		Availability:          availability, // required
		Brand:                 brand,        // required if product has a recognized brand
		Condition:             condition,    // "new", "used", or "refurbished"
		Mpn:                   mpn,          // manufacturer part number if no GTIN
		IdentifierExists:      identifierExists,
		GoogleProductCategory: googleProductCategory, // numeric ID or full path if known
		Multipack:             multipack,             // number of identical products in a single multipack
		IsBundle:              isBundle,              // if it's a merchant-defined bundle
		Adult:                 adult,                 // "true" if adult content
		Shipping:              shipping,              // override shipping if needed
		ProductWeight:         productWeight,
		ItemGroupId:           itemGroupID, // For variants
	}

	// 2) Invoke your application layer’s UpsertProduct command
	return h.app.UpsertProduct(ctx, application.UpsertProductCommand{
		Product: pr,
	})
}

func (h integrationHandlers[T]) onProductRemoved(ctx context.Context, event T) error {
	payload := event.Payload().(*productspb.ProductRemoved)
	return h.app.DeleteProduct(ctx, application.DeleteProductCommand{
		ProductID: payload.Id,
	})
}

func (h integrationHandlers[T]) onProductUpdated(ctx context.Context, event T) error {
	// Extract product ID from the event
	var productID string
	switch p := event.Payload().(type) {
	case *productspb.ProductUpdated:
		productID = p.Id
	case *productspb.ProductRebranded:
		productID = p.Id
	case *productspb.ProductPriceIncreased:
		productID = p.Id
	case *productspb.ProductPriceDecreased:
		productID = p.Id
	case *productspb.ProductStockAdjusted:
		productID = p.Id
	case *productspb.ProductThumbnailUpdated:
		productID = p.Id
	default:
		return fmt.Errorf("unknown event payload type: %T", p)
	}
	
	// For updates, we need to fetch the full product data and re-sync
	// In a real implementation, you would fetch this from your product service
	// For now, we'll just log the product ID that needs updating
	
	// TODO: Implement fetching full product data from product service
	// and then sync to Google Merchant Center
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Str("product_id", productID).
		Str("event_type", event.EventName()).
		Msg("product update event received, sync required")
	
	return nil
}
