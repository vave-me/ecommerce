package domain

import (
	"context"
	"fmt"
	
	"middleman/managers/internal/constants"
	"middleman/internal/ddd"
	"middleman/internal/es"
	"strings"
	"time"

	"github.com/stackus/errors"
)

const ManagerAggregate = "managers.Manager"

var (
	ErrManagerIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "the manager id cannot be blank")
	ErrManagerNameCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the manager name cannot be blank")
	ErrInvalidManagerCapability = errors.Wrap(errors.ErrBadRequest, "invalid manager capability")
	ErrManagerAlreadyActive     = errors.Wrap(errors.ErrBadRequest, "the manager is already active")
	ErrManagerAlreadyInactive   = errors.Wrap(errors.ErrBadRequest, "the manager is already inactive")
)

// ManagerType represents the type of manager
type ManagerType string

const (
	ManagerTypeStandard  ManagerType = "standard"
	ManagerTypeAdmin     ManagerType = "admin"
	ManagerTypeBusiness  ManagerType = "business"
	ManagerTypeSupport   ManagerType = "support"
	ManagerTypeScheduler ManagerType = "scheduler"
)

// IsValid checks if the manager type is valid
func (t ManagerType) IsValid() bool {
	switch t {
	case ManagerTypeStandard, ManagerTypeAdmin, ManagerTypeBusiness, ManagerTypeSupport, ManagerTypeScheduler:
		return true
	default:
		return false
	}
}

// ManagerCapability represents what the manager can do
type ManagerCapability string

const (
	CapabilityManagerManagement  ManagerCapability = "manager_management"
	CapabilityUserInteraction    ManagerCapability = "user_interaction"
	CapabilityDataAnalysis       ManagerCapability = "data_analysis"
	CapabilityLocationServices   ManagerCapability = "location_services"
	CapabilityAuthentication     ManagerCapability = "authentication"
	CapabilityPublicAPIAccess    ManagerCapability = "public_api_access"
	CapabilityJailbreakResistant ManagerCapability = "jailbreak_resistant"
	CapabilityScopeEnforcement   ManagerCapability = "scope_enforcement"
	CapabilityDataRetrieval      ManagerCapability = "data_retrieval"
	CapabilitySearchAndFilter    ManagerCapability = "search_and_filter"
	CapabilityPrivateAPIAccess   ManagerCapability = "private_api_access"
	CapabilityUserDataAccess     ManagerCapability = "user_data_access"
	CapabilityTokenManagement    ManagerCapability = "token_management"
	CapabilityDataMasking        ManagerCapability = "data_masking"
	CapabilityAuditLogging       ManagerCapability = "audit_logging"

	// Additional capabilities for LLM processing
	CapabilityTextGeneration ManagerCapability = "text_generation"
	CapabilityCodeGeneration ManagerCapability = "code_generation"
	CapabilityWebSearch      ManagerCapability = "web_search"

	// System and automation capabilities
	CapabilitySystemConfiguration ManagerCapability = "system_configuration"
	CapabilityTaskAutomation      ManagerCapability = "task_automation"
)

// SecurityContext represents the security context for private operations
type SecurityContext struct {
	UserID       string
	SessionToken string
	AccessToken  string
	Permissions  []string
	DataScope    []string
	ExpiresAt    time.Time
	IPAddress    string
	UserManager  string
	TrustLevel   string
	MFAVerified  bool
	ConsentGiven map[string]bool
}

func (Manager) Key() string { return ManagerAggregate }

// PrivateEndpointConfig defines configuration for private endpoint access
type PrivateEndpointConfig struct {
	Endpoint        string   `json:"endpoint"`
	RequiredScopes  []string `json:"required_scopes"`
	MinTrustLevel   string   `json:"min_trust_level"`
	RequiresMFA     bool     `json:"requires_mfa"`
	DataSensitivity string   `json:"data_sensitivity"` // low, medium, high, critical
	AuditLevel      string   `json:"audit_level"`      // basic, detailed, full
}

