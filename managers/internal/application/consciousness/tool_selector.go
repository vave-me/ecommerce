package consciousness

import (
	"context"
	"fmt"
	"strings"
	
	"github.com/rs/zerolog"
	"middleman/internal/ddd"
	"middleman/managers/internal/application/consciousness/tool_mappings"
)

// ToolSelector dynamically selects appropriate tools based on event patterns
type ToolSelector struct {
	logger               zerolog.Logger
	eventToolMapping     map[string][]string // event type -> tool names
	patternMapping       map[string][]string // pattern type -> tool names
	comprehensiveMapping *toolmappings.EventToolMapping
}

// NewToolSelector creates a new tool selector with comprehensive mappings
func NewToolSelector(logger zerolog.Logger) *ToolSelector {
	ts := &ToolSelector{
		logger:           logger,
		eventToolMapping: make(map[string][]string),
		patternMapping:   make(map[string][]string),
	}
	
	// Load comprehensive mappings from YAML
	comprehensiveMapping, err := toolmappings.LoadEventToolMappings()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load comprehensive tool mappings, using defaults")
		// Initialize default mappings as fallback
		ts.initializeEventMappings()
		ts.initializePatternMappings()
	} else {
		ts.comprehensiveMapping = comprehensiveMapping
		logger.Info().Msg("Loaded comprehensive tool mappings")
	}
	
	return ts
}

// SelectToolsForEvent returns relevant tool names based on event type and data
func (ts *ToolSelector) SelectToolsForEvent(ctx context.Context, event ddd.Event) []string {
	eventType := event.EventName()
	tools := make(map[string]bool)
	
	// First try comprehensive mappings if available
	if ts.comprehensiveMapping != nil {
		comprehensiveTools := ts.comprehensiveMapping.GetToolsForEvent(eventType)
		for _, tool := range comprehensiveTools {
			tools[tool] = true
		}
	}
	
	// Fallback to default mappings if no comprehensive tools found
	if len(tools) == 0 {
		if eventTools, exists := ts.eventToolMapping[eventType]; exists {
			for _, tool := range eventTools {
				tools[tool] = true
			}
		}
	}
	
	// Analyze event data for additional tool selection
	additionalTools := ts.analyzeEventData(event)
	for _, tool := range additionalTools {
		tools[tool] = true
	}
	
	// Convert map to slice
	result := make([]string, 0, len(tools))
	for tool := range tools {
		result = append(result, tool)
	}
	
	ts.logger.Debug().
		Str("event_type", eventType).
		Int("tool_count", len(result)).
		Strs("tools", result).
		Msg("Selected tools for event")
	
	return result
}

// SelectToolsForPattern returns relevant tool names based on detected pattern
func (ts *ToolSelector) SelectToolsForPattern(pattern *Pattern) []string {
	tools := make(map[string]bool)
	
	// Get tools mapped to pattern type
	if patternTools, exists := ts.patternMapping[pattern.Type]; exists {
		for _, tool := range patternTools {
			tools[tool] = true
		}
	}
	
	// Add tools based on pattern properties
	additionalTools := ts.analyzePatternProperties(pattern)
	for _, tool := range additionalTools {
		tools[tool] = true
	}
	
	// Convert map to slice
	result := make([]string, 0, len(tools))
	for tool := range tools {
		result = append(result, tool)
	}
	
	return result
}

