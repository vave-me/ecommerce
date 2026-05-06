package tools

import (
	ai2 "middleman/internal/ai"
)

// CreateOpenAICompliantTools creates only essential tool definitions to stay under 128 limit
func CreateOpenAICompliantTools() []ai2.Tool {
	tools := []ai2.Tool{}

	// CRUCIAL TOOLS (as specified by user)
	// Keep these but with reduced methods
	tools = append(tools, createCategoryTools()...)     // ~4 tools
	tools = append(tools, createNewsletterTools()...)   // ~3 tools
	tools = append(tools, createUserTools()...)         // ~4 tools
	tools = append(tools, createProductTools()...)      // ~6 tools
	tools = append(tools, createPostTools()...)         // ~5 tools
	tools = append(tools, createServiceTools()...)      // ~4 tools
	tools = append(tools, createBasketTools()...)       // ~5 tools
	tools = append(tools, createOrderTools()...)        // ~5 tools
	tools = append(tools, createShippingTools()...)     // ~4 tools
	
	// ALSO CRUCIAL (as specified by user)
	tools = append(tools, createMediaTools()...)        // ~3 tools
	tools = append(tools, createCommentTools()...)      // ~3 tools
	tools = append(tools, createWishlistTools()...)     // ~3 tools
	tools = append(tools, createNotificationTools()...) // ~3 tools
	tools = append(tools, createMailerTools()...)       // ~2 tools
	tools = append(tools, createActivityTools()...)     // ~2 tools

	// Total: ~56 essential tools (well under 128 limit)
	
	// COMMENTED OUT - Less crucial tools
	// tools = append(tools, createFollowingTools()...)
	// tools = append(tools, createGeocodingTools()...)
	// tools = append(tools, createMessageTools()...)
	// tools = append(tools, createMetricTools()...)
	// tools = append(tools, createOfferTools()...)
	// tools = append(tools, createPaymentTools()...)
	// tools = append(tools, createReviewTools()...)
	// tools = append(tools, createSupportTools()...)
	// tools = append(tools, createVectorTools()...)

	return tools
}