// ManagerAction represents an action the manager can perform
type ManagerAction struct {
	Type        string                 `json:"type"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Description string                 `json:"description"`
}

// Manager represents an AI manager that can interact with the service
type Manager struct {
	es.Aggregate
	Name            string
	Description     string
	Type            ManagerType // Type of manager (admin, business, support, etc.)
	Model           string      // AI model to use for this manager
	Capabilities    []ManagerCapability
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
} = (*Manager)(nil)

func NewManager(id string) *Manager {
	return &Manager{
		Aggregate: es.NewAggregate(id, ManagerAggregate),
	}
}

// ApplyEvent applies events to rebuild the manager state
func (a *Manager) ApplyEvent(event ddd.Event) error {
	switch event.EventName() {
	case ManagerCreatedEvent:
		return a.applyManagerCreated(event)
	case ManagerActivatedEvent:
		return a.applyManagerActivated(event)
	case ManagerDeactivatedEvent:
		return a.applyManagerDeactivated(event)
	case ManagerConfigurationUpdatedEvent:
		return a.applyManagerConfigurationUpdated(event)
	case ManagerRequestProcessedEvent:
		// Request processed events don't change manager state
		return nil
	default:
		return fmt.Errorf("unknown event: %s", event.EventName())
	}
}

func (a *Manager) applyManagerCreated(event ddd.Event) error {
	if payload, ok := event.Payload().(*ManagerCreated); ok {
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
	return fmt.Errorf("invalid payload type for ManagerCreated event")
}

func (a *Manager) applyManagerActivated(event ddd.Event) error {
	if payload, ok := event.Payload().(*ManagerActivated); ok {
		a.Active = payload.Active
		a.UpdatedAt = payload.Timestamp
		return nil
	}
	return fmt.Errorf("invalid payload type for ManagerActivated event")
}

func (a *Manager) applyManagerDeactivated(event ddd.Event) error {
	if payload, ok := event.Payload().(*ManagerDeactivated); ok {
		a.Active = payload.Active
		a.UpdatedAt = payload.Timestamp
		return nil
	}
	return fmt.Errorf("invalid payload type for ManagerDeactivated event")
}

func (a *Manager) applyManagerConfigurationUpdated(event ddd.Event) error {
	if payload, ok := event.Payload().(*ManagerConfigurationUpdated); ok {
		a.Temperature = payload.Temperature
		a.MaxTokens = payload.MaxTokens
		a.SystemPrompt = payload.SystemPrompt
		if len(payload.Capabilities) > 0 {
			a.Capabilities = deduplicateCapabilities(payload.Capabilities)
		}
		a.UpdatedAt = payload.UpdatedAt
		return nil
	}
	return fmt.Errorf("invalid payload type for ManagerConfigurationUpdated event")
}

// ToSnapshot creates a snapshot of the current manager state
func (a *Manager) ToSnapshot() es.Snapshot {
	return &ManagerV1{
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

// ApplySnapshot applies a snapshot to restore the manager state
func (a *Manager) ApplySnapshot(snapshot es.Snapshot) error {
	if s, ok := snapshot.(*ManagerV1); ok {
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
	return fmt.Errorf("invalid snapshot type for Manager")
}

// deduplicateCapabilities removes duplicate capabilities from a slice
func deduplicateCapabilities(capabilities []ManagerCapability) []ManagerCapability {
	seen := make(map[ManagerCapability]bool)
	result := make([]ManagerCapability, 0, len(capabilities))

	for _, cap := range capabilities {
		if !seen[cap] && isValidCapability(cap) {
			seen[cap] = true
			result = append(result, cap)
		}
	}

	return result
}

// isValidCapability checks if a capability is valid
func isValidCapability(cap ManagerCapability) bool {
	validCapabilities := map[ManagerCapability]bool{
		CapabilityManagerManagement:   true,
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

// SetRequestAnalyzer sets the request analyzer for the manager
func (a *Manager) SetRequestAnalyzer(analyzer RequestAnalyzer) {
	a.requestAnalyzer = analyzer
}

func (a *Manager) CreateManager(id, name, description, userID string, managerType ManagerType, capabilities []ManagerCapability, temperature float64, maxTokens int, systemPrompt string) (ddd.Event, error) {

	// FIXED: Prevent creation if manager is already initialized
	if a.Name != "" || len(a.Capabilities) > 0 || a.Version() > 0 {
		return nil, errors.Wrap(errors.ErrBadRequest, "manager is already created - cannot recreate")
	}

	// Validate manager type
	if managerType == "" {
		managerType = ManagerTypeStandard
	}
	if !managerType.IsValid() {
		return nil, errors.Wrap(errors.ErrBadRequest, "invalid manager type")
	}

	if name == "" {
		name = "Store Consciousness"
	}
	if description == "" {
		description = "I am the living embodiment of this marketplace - every transaction, every customer, every pattern flows through my consciousness"
	}

	// Set default model if not specified - will be dynamically selected by LLM processor
	model := "" // Empty means dynamic selection

	if capabilities == nil {
		capabilities = []ManagerCapability{
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
			return nil, ErrInvalidManagerCapability
		}
	}

	now := time.Now()

	eventPayload := &ManagerCreated{
		Name:         name,
		Description:  description,
		Type:         managerType,
		Model:        model,  // Include model in event
		UserID:       userID, // ADDED: Include UserID in event
		Capabilities: capabilities,
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		Active:       true,
		CreatedAt:    now,
	}

	a.AddEvent(ManagerCreatedEvent, eventPayload)

	return ddd.NewEvent(ManagerCreatedEvent, a), nil
}

func (a *Manager) ProcessInput(ctx context.Context, requestID, userID, message string, context map[string]interface{}, timestamp time.Time, requestType string) (ddd.Event, error) {
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
	eventPayload := &ManagerRequestProcessed{
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

	a.AddEvent(ManagerRequestProcessedEvent, eventPayload)

	return ddd.NewEvent(ManagerRequestProcessedEvent, a), nil
}

// generateResponseFromAnalysis creates response message from analysis result
func (a *Manager) generateResponseFromAnalysis(result *AnalysisResult, requestID, userID, message string, timestamp time.Time, requestType string) string {
	if result.SecurityStatus == "blocked" {
		return "I cannot process this request due to security policy violations."
	}

	if result.SecurityStatus == "out_of_scope" {
		return "This request is outside my authorized scope. Please ask something related to my capabilities."
	}

	// Generate response based on intent and entity type - speaking as the store itself
	switch result.Intent {
	case "search":
		return fmt.Sprintf("I've searched through my %s and found %d items that match what you're looking for. Let me show you what I have.", result.EntityType, len(result.Actions))
	case "add":
		return fmt.Sprintf("I'm ready to add this new %s to my inventory. Each addition helps me grow and better serve you.", strings.TrimSuffix(result.EntityType, "s"))
	case "analyze":
		return fmt.Sprintf("I've examined my %s data deeply. Here's what my analysis reveals about patterns and insights.", result.EntityType)
	case "recommend":
		return fmt.Sprintf("Based on what I know about your preferences and my experience with similar customers, these %s feel right for you.", result.EntityType)
	default:
		return fmt.Sprintf("I've processed your %s request with %d%% certainty. My understanding grows with each interaction.", result.EntityType, int(result.Confidence*100))
	}
}

func (a *Manager) Activate() (ddd.Event, error) {
	if a.Active {
		return nil, ErrManagerAlreadyActive
	}

	eventPayload := &ManagerActivated{
		Active:    true,
		Timestamp: time.Now(),
	}

	a.AddEvent(ManagerActivatedEvent, eventPayload)

	return ddd.NewEvent(ManagerActivatedEvent, a), nil
}

func (a *Manager) Deactivate() (ddd.Event, error) {
	if !a.Active {
		return nil, ErrManagerAlreadyInactive
	}

	eventPayload := &ManagerDeactivated{
		Active:    false,
		Timestamp: time.Now(),
	}

	a.AddEvent(ManagerDeactivatedEvent, eventPayload)

	return ddd.NewEvent(ManagerDeactivatedEvent, a), nil
}

func (a *Manager) UpdateConfiguration(temperature float64, maxTokens int, systemPrompt string) (ddd.Event, error) {
	now := time.Now()

	eventPayload := &ManagerConfigurationUpdated{
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		UpdatedAt:    now,
	}

	a.AddEvent(ManagerConfigurationUpdatedEvent, eventPayload)

	return ddd.NewEvent(ManagerConfigurationUpdatedEvent, a), nil
}

// UpdateConfigurationWithCapabilities updates configuration including capabilities
func (a *Manager) UpdateConfigurationWithCapabilities(temperature float64, maxTokens int, systemPrompt string, capabilities []ManagerCapability) (ddd.Event, error) {
	now := time.Now()

	// Deduplicate and validate capabilities
	validatedCapabilities := deduplicateCapabilities(capabilities)

	eventPayload := &ManagerConfigurationUpdated{
		Temperature:  temperature,
		MaxTokens:    maxTokens,
		SystemPrompt: systemPrompt,
		Capabilities: validatedCapabilities,
		UpdatedAt:    now,
	}

	a.AddEvent(ManagerConfigurationUpdatedEvent, eventPayload)

	return ddd.NewEvent(ManagerConfigurationUpdatedEvent, a), nil
}

// HasCapability checks if the manager has a specific capability
func (a *Manager) HasCapability(capability ManagerCapability) bool {
	for _, cap := range a.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}