// initializeEventMappings sets up comprehensive event-to-tool mappings
func (ts *ToolSelector) initializeEventMappings() {
	// Order events
	ts.eventToolMapping["OrderCreated"] = []string{
		"order_get_by_id", "order_create", "order_update_status",
		"user_get_by_id", "user_get_orders", "user_update_stats",
		"product_get_by_id", "product_update_stock",
		"notification_send", "notification_create_template",
		"activity_record_event", "activity_get_user_activity",
		"payment_create_intent", "payment_get_methods",
		"shipping_calculate_rates", "shipping_create_label",
	}
	
	ts.eventToolMapping["OrderCancelled"] = []string{
		"order_get_by_id", "order_update_status", "order_refund",
		"payment_create_refund", "payment_update_status",
		"product_update_stock", "product_restore_inventory",
		"notification_send", "notification_cancel_scheduled",
		"user_update_stats", "user_add_credit",
		"shipping_cancel_shipment",
	}
	
	// Payment events
	ts.eventToolMapping["PaymentProcessed"] = []string{
		"payment_get_by_id", "payment_update_status",
		"order_update_payment_status", "order_mark_paid",
		"notification_send_receipt", "notification_send",
		"user_update_payment_history", "user_update_loyalty_points",
		"activity_record_payment",
	}
	
	ts.eventToolMapping["PaymentFailed"] = []string{
		"payment_get_by_id", "payment_retry", "payment_get_failure_reason",
		"order_update_payment_status", "order_hold",
		"notification_send_payment_failure", "notification_send",
		"support_create_ticket", "support_escalate_payment_issue",
		"user_flag_payment_risk",
	}
	
	// User events
	ts.eventToolMapping["UserRegistered"] = []string{
		"user_get_by_id", "user_create_profile", "user_send_welcome",
		"notification_send_welcome", "notification_subscribe_defaults",
		"newsletter_subscribe", "newsletter_send_welcome",
		"offer_create_welcome", "offer_assign_to_user",
		"activity_record_registration",
		"vector_create_user_embedding",
	}
	
	ts.eventToolMapping["UserUpdated"] = []string{
		"user_get_by_id", "user_update_profile",
		"activity_record_profile_update",
		"vector_update_user_embedding",
	}
	
	// Basket events
	ts.eventToolMapping["BasketItemAdded"] = []string{
		"basket_get_by_id", "basket_add_item", "basket_calculate_total",
		"product_get_by_id", "product_check_availability",
		"offer_check_applicable", "offer_apply_to_basket",
		"user_update_browsing_history",
		"activity_record_basket_action",
	}
	
	ts.eventToolMapping["BasketAbandoned"] = []string{
		"basket_get_by_id", "basket_get_items", "basket_calculate_value",
		"user_get_by_id", "user_get_email", "user_get_preferences",
		"notification_send_cart_reminder", "notification_schedule",
		"offer_create_recovery", "offer_send_to_user",
		"activity_record_abandonment",
		"metric_track_abandonment",
	}
	
	// Product events
	ts.eventToolMapping["ProductCreated"] = []string{
		"product_get_by_id", "product_create_variants",
		"category_add_product", "category_update_count",
		"media_upload_images", "media_process_product_images",
		"vector_create_product_embedding", "vector_index_product",
		"activity_record_product_creation",
	}
	
	ts.eventToolMapping["ProductUpdated"] = []string{
		"product_get_by_id", "product_update_details",
		"vector_update_product_embedding",
		"notification_notify_wishlist_users",
		"activity_record_product_update",
	}
	
	ts.eventToolMapping["ProductOutOfStock"] = []string{
		"product_get_by_id", "product_update_status",
		"wishlist_get_users_for_product", "wishlist_notify_users",
		"notification_send_restock_alert",
		"basket_remove_unavailable_items",
		"activity_record_stock_event",
	}
	
	// Support events
	ts.eventToolMapping["TicketCreated"] = []string{
		"support_get_ticket", "support_assign_agent",
		"user_get_by_id", "user_get_support_history",
		"notification_send_ticket_confirmation",
		"activity_record_support_interaction",
	}
	
	ts.eventToolMapping["TicketEscalated"] = []string{
		"support_get_ticket", "support_escalate_to_manager",
		"support_get_escalation_history",
		"notification_send_escalation_alert",
		"user_flag_priority_support",
	}
	
	// Review events
	ts.eventToolMapping["ReviewSubmitted"] = []string{
		"review_get_by_id", "review_moderate",
		"product_update_rating", "product_calculate_average_rating",
		"user_update_review_count", "user_award_points",
		"notification_thank_for_review",
		"activity_record_review",
	}
	
	// Shipping events
	ts.eventToolMapping["ShipmentCreated"] = []string{
		"shipping_get_shipment", "shipping_create_tracking",
		"order_update_shipping_status",
		"notification_send_shipping_confirmation",
		"activity_record_shipment",
	}
	
	ts.eventToolMapping["ShipmentDelivered"] = []string{
		"shipping_get_shipment", "shipping_confirm_delivery",
		"order_mark_delivered", "order_complete",
		"notification_send_delivery_confirmation",
		"review_request_product_review",
		"activity_record_delivery",
	}
	
	// Newsletter events
	ts.eventToolMapping["NewsletterSubscribed"] = []string{
		"newsletter_add_subscriber", "newsletter_send_confirmation",
		"user_update_preferences",
		"activity_record_subscription",
	}
	
	// Offer events
	ts.eventToolMapping["OfferCreated"] = []string{
		"offer_get_by_id", "offer_set_rules",
		"user_get_eligible_users", "user_notify_offer",
		"notification_send_offer_announcement",
		"activity_record_offer_creation",
	}
	
	// Message events
	ts.eventToolMapping["MessageSent"] = []string{
		"message_get_by_id", "message_deliver",
		"user_get_recipient", "user_check_blocked",
		"notification_send_message_alert",
		"activity_record_message",
	}
	
	// Comment events
	ts.eventToolMapping["CommentPosted"] = []string{
		"comment_get_by_id", "comment_moderate",
		"post_update_comment_count",
		"user_notify_mentions", "user_update_engagement",
		"notification_send_comment_alert",
		"activity_record_comment",
	}
}

