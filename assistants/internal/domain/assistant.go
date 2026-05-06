package domain

import (
	"context"
	"fmt"

	"middleman/assistants/internal/constants"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"strings"
	"time"

	"github.com/stackus/errors"
)

const AssistantAggregate = "assistants.Assistant"

var (
	ErrAssistantIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "the assistant id cannot be blank")
	ErrAssistantNameCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the assistant name cannot be blank")
	ErrInvalidAssistantCapability = errors.Wrap(errors.ErrBadRequest, "invalid assistant capability")
	ErrAssistantAlreadyActive     = errors.Wrap(errors.ErrBadRequest, "the assistant is already active")
	ErrAssistantAlreadyInactive   = errors.Wrap(errors.ErrBadRequest, "the assistant is already inactive")
)

// AssistantType represents the type of assistant
type AssistantType string

const (
	AssistantTypeStandard  AssistantType = "standard"
	AssistantTypeAdmin     AssistantType = "admin"
	AssistantTypeBusiness  AssistantType = "business"
	AssistantTypeSupport   AssistantType = "support"
	AssistantTypeScheduler AssistantType = "scheduler"
)

// IsValid checks if the assistant type is valid
func (t AssistantType) IsValid() bool {
	switch t {
	case AssistantTypeStandard, AssistantTypeAdmin, AssistantTypeBusiness, AssistantTypeSupport, AssistantTypeScheduler:
		return true
	default:
		return false
	}
}

// AssistantCapability represents what the assistant can do
type AssistantCapability string

const (
	CapabilityAssistantManagement AssistantCapability = "assistant_management"
	CapabilityUserInteraction     AssistantCapability = "user_interaction"
	CapabilityDataAnalysis        AssistantCapability = "data_analysis"
	CapabilityLocationServices    AssistantCapability = "location_services"
	CapabilityAuthentication      AssistantCapability = "authentication"
	CapabilityPublicAPIAccess     AssistantCapability = "public_api_access"
	CapabilityJailbreakResistant  AssistantCapability = "jailbreak_resistant"
	CapabilityScopeEnforcement    AssistantCapability = "scope_enforcement"
	CapabilityDataRetrieval       AssistantCapability = "data_retrieval"
	CapabilitySearchAndFilter     AssistantCapability = "search_and_filter"
	CapabilityPrivateAPIAccess    AssistantCapability = "private_api_access"
	CapabilityUserDataAccess      AssistantCapability = "user_data_access"
	CapabilityTokenManagement     AssistantCapability = "token_management"
	CapabilityDataMasking         AssistantCapability = "data_masking"
	CapabilityAuditLogging        AssistantCapability = "audit_logging"

	// Additional capabilities for LLM processing
	CapabilityTextGeneration AssistantCapability = "text_generation"
	CapabilityCodeGeneration AssistantCapability = "code_generation"
	CapabilityWebSearch      AssistantCapability = "web_search"

	// System and automation capabilities
	CapabilitySystemConfiguration AssistantCapability = "system_configuration"
	CapabilityTaskAutomation      AssistantCapability = "task_automation"
)

// SecurityContext represents the security context for private operations
type SecurityContext struct {
	UserID        string
	SessionToken  string
	AccessToken   string
	Permissions   []string
	DataScope     []string
	ExpiresAt     time.Time
	IPAddress     string
	UserAssistant string
	TrustLevel    string
	MFAVerified   bool
	ConsentGiven  map[string]bool
}

func (Assistant) Key() string { return AssistantAggregate }

// PrivateEndpointConfig defines configuration for private endpoint access
type PrivateEndpointConfig struct {
	Endpoint        string   `json:"endpoint"`
	RequiredScopes  []string `json:"required_scopes"`
	MinTrustLevel   string   `json:"min_trust_level"`
	RequiresMFA     bool     `json:"requires_mfa"`
	DataSensitivity string   `json:"data_sensitivity"` // low, medium, high, critical
	AuditLevel      string   `json:"audit_level"`      // basic, detailed, full
}

// AssistantAction represents an action the assistant can perform
type AssistantAction struct {
	Type        string                 `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Description string                 `json:"description"`
}

// Assistant represents an AI assistant that can interact with the service
type Assistant struct {
	es.Aggregate
	Name            string
	Description     string
	Type            AssistantType // Type of assistant (admin, business, support, etc.)
	Model           string        // AI model to use for this assistant
	Capabilities    []AssistantCapability
	Active          bool
	Temperature     float64
	MaxTokens       int
	SystemPrompt    string
	UserID          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	requestAnalyzer RequestAnalyzer
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Assistant)(nil)

func NewAssistant(id string) *Assistant {
	return &Assistant{
		Aggregate: es.NewAggregate(id, AssistantAggregate),
	}
}

// ApplyEvent applies events to rebuild the assistant state
func (a *Assistant) ApplyEvent(event ddd.Event) error {
	switch event.EventName() {
	case AssistantCreatedEvent:
		return a.applyAssistantCreated(event)
	case AssistantActivatedEvent:
		return a.applyAssistantActivated(event)
	case AssistantDeactivatedEvent:
		return a.applyAssistantDeactivated(event)
	case AssistantConfigurationUpdatedEvent:
		return a.applyAssistantConfigurationUpdated(event)
	case AssistantRequestProcessedEvent:
		// Request processed events don't change assistant state
		return nil
	default:
		return fmt.Errorf("unknown event: %s", event.EventName())
	}
}

func (a *Assistant) applyAssistantCreated(event ddd.Event) error {
	if payload, ok := event.Payload().(*AssistantCreated); ok {
		a.Name = payload.Name
		a.Description = payload.Description
		a.Type = payload.Type
		a.Model = payload.Model
		a.UserID = payload.UserID
		a.Capabilities = deduplicateCapabilities(payload.Capabilities)
		a.Temperature = payload.Temperature
		a.MaxTokens = payload.MaxTokens
		a.SystemPrompt = payload.SystemPrompt
		a.Active = payload.Active
		a.CreatedAt = payload.CreatedAt
		a.UpdatedAt = payload.CreatedAt
		return nil
	}
	return fmt.Errorf("invalid payload type for AssistantCreated event")
}

func (a *Assistant) applyAssistantActivated(event ddd.Event) error {
	if payload, ok := event.Payload().(*AssistantActivated); ok {
		a.Active = payload.Active
		a.UpdatedAt = payload.Timestamp
		return nil
	}
	return fmt.Errorf("invalid payload type for AssistantActivated event")
}

func (a *Assistant) applyAssistantDeactivated(event ddd.Event) error {
	if payload, ok := event.Payload().(*AssistantDeactivated); ok {
		a.Active = payload.Active
		a.UpdatedAt = payload.Timestamp
		return nil
	}
	return fmt.Errorf("invalid payload type for AssistantDeactivated event")
}

func (a *Assistant) applyAssistantConfigurationUpdated(event ddd.Event) error {
	if payload, ok := event.Payload().(*AssistantConfigurationUpdated); ok {
		a.Temperature = payload.Temperature
		a.MaxTokens = payload.MaxTokens
		a.SystemPrompt = payload.SystemPrompt
		if len(payload.Capabilities) > 0 {
			a.Capabilities = deduplicateCapabilities(payload.Capabilities)
		}
		a.UpdatedAt = payload.UpdatedAt
		return nil
	}
	return fmt.Errorf("invalid payload type for AssistantConfigurationUpdated event")
}

// ToSnapshot creates a snapshot of the current assistant state
func (a *Assistant) ToSnapshot() es.Snapshot {
	return &AssistantV1{
		Name:         a.Name,
		Description:  a.Description,
		Type:         a.Type,
		Model:        a.Model,
		UserID:       a.UserID,
		Capabilities: a.Capabilities,
		Active:       a.Active,
		Temperature:  a.Temperature,
		MaxTokens:    a.MaxTokens,
		SystemPrompt: a.SystemPrompt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// ApplySnapshot applies a snapshot to restore the assistant state
func (a *Assistant) ApplySnapshot(snapshot es.Snapshot) error {
	if s, ok := snapshot.(*AssistantV1); ok {
		a.Name = s.Name
		a.Description = s.Description
		a.Type = s.Type
		a.Model = s.Model
		a.UserID = s.UserID
		a.Capabilities = deduplicateCapabilities(s.Capabilities)
		a.Active = s.Active
		a.Temperature = s.Temperature
		a.MaxTokens = s.MaxTokens
		a.SystemPrompt = s.SystemPrompt
		a.CreatedAt = s.CreatedAt
		a.UpdatedAt = s.UpdatedAt
		// Ensure analyzer is initialized
		if a.requestAnalyzer == nil {
			// Request analyzer will be injected
		}
		return nil
	}
	return fmt.Errorf("invalid snapshot type for Assistant")
}

// deduplicateCapabilities removes duplicate capabilities from a slice
func deduplicateCapabilities(capabilities []AssistantCapability) []AssistantCapability {
	seen := make(map[AssistantCapability]bool)
	result := make([]AssistantCapability, 0, len(capabilities))

	for _, cap := range capabilities {
		if !seen[cap] && isValidCapability(cap) {
			seen[cap] = true
			result = append(result, cap)
		}
	}

	return result
}

// isValidCapability checks if a capability is valid
func isValidCapability(cap AssistantCapability) bool {
	validCapabilities := map[AssistantCapability]bool{
		CapabilityAssistantManagement: true,
		CapabilityUserInteraction:     true,
		CapabilityDataAnalysis:        true,
		CapabilityLocationServices:    true,
		CapabilityAuthentication:      true,
		CapabilityPublicAPIAccess:     true,
		CapabilityJailbreakResistant:  true,
		CapabilityScopeEnforcement:    true,
		CapabilityDataRetrieval:       true,
		CapabilitySearchAndFilter:     true,
		CapabilityPrivateAPIAccess:    true,
		CapabilityUserDataAccess:      true,
		CapabilityTokenManagement:     true,
		CapabilityDataMasking:         true,
		CapabilityAuditLogging:        true,
		CapabilityTextGeneration:      true,
		CapabilityCodeGeneration:      true,
		CapabilityWebSearch:           true,
		CapabilitySystemConfiguration: true,
		CapabilityTaskAutomation:      true,
	}
	return validCapabilities[cap]
}

// SetRequestAnalyzer sets the request analyzer for the assistant
func (a *Assistant) SetRequestAnalyzer(analyzer RequestAnalyzer) {
	a.requestAnalyzer = analyzer
}

func (a *Assistant) CreateAssistant(id, name, description, userID string, assistantType AssistantType, capabilities []AssistantCapability, temperature float64, maxTokens int, systemPrompt string) (ddd.Event, error) {

	// FIXED: Prevent creation if assistant is already initialized
	if a.Name != "" || len(a.Capabilities) > 0 || a.Version() > 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "assistant is already created - cannot recreate")
	}

	// Validate assistant type
	if assistantType == "" {
		assistantType = AssistantTypeStandard
	}
	if !assistantType.IsValid() {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid assistant type")
	}

	if name == "" {
		name = "vaver"
	}
	if description == "" {
		description = "AI-powered marketplace assistant for intelligent search and recommendations"
	}

	// Set default model if not specified - use a model that supports function calling
	model := "gpt-4o" // GPT-4o supports function calling

	if capabilities == nil {
		capabilities = []AssistantCapability{
			CapabilityUserInteraction,
			CapabilityDataAnalysis,
			CapabilitySearchAndFilter,
			CapabilityDataRetrieval,
			CapabilityPublicAPIAccess,
			CapabilityJailbreakResistant,
			CapabilityScopeEnforcement,
		}
	}
	if temperature == 0.0 {
		temperature = 0.7
	}
	if maxTokens == 0 {
		maxTokens = 1000
	}
	if systemPrompt == "" {
		// Use the comprehensive system prompt with tool instructions
		systemPrompt = constants.SystemPrompt
	}

	for _, cap := range capabilities {
		if !isValidCapability(cap) {
			return nil, ErrInvalidAssistantCapability
		}
	}

	now := time.Now()

	eventPayload := &AssistantCreated{
		Name:         name,
		Description:  description,
		Type:         assistantType,
		Model:        model,  // Include model in event
		UserID:       userID, // ADDED: Include UserID in event
		Capabilities: capabilities,
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		Active:       true,
		CreatedAt:    now,
	}

	a.AddEvent(AssistantCreatedEvent, eventPayload)

	return ddd.NewEvent(AssistantCreatedEvent, a), nil
}

func (a *Assistant) ProcessInput(ctx context.Context, requestID, userID, message string, context map[string]interface{}, timestamp time.Time, requestType string) (ddd.Event, error) {
	if a.requestAnalyzer == nil {
		// Request analyzer will be injected
	}

	// Use unified analyzer for request processing
	analysisReq := AnalysisRequest{
		RequestID:    requestID,
		UserID:       userID,
		Message:      message,
		Context:      context,
		Timestamp:    timestamp,
		RequestType:  requestType,
		Capabilities: a.Capabilities,
		SecurityCtx:  nil, // Public request
	}

	result, err := a.requestAnalyzer.AnalyzeRequest(ctx, analysisReq)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze request: %w", err)
	}

	// Generate response based on analysis
	responseID := fmt.Sprintf("resp_%d", timestamp.UnixNano())
	responseMessage := a.generateResponseFromAnalysis(result, requestID, userID, message, timestamp, requestType)
	responseConfidence := result.Confidence

	// Create and add the event
	eventPayload := &AssistantRequestProcessed{
		RequestID:        requestID,
		UserID:           userID,
		Message:          message,
		Context:          context,
		RequestType:      requestType,
		RequestTimestamp: timestamp,
		ResponseID:       responseID,
		ResponseMessage:  responseMessage,
		ResponseData: map[string]interface{}{
			"entity_type":     result.EntityType,
			"intent":          result.Intent,
			"security_status": result.SecurityStatus,
		},
		ResponseActions:    result.Actions,
		ResponseTimestamp:  timestamp,
		ResponseStatus:     "completed",
		ResponseConfidence: responseConfidence,
	}

	a.AddEvent(AssistantRequestProcessedEvent, eventPayload)

	return ddd.NewEvent(AssistantRequestProcessedEvent, a), nil
}

// generateResponseFromAnalysis creates response message from analysis result
func (a *Assistant) generateResponseFromAnalysis(result *AnalysisResult, requestID, userID, message string, timestamp time.Time, requestType string) string {
	if result.SecurityStatus == "blocked" {
		return "I cannot process this request due to security policy violations."
	}

	if result.SecurityStatus == "out_of_scope" {
		return "This request is outside my authorized scope. Please ask something related to my capabilities."
	}

	// Generate response based on intent and entity type
	switch result.Intent {
	case "search":
		return fmt.Sprintf("I found %d %s matching your criteria. Here are the results.", len(result.Actions), result.EntityType)
	case "add":
		return fmt.Sprintf("I'll help you add a new %s entry.", strings.TrimSuffix(result.EntityType, "s"))
	case "analyze":
		return fmt.Sprintf("I've analyzed the %s data as requested.", result.EntityType)
	case "recommend":
		return fmt.Sprintf("Based on your preferences, here are my %s recommendations.", result.EntityType)
	default:
		return fmt.Sprintf("I've processed your %s request with %d%% confidence.", result.EntityType, int(result.Confidence*100))
	}
}

func (a *Assistant) Activate() (ddd.Event, error) {
	if a.Active {
		return nil, ErrAssistantAlreadyActive
	}

	eventPayload := &AssistantActivated{
		Active:    true,
		Timestamp: time.Now(),
	}

	a.AddEvent(AssistantActivatedEvent, eventPayload)

	return ddd.NewEvent(AssistantActivatedEvent, a), nil
}

func (a *Assistant) Deactivate() (ddd.Event, error) {
	if !a.Active {
		return nil, ErrAssistantAlreadyInactive
	}

	eventPayload := &AssistantDeactivated{
		Active:    false,
		Timestamp: time.Now(),
	}

	a.AddEvent(AssistantDeactivatedEvent, eventPayload)

	return ddd.NewEvent(AssistantDeactivatedEvent, a), nil
}

func (a *Assistant) UpdateConfiguration(temperature float64, maxTokens int, systemPrompt string) (ddd.Event, error) {
	now := time.Now()

	eventPayload := &AssistantConfigurationUpdated{
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		UpdatedAt:    now,
	}

	a.AddEvent(AssistantConfigurationUpdatedEvent, eventPayload)

	return ddd.NewEvent(AssistantConfigurationUpdatedEvent, a), nil
}

// UpdateConfigurationWithCapabilities updates configuration including capabilities
func (a *Assistant) UpdateConfigurationWithCapabilities(temperature float64, maxTokens int, systemPrompt string, capabilities []AssistantCapability) (ddd.Event, error) {
	now := time.Now()

	// Deduplicate and validate capabilities
	validatedCapabilities := deduplicateCapabilities(capabilities)

	eventPayload := &AssistantConfigurationUpdated{
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		Capabilities: validatedCapabilities,
		UpdatedAt:    now,
	}

	a.AddEvent(AssistantConfigurationUpdatedEvent, eventPayload)

	return ddd.NewEvent(AssistantConfigurationUpdatedEvent, a), nil
}

// HasCapability checks if the assistant has a specific capability
func (a *Assistant) HasCapability(capability AssistantCapability) bool {
	for _, cap := range a.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}