// initializePatternMappings sets up pattern-to-tool mappings
func (ts *ToolSelector) initializePatternMappings() {
	// Fraud patterns
	ts.patternMapping["fraud_risk"] = []string{
		"order_flag_for_review", "order_hold",
		"payment_verify", "payment_block",
		"user_flag_suspicious", "user_check_history",
		"activity_get_suspicious_patterns",
		"support_create_fraud_case",
		"notification_alert_security_team",
	}
	
	// Abandonment patterns
	ts.patternMapping["cart_abandonment"] = []string{
		"basket_get_abandoned", "basket_calculate_recovery_value",
		"offer_create_recovery", "offer_send_personalized",
		"notification_send_recovery_email",
		"user_get_abandonment_history",
		"metric_track_recovery_attempt",
	}
	
	// Support crisis patterns
	ts.patternMapping["support_crisis"] = []string{
		"support_get_open_tickets", "support_batch_assign",
		"support_enable_chat_bot", "support_alert_all_agents",
		"notification_send_wait_time_alert",
		"metric_track_response_times",
		"activity_log_crisis_event",
	}
	
	// Activity surge patterns
	ts.patternMapping["activity_surge"] = []string{
		"activity_get_surge_details", "activity_analyze_patterns",
		"metric_track_surge", "metric_alert_operations",
		"notification_alert_team",
		"service_scale_resources",
	}
	
	// User surge patterns
	ts.patternMapping["user_surge"] = []string{
		"user_get_registration_rate", "user_batch_process",
		"notification_queue_welcome_emails",
		"offer_prepare_bulk_distribution",
		"service_scale_authentication",
		"metric_track_user_surge",
	}
	
	// Cancellation patterns
	ts.patternMapping["cancellation_spike"] = []string{
		"order_get_cancellation_reasons", "order_analyze_patterns",
		"product_check_issues", "product_flag_problematic",
		"support_proactive_outreach",
		"notification_send_retention_offers",
		"metric_analyze_cancellation_trends",
	}
}

// analyzeEventData performs deep analysis of event data to select additional tools
func (ts *ToolSelector) analyzeEventData(event ddd.Event) []string {
	tools := []string{}
	
	// Extract event data from payload
	var eventData map[string]interface{}
	
	// Try to extract data from the payload
	if payload := event.Payload(); payload != nil {
		// Try to cast payload to map
		if data, ok := payload.(map[string]interface{}); ok {
			eventData = data
		} else {
			// If payload is a struct, we'd need reflection or specific type assertions
			// For now, return empty tools
			return tools
		}
	}
	
	if eventData == nil {
		return tools
	}
	
	// Check for specific fields that indicate tool needs
	if userID, ok := eventData["user_id"].(string); ok && userID != "" {
		tools = append(tools, "user_get_by_id", "user_get_preferences")
	}
	
	if productID, ok := eventData["product_id"].(string); ok && productID != "" {
		tools = append(tools, "product_get_by_id", "product_get_availability")
	}
	
	if amount, ok := eventData["amount"].(float64); ok && amount > 1000 {
		// High value transaction - add fraud detection tools
		tools = append(tools, "payment_verify_high_value", "user_check_payment_history")
	}
	
	if priority, ok := eventData["priority"].(string); ok && priority == "urgent" {
		// Urgent items need immediate action tools
		tools = append(tools, "notification_send_immediate", "support_escalate_urgent")
	}
	
	// Check for sentiment in support/review events
	if sentiment, ok := eventData["sentiment_score"].(float64); ok && sentiment < -0.5 {
		// Negative sentiment - add retention tools
		tools = append(tools, "support_flag_negative", "offer_create_retention")
	}
	
	return tools
}

// analyzePatternProperties analyzes pattern properties for additional tool selection
func (ts *ToolSelector) analyzePatternProperties(pattern *Pattern) []string {
	tools := []string{}
	
	// High confidence patterns may need immediate action tools
	if pattern.Confidence > 0.9 {
		tools = append(tools, "notification_send_immediate", "activity_log_high_confidence")
	}
	
	// Check pattern properties
	for key, value := range pattern.Properties {
		switch key {
		case "urgency":
			if urgency, ok := value.(string); ok && urgency == "high" {
				tools = append(tools, "support_escalate_immediate", "notification_alert_manager")
			}
		case "value":
			if val, ok := value.(float64); ok && val > 5000 {
				tools = append(tools, "payment_verify_high_value", "order_flag_vip")
			}
		case "user_type":
			if userType, ok := value.(string); ok && userType == "vip" {
				tools = append(tools, "user_get_vip_preferences", "offer_create_exclusive")
			}
		}
	}
	
	// Multi-dimensional patterns may need complex analysis tools
	if len(pattern.Dimensions) > 3 {
		tools = append(tools, "metric_complex_analysis", "vector_multi_dimensional_search")
	}
	
	return tools
}

// GetAllAvailableTools returns all 415+ tool names that can be dynamically selected
func (ts *ToolSelector) GetAllAvailableTools() []string {
	// This would connect to the ComprehensiveToolRegistry to get all available tools
	// For now, returning a comprehensive list based on the handlers
	allTools := []string{}
	
	// Collect all unique tools from mappings
	toolSet := make(map[string]bool)
	
	for _, tools := range ts.eventToolMapping {
		for _, tool := range tools {
			toolSet[tool] = true
		}
	}
	
	for _, tools := range ts.patternMapping {
		for _, tool := range tools {
			toolSet[tool] = true
		}
	}
	
	// Convert to slice
	for tool := range toolSet {
		allTools = append(allTools, tool)
	}
	
	return allTools
}

// IsToolRelevant checks if a tool is relevant for a given context
func (ts *ToolSelector) IsToolRelevant(toolName string, eventType string, pattern *Pattern) bool {
	// Check event mapping
	if tools, exists := ts.eventToolMapping[eventType]; exists {
		for _, tool := range tools {
			if tool == toolName {
				return true
			}
		}
	}
	
	// Check pattern mapping
	if pattern != nil && pattern.Type != "" {
		if tools, exists := ts.patternMapping[pattern.Type]; exists {
			for _, tool := range tools {
				if tool == toolName {
					return true
				}
			}
		}
	}
	
	// Check if tool name contains relevant keywords
	relevantKeywords := ts.extractKeywords(eventType, pattern)
	for _, keyword := range relevantKeywords {
		if strings.Contains(strings.ToLower(toolName), keyword) {
			return true
		}
	}
	
	return false
}

// extractKeywords extracts relevant keywords from event and pattern
func (ts *ToolSelector) extractKeywords(eventType string, pattern *Pattern) []string {
	keywords := []string{}
	
	// Extract from event type
	eventLower := strings.ToLower(eventType)
	if strings.Contains(eventLower, "order") {
		keywords = append(keywords, "order", "payment", "shipping")
	}
	if strings.Contains(eventLower, "user") {
		keywords = append(keywords, "user", "profile", "preference")
	}
	if strings.Contains(eventLower, "product") {
		keywords = append(keywords, "product", "inventory", "stock")
	}
	
	// Extract from pattern
	if pattern != nil {
		patternLower := strings.ToLower(pattern.Type)
		if strings.Contains(patternLower, "fraud") {
			keywords = append(keywords, "verify", "flag", "suspicious")
		}
		if strings.Contains(patternLower, "abandon") {
			keywords = append(keywords, "recovery", "reminder", "offer")
		}
	}
	
	return keywords
}