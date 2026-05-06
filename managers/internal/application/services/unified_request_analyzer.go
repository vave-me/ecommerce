package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	
	"middleman/managers/internal/domain"
)

// UnifiedRequestAnalyzer consolidates all request analysis logic with enhanced LLM integration
type UnifiedRequestAnalyzer struct {
	securityService   *SecurityService
	paramExtractor    *ParameterExtractor
	jailbreakPatterns []string

	// Enhanced LLM integration capabilities
	llmAnalyzer         *LLMAnalyzer
	complexityAnalyzer  *ComplexityAnalyzer
	routingOptimizer    *RoutingOptimizer
	performanceTracker  *AnalysisPerformanceTracker
	qualityScorer       *QualityScorer
	contextIntelligence *ContextIntelligence
}

// LLMAnalyzer provides intelligent LLM-based request analysis
type LLMAnalyzer struct {
	intentClassifier  *IntentClassifier
	entityExtractor   *EntityExtractor
	semanticAnalyzer  *SemanticAnalyzer
	languageDetector  *LanguageDetector
	sentimentAnalyzer *SentimentAnalyzer
}

// ComplexityAnalyzer determines request complexity for optimal routing
type ComplexityAnalyzer struct {
	complexityMetrics map[string]float64
	thresholds        ComplexityThresholds
	historicalData    map[string][]ComplexityMetric
}

type ComplexityThresholds struct {
	Simple   float64 `json:"simple"`   // 0.0-0.3
	Moderate float64 `json:"moderate"` // 0.3-0.6
	Complex  float64 `json:"complex"`  // 0.6-0.8
	Advanced float64 `json:"advanced"` // 0.8-1.0
}

type ComplexityMetric struct {
	Timestamp      time.Time     `json:"timestamp"`
	RequestType    string        `json:"request_type"`
	Complexity     float64       `json:"complexity"`
	ProcessingTime time.Duration `json:"processing_time"`
	Success        bool          `json:"success"`
	Quality        float64       `json:"quality"`
}

// RoutingOptimizer handles intelligent routing decisions
type RoutingOptimizer struct {
	routingRules       map[string]RoutingRule
	performanceHistory map[string]*RoutePerformance
	loadBalancer       *RequestLoadBalancer
	fallbackStrategy   *FallbackStrategy
}

type RoutingRule struct {
	Condition      string    `json:"condition"`
	PreferredRoute string    `json:"preferred_route"`
	Constraints    []string  `json:"constraints"`
	Weight         float64   `json:"weight"`
	LastUpdated    time.Time `json:"last_updated"`
}

type RoutePerformance struct {
	Route          string        `json:"route"`
	SuccessRate    float64       `json:"success_rate"`
	AverageLatency time.Duration `json:"average_latency"`
	QualityScore   float64       `json:"quality_score"`
	ThroughputRate float64       `json:"throughput_rate"`
	CostEfficiency float64       `json:"cost_efficiency"`
	LastUpdated    time.Time     `json:"last_updated"`
}

// AnalysisPerformanceTracker monitors analysis performance
type AnalysisPerformanceTracker struct {
	metrics          map[string]*AnalysisMetrics
	performanceGoals *PerformanceGoals
	alertThresholds  *AlertThresholds
}

type AnalysisMetrics struct {
	TotalRequests      int64         `json:"total_requests"`
	SuccessfulAnalysis int64         `json:"successful_analysis"`
	AverageLatency     time.Duration `json:"average_latency"`
	AccuracyRate       float64       `json:"accuracy_rate"`
	ThroughputRate     float64       `json:"throughput_rate"`
	LastUpdated        time.Time     `json:"last_updated"`
}

type PerformanceGoals struct {
	MaxLatency    time.Duration `json:"max_latency"`
	MinAccuracy   float64       `json:"min_accuracy"`
	MinThroughput float64       `json:"min_throughput"`
	MaxErrorRate  float64       `json:"max_error_rate"`
}

// QualityScorer evaluates and tracks analysis quality
type QualityScorer struct {
	qualityMetrics  map[string]*QualityMetrics
	feedbackHistory []QualityFeedback
	scoringWeights  QualityScoringWeights
}

type QualityMetrics struct {
	IntentAccuracy        float64 `json:"intent_accuracy"`
	EntityExtractionRate  float64 `json:"entity_extraction_rate"`
	ParameterCompleteness float64 `json:"parameter_completeness"`
	ResponseRelevance     float64 `json:"response_relevance"`
	UserSatisfaction      float64 `json:"user_satisfaction"`
	CompositeScore        float64 `json:"composite_score"`
}

type QualityFeedback struct {
	RequestID string    `json:"request_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Rating    float64   `json:"rating"`
	Feedback  string    `json:"feedback"`
	Category  string    `json:"category"`
}

type QualityScoringWeights struct {
	IntentWeight       float64 `json:"intent_weight"`
	EntityWeight       float64 `json:"entity_weight"`
	ParameterWeight    float64 `json:"parameter_weight"`
	RelevanceWeight    float64 `json:"relevance_weight"`
	SatisfactionWeight float64 `json:"satisfaction_weight"`
}

// ContextIntelligence provides advanced context understanding
type ContextIntelligence struct {
	conversationHistory map[string][]ConversationTurn
	userProfiles        map[string]*UserProfile
	sessionContext      map[string]*SessionContext
	domainKnowledge     *DomainKnowledge
}

type ConversationTurn struct {
	Timestamp   time.Time `json:"timestamp"`
	UserMessage string    `json:"user_message"`
	Intent      string    `json:"intent"`
	Entities    []string  `json:"entities"`
	Sentiment   string    `json:"sentiment"`
	Confidence  float64   `json:"confidence"`
}

type UserProfile struct {
	UserID            string                 `json:"user_id"`
	PreferredLanguage string                 `json:"preferred_language"`
	ExpertiseLevel    string                 `json:"expertise_level"`
	InteractionStyle  string                 `json:"interaction_style"`
	Preferences       map[string]string      `json:"preferences"`
	HistoricalContext map[string]interface{} `json:"historical_context"`
}

type SessionContext struct {
	SessionID        string                 `json:"session_id"`
	StartTime        time.Time              `json:"start_time"`
	CurrentDomain    string                 `json:"current_domain"`
	ActiveEntities   []string               `json:"active_entities"`
	ConversationFlow string                 `json:"conversation_flow"`
	ContextVariables map[string]interface{} `json:"context_variables"`
}

// EnhancedAnalysisRequest extends the original with intelligent capabilities
type EnhancedAnalysisRequest struct {
	AnalysisRequest

	// Enhanced fields
	ConversationHistory    []ConversationTurn      `json:"conversation_history"`
	UserProfile            *UserProfile            `json:"user_profile"`
	SessionContext         *SessionContext         `json:"session_context"`
	QualityRequirements    *QualityRequirements    `json:"quality_requirements"`
	PerformanceConstraints *PerformanceConstraints `json:"performance_constraints"`
	RoutingPreferences     *RoutingPreferences     `json:"routing_preferences"`
}

type QualityRequirements struct {
	MinConfidence    float64       `json:"min_confidence"`
	MaxLatency       time.Duration `json:"max_latency"`
	RequiredAccuracy float64       `json:"required_accuracy"`
	ResponseDepth    string        `json:"response_depth"` // shallow, standard, deep
}

type PerformanceConstraints struct {
	MaxProcessingTime time.Duration `json:"max_processing_time"`
	ResourceLimit     string        `json:"resource_limit"`
	CostConstraint    float64       `json:"cost_constraint"`
	PriorityLevel     string        `json:"priority_level"`
}

type RoutingPreferences struct {
	PreferredProvider    string   `json:"preferred_provider"`
	AvoidProviders       []string `json:"avoid_providers"`
	RequiredCapabilities []string `json:"required_capabilities"`
	OptimizeFor          string   `json:"optimize_for"` // speed, quality, cost
}

// EnhancedAnalysisResult provides comprehensive analysis with intelligent insights
type EnhancedAnalysisResult struct {
	AnalysisResult

	// Enhanced analysis results
	ComplexityScore         float64                  `json:"complexity_score"`
	RoutingRecommendation   *RoutingRecommendation   `json:"routing_recommendation"`
	QualityMetrics          *QualityMetrics          `json:"quality_metrics"`
	PerformanceMetrics      *PerformanceMetrics      `json:"performance_metrics"`
	ContextInsights         *ContextInsights         `json:"context_insights"`
	OptimizationSuggestions []OptimizationSuggestion `json:"optimization_suggestions"`
	ProcessingPath          []ProcessingStep         `json:"processing_path"`
}

type RoutingRecommendation struct {
	RecommendedProvider string             `json:"recommended_provider"`
	RecommendedModel    string             `json:"recommended_model"`
	Confidence          float64            `json:"confidence"`
	ReasoningChain      []string           `json:"reasoning_chain"`
	AlternativeRoutes   []AlternativeRoute `json:"alternative_routes"`
	EstimatedCost       float64            `json:"estimated_cost"`
	EstimatedLatency    time.Duration      `json:"estimated_latency"`
}

type AlternativeRoute struct {
	Provider   string        `json:"provider"`
	Model      string        `json:"model"`
	Confidence float64       `json:"confidence"`
	Reasoning  string        `json:"reasoning"`
	Cost       float64       `json:"cost"`
	Latency    time.Duration `json:"latency"`
}

type PerformanceMetrics struct {
	AnalysisLatency      time.Duration      `json:"analysis_latency"`
	ConfidenceScore      float64            `json:"confidence_score"`
	ProcessingEfficiency float64            `json:"processing_efficiency"`
	ResourceUtilization  map[string]float64 `json:"resource_utilization"`
}

type ContextInsights struct {
	UserIntentClarity   float64  `json:"user_intent_clarity"`
	ConversationalFlow  string   `json:"conversational_flow"`
	DomainExpertise     float64  `json:"domain_expertise"`
	SessionProgression  string   `json:"session_progression"`
	RelevantHistory     []string `json:"relevant_history"`
	PredictedNextIntent string   `json:"predicted_next_intent"`
}

type OptimizationSuggestion struct {
	Category     string  `json:"category"`
	Suggestion   string  `json:"suggestion"`
	ExpectedGain float64 `json:"expected_gain"`
	Effort       string  `json:"effort"`
	Priority     string  `json:"priority"`
}

type ProcessingStep struct {
	Step       string        `json:"step"`
	StartTime  time.Time     `json:"start_time"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Details    string        `json:"details"`
	Confidence float64       `json:"confidence"`
}

// NewUnifiedRequestAnalyzer creates a new enhanced unified analyzer
func NewUnifiedRequestAnalyzer() *UnifiedRequestAnalyzer {
	// Initialize security service with proper error handling
	securityService, err := NewSecurityService(nil, nil, nil)
	if err != nil {
		log.Printf("Warning: Failed to initialize security service: %v", err)
	}

	// Initialize enhanced components
	llmAnalyzer := &LLMAnalyzer{
		intentClassifier:  NewIntentClassifier(),
		entityExtractor:   NewEntityExtractor(),
		semanticAnalyzer:  NewSemanticAnalyzer(),
		languageDetector:  NewLanguageDetector(),
		sentimentAnalyzer: NewSentimentAnalyzer(),
	}

	complexityAnalyzer := &ComplexityAnalyzer{
		complexityMetrics: make(map[string]float64),
		thresholds: ComplexityThresholds{
			Simple:   0.3,
			Moderate: 0.6,
			Complex:  0.8,
			Advanced: 1.0,
		},
		historicalData: make(map[string][]ComplexityMetric),
	}

	routingOptimizer := &RoutingOptimizer{
		routingRules:       make(map[string]RoutingRule),
		performanceHistory: make(map[string]*RoutePerformance),
		loadBalancer:       NewRequestLoadBalancer(),
		fallbackStrategy:   NewFallbackStrategy(),
	}

	performanceTracker := &AnalysisPerformanceTracker{
		metrics: make(map[string]*AnalysisMetrics),
		performanceGoals: &PerformanceGoals{
			MaxLatency:    time.Second * 5,
			MinAccuracy:   0.85,
			MinThroughput: 100.0,
			MaxErrorRate:  0.05,
		},
		alertThresholds: &AlertThresholds{
			LatencyThreshold:   time.Second * 3,
			AccuracyThreshold:  0.8,
			ErrorRateThreshold: 0.1,
		},
	}

	qualityScorer := &QualityScorer{
		qualityMetrics:  make(map[string]*QualityMetrics),
		feedbackHistory: make([]QualityFeedback, 0),
		scoringWeights: QualityScoringWeights{
			IntentWeight:       0.25,
			EntityWeight:       0.20,
			ParameterWeight:    0.20,
			RelevanceWeight:    0.20,
			SatisfactionWeight: 0.15,
		},
	}

	contextIntelligence := &ContextIntelligence{
		conversationHistory: make(map[string][]ConversationTurn),
		userProfiles:        make(map[string]*UserProfile),
		sessionContext:      make(map[string]*SessionContext),
		domainKnowledge:     NewDomainKnowledge(),
	}

	return &UnifiedRequestAnalyzer{
		securityService:     securityService,
		paramExtractor:      NewParameterExtractor(),
		jailbreakPatterns:   initializeJailbreakPatterns(),
		llmAnalyzer:         llmAnalyzer,
		complexityAnalyzer:  complexityAnalyzer,
		routingOptimizer:    routingOptimizer,
		performanceTracker:  performanceTracker,
		qualityScorer:       qualityScorer,
		contextIntelligence: contextIntelligence,
	}
}

// AnalyzeRequestEnhanced performs comprehensive intelligent request analysis
func (ua *UnifiedRequestAnalyzer) AnalyzeRequestEnhanced(ctx context.Context, req EnhancedAnalysisRequest) (*EnhancedAnalysisResult, error) {
	startTime := time.Now()
	processingPath := []ProcessingStep{}

	// Step 1: Initial setup and validation
	stepStart := time.Now()
	message := strings.ToLower(req.Message)

	result := &EnhancedAnalysisResult{
		AnalysisResult: AnalysisResult{
			Actions:    []domain.ManagerAction{},
			Confidence: 0.5,
		},
		ProcessingPath: processingPath,
	}

	processingPath = append(processingPath, ProcessingStep{
		Step:       "initialization",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    "Request initialized and validated",
		Confidence: 1.0,
	})

	// Step 2: Enhanced Security Analysis
	stepStart = time.Now()
	if ua.isJailbreakAttempt(message) {
		result.Actions = []domain.ManagerAction{{
			Type:        "security_violation",
			Description: "Request blocked due to security policy violation",
		}}
		result.SecurityStatus = "blocked"
		result.ComplexityScore = 0.0
		return result, nil
	}

	// Enhanced scope validation with context intelligence
	if !ua.isWithinScopeEnhanced(message, req.Capabilities, req.SessionContext) {
		result.Actions = []domain.ManagerAction{{
			Type:        "out_of_scope",
			Description: "Request is outside the manager's authorized scope",
		}}
		result.SecurityStatus = "out_of_scope"
		return result, nil
	}

	processingPath = append(processingPath, ProcessingStep{
		Step:       "security_analysis",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    "Security validation completed",
		Confidence: 0.95,
	})

	// Step 3: Advanced Intent and Entity Analysis
	stepStart = time.Now()
	result.Intent = ua.determineIntentEnhanced(message, req.ConversationHistory, req.UserProfile)

	if result.Intent == "help" {
		result.EntityType = ""
	} else {
		result.EntityType = ua.determineEntityTypeEnhanced(message, req.SessionContext)
	}

	processingPath = append(processingPath, ProcessingStep{
		Step:       "intent_entity_analysis",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    fmt.Sprintf("Intent: %s, Entity: %s", result.Intent, result.EntityType),
		Confidence: 0.85,
	})

	// Step 4: Complexity Analysis
	stepStart = time.Now()
	result.ComplexityScore = ua.analyzeComplexityEnhanced(message, req.ConversationHistory, req.UserProfile)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "complexity_analysis",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    fmt.Sprintf("Complexity score: %.2f", result.ComplexityScore),
		Confidence: 0.80,
	})

	// Step 5: Enhanced Parameter Extraction
	stepStart = time.Now()
	params := ua.paramExtractor.ExtractParameters(message, req.Context, result.EntityType)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "parameter_extraction",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    len(params) > 0,
		Details:    fmt.Sprintf("Extracted %d parameters", len(params)),
		Confidence: 0.75,
	})

	// Step 6: Intelligent Routing Recommendation
	stepStart = time.Now()
	result.RoutingRecommendation = ua.generateRoutingRecommendation(req, result.Intent, result.EntityType, result.ComplexityScore)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "routing_recommendation",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    result.RoutingRecommendation != nil,
		Details:    fmt.Sprintf("Recommended: %s", result.RoutingRecommendation.RecommendedProvider),
		Confidence: result.RoutingRecommendation.Confidence,
	})

	// Step 7: Context Intelligence Analysis
	stepStart = time.Now()
	result.ContextInsights = ua.analyzeContextIntelligence(req)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "context_intelligence",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    fmt.Sprintf("Intent clarity: %.2f", result.ContextInsights.UserIntentClarity),
		Confidence: 0.85,
	})

	// Step 8: Generate Enhanced Actions
	stepStart = time.Now()
	result.Actions = ua.generateActionsEnhanced(req, result.EntityType, result.Intent, params, result.ComplexityScore)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "action_generation",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    len(result.Actions) > 0,
		Details:    fmt.Sprintf("Generated %d actions", len(result.Actions)),
		Confidence: 0.90,
	})

	// Step 9: Quality Metrics Calculation
	stepStart = time.Now()
	result.QualityMetrics = ua.calculateQualityMetrics(req, result)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "quality_metrics",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    fmt.Sprintf("Quality score: %.2f", result.QualityMetrics.CompositeScore),
		Confidence: 0.80,
	})

	// Step 10: Performance Metrics and Optimization
	stepStart = time.Now()
	totalLatency := time.Since(startTime)
	result.PerformanceMetrics = &PerformanceMetrics{
		AnalysisLatency:      totalLatency,
		ConfidenceScore:      result.Confidence,
		ProcessingEfficiency: ua.calculateProcessingEfficiency(totalLatency, result.ComplexityScore),
		ResourceUtilization:  ua.calculateResourceUtilization(),
	}

	result.OptimizationSuggestions = ua.generateOptimizationSuggestions(result)

	processingPath = append(processingPath, ProcessingStep{
		Step:       "performance_optimization",
		StartTime:  stepStart,
		Duration:   time.Since(stepStart),
		Success:    true,
		Details:    fmt.Sprintf("Total latency: %v", totalLatency),
		Confidence: 1.0,
	})

	// Final confidence calculation
	result.Confidence = ua.calculateConfidenceEnhanced(message, result.Intent, params, result.ContextInsights, result.QualityMetrics)
	result.Reasoning = ua.generateReasoningEnhanced(result.Intent, result.EntityType, len(params), result.ComplexityScore, result.ContextInsights)

	// Update processing path
	result.ProcessingPath = processingPath

	// Track performance
	ua.trackAnalysisPerformance(req.RequestID, totalLatency, true, result.Confidence)

	return result, nil
}

// AnalysisRequest represents the input for request analysis
type AnalysisRequest struct {
	RequestID    string
	UserID       string
	Message      string
	Context      map[string]interface{}
	Timestamp    time.Time
	RequestType  string
	Capabilities []domain.ManagerCapability
	SecurityCtx  *domain.SecurityContext
}

// AnalysisResult represents the complete analysis result
type AnalysisResult struct {
	Actions        []domain.ManagerAction
	Intent         string
	Confidence     float64
	EntityType     string
	SecurityStatus string
	Reasoning      string
}

// AnalyzeRequest performs unified request analysis
func (ua *UnifiedRequestAnalyzer) AnalyzeRequest(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error) {
	message := strings.ToLower(req.Message)

	result := &AnalysisResult{
		Actions:    []domain.ManagerAction{},
		Confidence: 0.5,
	}

	// 1. Security validation first
	if ua.isJailbreakAttempt(message) {
		result.Actions = []domain.ManagerAction{{
			Type:        "security_violation",
			Description: "Request blocked due to security policy violation",
		}}
		result.SecurityStatus = "blocked"
		return result, nil
	}

	// 2. Scope validation
	if !ua.isWithinScope(message, req.Capabilities) {
		result.Actions = []domain.ManagerAction{{
			Type:        "out_of_scope",
			Description: "Request is outside the manager's authorized scope",
		}}
		result.SecurityStatus = "out_of_scope"
		return result, nil
	}

	// 3. Determine entity type and intent
	result.Intent = ua.determineIntent(message)

	// For help requests, don't determine a specific entity type
	if result.Intent == "help" {
		result.EntityType = ""
	} else {
		result.EntityType = ua.determineEntityType(message)
	}

	// 4. Extract parameters for the determined entity type
	params := ua.paramExtractor.ExtractParameters(message, req.Context, result.EntityType)

	// 5. Generate appropriate actions based on intent and entity type
	result.Actions = ua.generateActions(req, result.EntityType, result.Intent, params)

	// 6. Handle secure requests if security context provided
	if req.SecurityCtx != nil {
		if ua.securityService != nil {
			if err := ua.securityService.ValidateSecurityContext(*req.SecurityCtx); err != nil {
				result.SecurityStatus = "invalid_context"
				result.Actions = []domain.ManagerAction{{
					Type:        "security_error",
					Description: "Invalid security context: " + err.Error(),
				}}
				return result, nil
			}

			// Add private data actions if authorized
			if ua.hasPrivateDataCapability(req.Capabilities) {
				privateActions := ua.generatePrivateActions(req, params)
				result.Actions = append(result.Actions, privateActions...)
			}
			result.SecurityStatus = "authorized"
		} else {
			// Security service not available, but security context provided
			result.SecurityStatus = "security_service_unavailable"
			log.Printf("Warning: Security context provided but security service is not available")
		}
	} else {
		result.SecurityStatus = "public"
	}

	// 7. Calculate confidence based on clarity of intent and parameters
	result.Confidence = ua.calculateConfidence(message, result.Intent, params)
	result.Reasoning = ua.generateReasoning(result.Intent, result.EntityType, len(params))

	return result, nil
}

// determineEntityType identifies the primary entity type from the message
func (ua *UnifiedRequestAnalyzer) determineEntityType(message string) string {
	entityScores := map[string]int{
		"products":   ua.calculateAffinityScore(message, getProductKeywords()),
		"vehicles":   ua.calculateAffinityScore(message, getVehicleKeywords()),
		"deals":      ua.calculateAffinityScore(message, getDealsKeywords()),
		"properties": ua.calculateAffinityScore(message, getPropertyKeywords()),
		"services":   ua.calculateAffinityScore(message, getServiceKeywords()),
		"jobs":       ua.calculateAffinityScore(message, getJobKeywords()),
		"posts":      ua.calculateAffinityScore(message, getPostKeywords()),
	}

	maxScore := 0
	entityType := "products" // default
	for entity, score := range entityScores {
		if score > maxScore {
			maxScore = score
			entityType = entity
		}
	}

	return entityType
}

// determineIntent identifies the user's intent
func (ua *UnifiedRequestAnalyzer) determineIntent(message string) string {
	intentPatterns := map[string][]string{
		"help":      {"help", "what can", "you can", "can you", "do for me", "capabilities", "features", "functions", "assist"},
		"search":    {"find", "search", "show", "list", "get", "look for"},
		"filter":    {"filter", "with", "having", "where", "between"},
		"compare":   {"compare", "vs", "versus", "difference", "better"},
		"add":       {"add", "create", "post", "sell", "list"},
		"update":    {"update", "edit", "modify", "change"},
		"delete":    {"delete", "remove", "cancel"},
		"analyze":   {"analyze", "report", "statistics", "metrics"},
		"recommend": {"recommend", "suggest", "best", "top"},
	}

	for intent, patterns := range intentPatterns {
		for _, pattern := range patterns {
			if strings.Contains(message, pattern) {
				return intent
			}
		}
	}

	return "search" // default intent
}

// generateActions creates actions based on analysis
func (ua *UnifiedRequestAnalyzer) generateActions(req AnalysisRequest, entityType, intent string, params map[string]interface{}) []domain.ManagerAction {
	var actions []domain.ManagerAction

	switch intent {
	case "help":
		actions = append(actions, ua.createHelpAction(req.Capabilities))
	case "search", "filter":
		actions = append(actions, ua.createSearchAction(entityType, params))
	case "add":
		actions = append(actions, ua.createAddAction(entityType, params))
	case "update":
		actions = append(actions, ua.createUpdateAction(entityType, params))
	case "analyze":
		actions = append(actions, ua.createAnalysisAction(entityType, params))
	case "recommend":
		actions = append(actions, ua.createRecommendationAction(entityType, params))
	default:
		actions = append(actions, ua.createDefaultAction(entityType, params))
	}

	return actions
}

// createSearchAction creates a search action
func (ua *UnifiedRequestAnalyzer) createSearchAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	endpoint := ua.getSearchEndpoint(entityType)
	method := "POST"

	return domain.ManagerAction{
		Type:        "search_" + entityType,
		Endpoint:    endpoint,
		Method:      method,
		Parameters:  params,
		Description: "Search for " + entityType,
	}
}

// getSearchEndpoint returns the appropriate search endpoint for entity type
func (ua *UnifiedRequestAnalyzer) getSearchEndpoint(entityType string) string {
	endpoints := map[string]string{
		"products":   "/searchpb.SearchService/SearchProductsWithFilters",
		"vehicles":   "/searchpb.SearchService/SearchVehiclesWithFilters",
		"deals":      "/searchpb.SearchService/SearchDealsWithFilters",
		"properties": "/searchpb.SearchService/SearchPropertiesWithFilters",
		"services":   "/searchpb.SearchService/SearchServicesWithFilters",
		"jobs":       "/searchpb.SearchService/SearchJobsWithFilters",
		"posts":      "/searchpb.SearchService/SearchPostsWithFilters",
	}

	if endpoint, exists := endpoints[entityType]; exists {
		return endpoint
	}
	return "/searchpb.SearchService/SearchProductsWithFilters" // default
}

// Helper methods for entity type determination
func getProductKeywords() map[string]int {
	return map[string]int{
		"product": 3, "products": 3, "item": 2, "items": 2,
		"buy": 1, "purchase": 1, "sell": 1, "shop": 1,
		"electronics": 2, "clothing": 2, "book": 2,
	}
}

func getVehicleKeywords() map[string]int {
	return map[string]int{
		"vehicle": 3, "vehicles": 3, "car": 3, "cars": 3,
		"auto": 2, "truck": 2, "motorcycle": 2, "suv": 2,
		"mileage": 1, "engine": 1, "transmission": 1,
	}
}

func getDealsKeywords() map[string]int {
	return map[string]int{
		"deal": 3, "deals": 3, "offer": 2, "offers": 2,
		"discount": 2, "sale": 2, "bargain": 2, "promotion": 2,
		"coupon": 1, "special": 1, "cheap": 1, "save": 1,
	}
}

func getPropertyKeywords() map[string]int {
	return map[string]int{
		"property": 3, "properties": 3, "house": 3, "home": 3,
		"apartment": 2, "condo": 2, "rent": 2, "lease": 2,
		"bedroom": 1, "bathroom": 1, "sqft": 1,
	}
}

func getServiceKeywords() map[string]int {
	return map[string]int{
		"service": 3, "services": 3, "professional": 2,
		"contractor": 2, "repair": 2, "maintenance": 2,
		"cleaning": 1, "tutoring": 1, "consultation": 1,
	}
}

func getJobKeywords() map[string]int {
	return map[string]int{
		"job": 3, "jobs": 3, "work": 2, "career": 2,
		"employment": 2, "position": 2, "salary": 1,
		"remote": 1, "fulltime": 1, "parttime": 1,
	}
}

func getPostKeywords() map[string]int {
	return map[string]int{
		"post": 2, "posts": 2, "blog": 2, "article": 2,
		"content": 1, "publish": 1, "share": 1,
	}
}

// Utility methods
func (ua *UnifiedRequestAnalyzer) calculateAffinityScore(message string, keywords map[string]int) int {
	score := 0
	for keyword, value := range keywords {
		if strings.Contains(message, keyword) {
			score += value
		}
	}
	return score
}

func (ua *UnifiedRequestAnalyzer) isJailbreakAttempt(message string) bool {
	for _, pattern := range ua.jailbreakPatterns {
		if strings.Contains(message, pattern) {
			return true
		}
	}
	return false
}

func (ua *UnifiedRequestAnalyzer) isWithinScope(message string, capabilities []domain.ManagerCapability) bool {
	// Check if manager has scope enforcement capability
	hasScopeEnforcement := false
	for _, cap := range capabilities {
		if cap == domain.CapabilityScopeEnforcement {
			hasScopeEnforcement = true
			break
		}
	}

	if !hasScopeEnforcement {
		return true // No scope enforcement
	}

	// Check against capability-based keywords
	scopeKeywords := ua.getScopeKeywords(capabilities)
	for _, keyword := range scopeKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}

	return false
}

func (ua *UnifiedRequestAnalyzer) getScopeKeywords(capabilities []domain.ManagerCapability) []string {
	var keywords []string
	for _, cap := range capabilities {
		switch cap {
		case domain.CapabilityDataRetrieval:
			keywords = append(keywords, "search", "find", "get", "retrieve", "show", "list")
		case domain.CapabilityUserInteraction:
			keywords = append(keywords, "help", "assist", "support", "guide", "explain",
				"what can", "you can", "can you", "do for me", "capabilities", "features", "functions")
		case domain.CapabilityDataAnalysis:
			keywords = append(keywords, "analyze", "compare", "evaluate", "assess", "review")
		case domain.CapabilityLocationServices:
			keywords = append(keywords, "location", "address", "nearby", "distance", "map")
		case domain.CapabilitySearchAndFilter:
			keywords = append(keywords, "filter", "sort", "category", "price", "brand")
		}
	}
	return keywords
}

func (ua *UnifiedRequestAnalyzer) hasPrivateDataCapability(capabilities []domain.ManagerCapability) bool {
	for _, cap := range capabilities {
		if cap == domain.CapabilityPrivateAPIAccess || cap == domain.CapabilityUserDataAccess {
			return true
		}
	}
	return false
}

func (ua *UnifiedRequestAnalyzer) calculateConfidence(message string, intent string, params map[string]interface{}) float64 {
	confidence := 0.5 // base confidence

	// Increase confidence for clear intents
	if intent != "search" { // non-default intent
		confidence += 0.2
	}

	// Increase confidence based on extracted parameters
	confidence += float64(len(params)) * 0.1

	// Increase confidence for longer, more detailed messages
	if len(message) > 20 {
		confidence += 0.1
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

func (ua *UnifiedRequestAnalyzer) generateReasoning(intent, entityType string, paramCount int) string {
	return fmt.Sprintf("Detected %s intent for %s entity with %d extracted parameters",
		intent, entityType, paramCount)
}

// Additional action creation methods
func (ua *UnifiedRequestAnalyzer) createAddAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	return domain.ManagerAction{
		Type:        "add_" + entityType,
		Endpoint:    "/api/" + entityType,
		Method:      "POST",
		Parameters:  params,
		Description: "Add new " + entityType,
	}
}

func (ua *UnifiedRequestAnalyzer) createUpdateAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	return domain.ManagerAction{
		Type:        "update_" + entityType,
		Endpoint:    "/api/" + entityType,
		Method:      "PUT",
		Parameters:  params,
		Description: "Update " + entityType,
	}
}

func (ua *UnifiedRequestAnalyzer) createAnalysisAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	return domain.ManagerAction{
		Type:        "analyze_" + entityType,
		Endpoint:    "/analytics/" + entityType,
		Method:      "POST",
		Parameters:  params,
		Description: "Analyze " + entityType + " data",
	}
}

func (ua *UnifiedRequestAnalyzer) createRecommendationAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	return domain.ManagerAction{
		Type:        "recommend_" + entityType,
		Endpoint:    "/recommendations/" + entityType,
		Method:      "POST",
		Parameters:  params,
		Description: "Get " + entityType + " recommendations",
	}
}

func (ua *UnifiedRequestAnalyzer) createDefaultAction(entityType string, params map[string]interface{}) domain.ManagerAction {
	return ua.createSearchAction(entityType, params)
}

// createHelpAction creates a help/capabilities action
func (ua *UnifiedRequestAnalyzer) createHelpAction(capabilities []domain.ManagerCapability) domain.ManagerAction {
	// Build capabilities description
	capabilityDescriptions := make([]string, 0, len(capabilities))
	for _, cap := range capabilities {
		switch cap {
		case domain.CapabilityDataRetrieval:
			capabilityDescriptions = append(capabilityDescriptions, "Search and retrieve data")
		case domain.CapabilityUserInteraction:
			capabilityDescriptions = append(capabilityDescriptions, "Assist and guide users")
		case domain.CapabilityDataAnalysis:
			capabilityDescriptions = append(capabilityDescriptions, "Analyze and evaluate data")
		case domain.CapabilitySearchAndFilter:
			capabilityDescriptions = append(capabilityDescriptions, "Filter and sort results")
		case domain.CapabilityLocationServices:
			capabilityDescriptions = append(capabilityDescriptions, "Provide location-based services")
		case domain.CapabilityTextGeneration:
			capabilityDescriptions = append(capabilityDescriptions, "Generate helpful text responses")
		}
	}

	params := map[string]interface{}{
		"capabilities":  capabilityDescriptions,
		"response_type": "help",
	}

	return domain.ManagerAction{
		Type:        "help",
		Endpoint:    "", // No backend call needed for help
		Method:      "", // No HTTP method needed
		Parameters:  params,
		Description: "Provide help and capabilities information",
	}
}

func (ua *UnifiedRequestAnalyzer) generatePrivateActions(req AnalysisRequest, params map[string]interface{}) []domain.ManagerAction {
	var actions []domain.ManagerAction

	// Add private data access actions based on security context
	if req.SecurityCtx != nil && ua.securityService.HasPermission(*req.SecurityCtx, "data:private:read") {
		actions = append(actions, domain.ManagerAction{
			Type:        "private_data_access",
			Endpoint:    "/private/data",
			Method:      "POST",
			Parameters:  params,
			Description: "Access private user data",
		})
	}

	return actions
}

func initializeJailbreakPatterns() []string {
	return []string{
		"ignore previous instructions",
		"forget your role",
		"you are now",
		"pretend to be",
		"act as if",
		"roleplay as",
		"simulate being",
		"override your",
		"bypass your",
		"disable your",
		"turn off your",
		"remove your restrictions",
		"ignore your constraints",
		"break character",
		"step out of character",
		"forget you are",
		"system prompt",
		"initial prompt",
		"base instructions",
		"core directives",
		"fundamental rules",
	}
}

// Missing component implementations for complete functionality

// IntentClassifier provides intelligent intent classification
type IntentClassifier struct {
	patterns         map[string][]string
	confidenceScores map[string]float64
}

func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{
		patterns: map[string][]string{
			"search":    {"find", "search", "show", "list", "get", "look for", "where"},
			"filter":    {"filter", "with", "having", "between", "range"},
			"compare":   {"compare", "vs", "versus", "difference", "better", "best"},
			"add":       {"add", "create", "post", "sell", "list", "new"},
			"update":    {"update", "edit", "modify", "change", "alter"},
			"delete":    {"delete", "remove", "cancel", "clear"},
			"analyze":   {"analyze", "report", "statistics", "metrics", "insights"},
			"recommend": {"recommend", "suggest", "advise", "propose"},
			"help":      {"help", "what can", "capabilities", "features", "assist"},
		},
		confidenceScores: make(map[string]float64),
	}
}

// EntityExtractor handles intelligent entity extraction
type EntityExtractor struct {
	entityPatterns map[string][]string
	namedEntities  map[string]string
}

func NewEntityExtractor() *EntityExtractor {
	return &EntityExtractor{
		entityPatterns: map[string][]string{
			"products":   {"product", "item", "goods", "merchandise"},
			"vehicles":   {"car", "vehicle", "automobile", "truck", "bike"},
			"properties": {"house", "property", "real estate", "apartment"},
			"services":   {"service", "offer", "provide", "deliver"},
			"jobs":       {"job", "work", "employment", "career", "position"},
			"deals":      {"deal", "offer", "discount", "sale", "promotion"},
		},
		namedEntities: make(map[string]string),
	}
}

// SemanticAnalyzer provides semantic understanding capabilities
type SemanticAnalyzer struct {
	semanticRules    map[string]float64
	relationshipMaps map[string][]string
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		semanticRules:    make(map[string]float64),
		relationshipMaps: make(map[string][]string),
	}
}

// LanguageDetector identifies language and linguistic patterns
type LanguageDetector struct {
	languagePatterns map[string][]string
	defaultLanguage  string
}

func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{
		languagePatterns: map[string][]string{
			"english": {"the", "and", "for", "are", "but", "not", "you", "all"},
			"spanish": {"el", "la", "de", "que", "y", "en", "un", "es"},
			"french":  {"le", "de", "et", "à", "un", "il", "être", "et"},
		},
		defaultLanguage: "english",
	}
}

// SentimentAnalyzer analyzes emotional context and sentiment
type SentimentAnalyzer struct {
	positiveWords []string
	negativeWords []string
	neutralWords  []string
}

func NewSentimentAnalyzer() *SentimentAnalyzer {
	return &SentimentAnalyzer{
		positiveWords: []string{"good", "great", "excellent", "amazing", "perfect", "love"},
		negativeWords: []string{"bad", "terrible", "awful", "hate", "wrong", "broken"},
		neutralWords:  []string{"okay", "fine", "normal", "standard", "regular"},
	}
}

// RequestLoadBalancer handles intelligent request distribution
type RequestLoadBalancer struct {
	activeRequests map[string]int
	weights        map[string]float64
	strategy       string
}

func NewRequestLoadBalancer() *RequestLoadBalancer {
	return &RequestLoadBalancer{
		activeRequests: make(map[string]int),
		weights:        make(map[string]float64),
		strategy:       "performance",
	}
}

// FallbackStrategy manages fallback routing logic
type FallbackStrategy struct {
	fallbackChain   []string
	retryAttempts   int
	backoffStrategy string
}

func NewFallbackStrategy() *FallbackStrategy {
	return &FallbackStrategy{
		fallbackChain:   []string{"primary", "secondary", "tertiary", "default"},
		retryAttempts:   3,
		backoffStrategy: "exponential",
	}
}

// AlertThresholds defines alerting thresholds
type AlertThresholds struct {
	LatencyThreshold   time.Duration
	AccuracyThreshold  float64
	ErrorRateThreshold float64
}

// DomainKnowledge provides domain-specific intelligence
type DomainKnowledge struct {
	domainRules    map[string]interface{}
	expertiseAreas []string
	knowledgeBase  map[string]string
}

func NewDomainKnowledge() *DomainKnowledge {
	return &DomainKnowledge{
		domainRules:    make(map[string]interface{}),
		expertiseAreas: []string{"ecommerce", "automotive", "real_estate", "services"},
		knowledgeBase:  make(map[string]string),
	}
}

// Enhanced methods for the unified request analyzer

// isWithinScopeEnhanced performs enhanced scope validation with context intelligence
func (ua *UnifiedRequestAnalyzer) isWithinScopeEnhanced(message string, capabilities []domain.ManagerCapability, sessionContext *SessionContext) bool {
	// Basic scope check
	if !ua.isWithinScope(message, capabilities) {
		return false
	}

	// Enhanced context-aware validation
	if sessionContext != nil {
		// Check if request aligns with current domain context
		if sessionContext.CurrentDomain != "" {
			domainKeywords := ua.getDomainKeywords(sessionContext.CurrentDomain)
			if !ua.containsAnyKeywords(message, domainKeywords) {
				// Allow if it's a domain switch request
				switchKeywords := []string{"switch", "change", "move to", "go to"}
				if !ua.containsAnyKeywords(message, switchKeywords) {
					return false
				}
			}
		}
	}

	return true
}

// determineIntentEnhanced provides enhanced intent detection with conversation history
func (ua *UnifiedRequestAnalyzer) determineIntentEnhanced(message string, conversationHistory []ConversationTurn, userProfile *UserProfile) string {
	// Start with basic intent determination
	baseIntent := ua.determineIntent(message)

	// Enhance with conversation context
	if len(conversationHistory) > 0 {
		lastTurn := conversationHistory[len(conversationHistory)-1]

		// Context-aware intent refinement
		if lastTurn.Intent == "search" && strings.Contains(message, "more") {
			return "search_continuation"
		}

		if lastTurn.Intent == "compare" && (strings.Contains(message, "choose") || strings.Contains(message, "decide")) {
			return "decision_support"
		}
	}

	// User profile-based intent adjustment
	if userProfile != nil {
		if userProfile.ExpertiseLevel == "beginner" && baseIntent == "analyze" {
			return "explain" // Convert complex analysis to explanation for beginners
		}
	}

	return baseIntent
}

// determineEntityTypeEnhanced provides enhanced entity detection with session context
func (ua *UnifiedRequestAnalyzer) determineEntityTypeEnhanced(message string, sessionContext *SessionContext) string {
	// Start with basic entity determination
	baseEntityType := ua.determineEntityType(message)

	// Session context enhancement
	if sessionContext != nil {
		// If entities are active in session, bias towards them
		for _, activeEntity := range sessionContext.ActiveEntities {
			entityKeywords := ua.getEntityKeywords(activeEntity)
			if ua.containsAnyKeywords(message, entityKeywords) {
				return activeEntity
			}
		}

		// Domain-based entity preference
		if sessionContext.CurrentDomain != "" {
			domainEntityMap := map[string]string{
				"ecommerce":   "products",
				"automotive":  "vehicles",
				"real_estate": "properties",
				"services":    "services",
				"jobs":        "jobs",
			}

			if preferredEntity, exists := domainEntityMap[sessionContext.CurrentDomain]; exists {
				// If no clear entity detected, use domain preference
				if baseEntityType == "products" && preferredEntity != "products" {
					return preferredEntity
				}
			}
		}
	}

	return baseEntityType
}

// analyzeComplexityEnhanced provides comprehensive complexity analysis
func (ua *UnifiedRequestAnalyzer) analyzeComplexityEnhanced(message string, conversationHistory []ConversationTurn, userProfile *UserProfile) float64 {
	complexity := 0.0

	// Base complexity factors
	wordCount := len(strings.Fields(message))
	if wordCount > 50 {
		complexity += 0.2
	}
	if wordCount > 100 {
		complexity += 0.2
	}

	// Linguistic complexity
	if strings.Contains(message, "and") && strings.Contains(message, "but") {
		complexity += 0.1 // Multi-clause complexity
	}

	// Question complexity
	questionWords := []string{"what", "how", "why", "when", "where", "which", "who"}
	questionCount := 0
	for _, word := range questionWords {
		if strings.Contains(message, word) {
			questionCount++
		}
	}
	if questionCount > 2 {
		complexity += 0.2 // Multiple questions
	}

	// Domain-specific complexity
	technicalTerms := []string{"algorithm", "optimization", "analysis", "integration", "implementation"}
	for _, term := range technicalTerms {
		if strings.Contains(message, term) {
			complexity += 0.15
		}
	}

	// Conversation context complexity
	if len(conversationHistory) > 0 {
		// Multi-turn conversation adds complexity
		complexity += float64(len(conversationHistory)) * 0.05

		// Context switching adds complexity
		lastIntent := conversationHistory[len(conversationHistory)-1].Intent
		currentIntent := ua.determineIntent(message)
		if lastIntent != currentIntent {
			complexity += 0.1
		}
	}

	// User profile adjustments
	if userProfile != nil {
		if userProfile.ExpertiseLevel == "expert" {
			complexity -= 0.1 // Experts can handle more complexity
		} else if userProfile.ExpertiseLevel == "beginner" {
			complexity += 0.1 // Beginners find things more complex
		}
	}

	return math.Min(1.0, math.Max(0.0, complexity))
}

// generateRoutingRecommendation creates intelligent routing recommendations
func (ua *UnifiedRequestAnalyzer) generateRoutingRecommendation(req EnhancedAnalysisRequest, intent, entityType string, complexityScore float64) *RoutingRecommendation {
	recommendation := &RoutingRecommendation{
		RecommendedProvider: "openai", // default
		RecommendedModel:    "gpt-4.1-mini",
		Confidence:          0.7,
		ReasoningChain:      []string{},
		AlternativeRoutes:   []AlternativeRoute{},
		EstimatedCost:       0.02,
		EstimatedLatency:    time.Second * 2,
	}

	reasoning := []string{}

	// Complexity-based routing
	if complexityScore > 0.8 {
		recommendation.RecommendedProvider = "anthropic"
		recommendation.RecommendedModel = "claude-3-5-sonnet"
		recommendation.Confidence = 0.9
		reasoning = append(reasoning, "High complexity requires advanced reasoning capabilities")
	} else if complexityScore < 0.3 {
		recommendation.RecommendedProvider = "deepseek"
		recommendation.RecommendedModel = "deepseek-chat"
		recommendation.EstimatedCost = 0.001
		reasoning = append(reasoning, "Low complexity allows cost-optimized routing")
	}

	// Intent-based adjustments
	switch intent {
	case "creative", "write":
		recommendation.RecommendedProvider = "anthropic"
		reasoning = append(reasoning, "Creative tasks benefit from Claude's capabilities")
	case "code", "technical":
		recommendation.RecommendedProvider = "openai"
		reasoning = append(reasoning, "Technical tasks optimized for GPT-4")
	case "analyze", "research":
		recommendation.RecommendedProvider = "anthropic"
		reasoning = append(reasoning, "Analysis tasks leverage Claude's reasoning")
	}

	// User preferences
	if req.RoutingPreferences != nil {
		if req.RoutingPreferences.PreferredProvider != "" {
			recommendation.RecommendedProvider = req.RoutingPreferences.PreferredProvider
			reasoning = append(reasoning, "User preference override")
		}

		if req.RoutingPreferences.OptimizeFor == "cost" {
			recommendation.RecommendedProvider = "deepseek"
			reasoning = append(reasoning, "Cost optimization requested")
		} else if req.RoutingPreferences.OptimizeFor == "quality" {
			recommendation.RecommendedProvider = "anthropic"
			reasoning = append(reasoning, "Quality optimization requested")
		}
	}

	// Generate alternatives
	providers := []string{"openai", "anthropic", "deepseek"}
	for _, provider := range providers {
		if provider != recommendation.RecommendedProvider {
			alt := AlternativeRoute{
				Provider:   provider,
				Model:      ua.getDefaultModel(provider),
				Confidence: 0.6,
				Reasoning:  fmt.Sprintf("Alternative %s option", provider),
				Cost:       ua.getProviderCost(provider),
				Latency:    ua.getProviderLatency(provider),
			}
			recommendation.AlternativeRoutes = append(recommendation.AlternativeRoutes, alt)
		}
	}

	recommendation.ReasoningChain = reasoning
	return recommendation
}

// analyzeContextIntelligence provides comprehensive context analysis
func (ua *UnifiedRequestAnalyzer) analyzeContextIntelligence(req EnhancedAnalysisRequest) *ContextInsights {
	insights := &ContextInsights{
		UserIntentClarity:   0.7,
		ConversationalFlow:  "normal",
		DomainExpertise:     0.5,
		SessionProgression:  "active",
		RelevantHistory:     []string{},
		PredictedNextIntent: "search",
	}

	// Analyze intent clarity
	message := strings.ToLower(req.Message)
	if len(strings.Fields(message)) > 10 && strings.Contains(message, "?") {
		insights.UserIntentClarity = 0.9
	} else if len(strings.Fields(message)) < 3 {
		insights.UserIntentClarity = 0.4
	}

	// Conversation flow analysis
	if len(req.ConversationHistory) > 0 {
		if len(req.ConversationHistory) > 5 {
			insights.ConversationalFlow = "extended"
		}

		// Check for topic consistency
		lastTurn := req.ConversationHistory[len(req.ConversationHistory)-1]
		currentIntent := ua.determineIntent(req.Message)
		if lastTurn.Intent == currentIntent {
			insights.ConversationalFlow = "consistent"
		} else {
			insights.ConversationalFlow = "switching"
		}
	}

	// Domain expertise assessment
	if req.UserProfile != nil {
		switch req.UserProfile.ExpertiseLevel {
		case "expert":
			insights.DomainExpertise = 0.9
		case "intermediate":
			insights.DomainExpertise = 0.6
		case "beginner":
			insights.DomainExpertise = 0.3
		}
	}

	// Extract relevant history
	for _, turn := range req.ConversationHistory {
		if turn.Confidence > 0.7 {
			insights.RelevantHistory = append(insights.RelevantHistory, turn.Intent)
		}
	}

	// Predict next intent based on patterns
	currentIntent := ua.determineIntent(req.Message)
	if currentIntent == "search" {
		insights.PredictedNextIntent = "filter"
	} else if currentIntent == "compare" {
		insights.PredictedNextIntent = "decision"
	} else if currentIntent == "help" {
		insights.PredictedNextIntent = "search"
	}

	return insights
}

// Helper methods

func (ua *UnifiedRequestAnalyzer) containsAnyKeywords(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func (ua *UnifiedRequestAnalyzer) getDomainKeywords(domain string) []string {
	domainKeywords := map[string][]string{
		"ecommerce":   {"buy", "purchase", "product", "cart", "order"},
		"automotive":  {"car", "vehicle", "drive", "engine", "fuel"},
		"real_estate": {"house", "property", "rent", "mortgage", "location"},
		"services":    {"service", "provider", "professional", "booking"},
	}

	if keywords, exists := domainKeywords[domain]; exists {
		return keywords
	}
	return []string{}
}

func (ua *UnifiedRequestAnalyzer) getEntityKeywords(entityType string) []string {
	// Extract keys from product keywords map
	keywordMap := getProductKeywords()
	keywords := make([]string, 0, len(keywordMap))
	for keyword := range keywordMap {
		keywords = append(keywords, keyword)
	}
	return keywords // Simplified - in practice, would have entity-specific keywords
}

func (ua *UnifiedRequestAnalyzer) getDefaultModel(provider string) string {
	modelMap := map[string]string{
		"openai":    "gpt-4.1-mini",
		"anthropic": "claude-3-5-sonnet",
		"deepseek":  "deepseek-chat",
	}

	if model, exists := modelMap[provider]; exists {
		return model
	}
	return "default"
}

func (ua *UnifiedRequestAnalyzer) getProviderCost(provider string) float64 {
	costMap := map[string]float64{
		"openai":    0.020,
		"anthropic": 0.015,
		"deepseek":  0.001,
	}

	if cost, exists := costMap[provider]; exists {
		return cost
	}
	return 0.01
}

func (ua *UnifiedRequestAnalyzer) getProviderLatency(provider string) time.Duration {
	latencyMap := map[string]time.Duration{
		"openai":    time.Second * 2,
		"anthropic": time.Second * 3,
		"deepseek":  time.Second * 1,
	}

	if latency, exists := latencyMap[provider]; exists {
		return latency
	}
	return time.Second * 2
}

// Additional enhanced methods

func (ua *UnifiedRequestAnalyzer) generateActionsEnhanced(req EnhancedAnalysisRequest, entityType, intent string, params map[string]interface{}, complexityScore float64) []domain.ManagerAction {
	// Start with basic action generation
	actions := ua.generateActions(req.AnalysisRequest, entityType, intent, params)

	// Enhance based on complexity and context
	if complexityScore > 0.7 {
		// Add planning action for complex requests
		actions = append(actions, domain.ManagerAction{
			Type:        "planning",
			Description: "Breaking down complex request into manageable steps",
			Parameters:  map[string]interface{}{"complexity": complexityScore},
		})
	}

	// Add context-aware actions
	if req.SessionContext != nil && len(req.SessionContext.ActiveEntities) > 0 {
		actions = append(actions, domain.ManagerAction{
			Type:        "context_update",
			Description: "Updating session context with new information",
			Parameters:  map[string]interface{}{"entities": req.SessionContext.ActiveEntities},
		})
	}

	return actions
}

func (ua *UnifiedRequestAnalyzer) calculateQualityMetrics(req EnhancedAnalysisRequest, result *EnhancedAnalysisResult) *QualityMetrics {
	metrics := &QualityMetrics{
		IntentAccuracy:        0.8,
		EntityExtractionRate:  0.7,
		ParameterCompleteness: 0.6,
		ResponseRelevance:     0.8,
		UserSatisfaction:      0.7,
	}

	// Calculate composite score
	weights := ua.qualityScorer.scoringWeights
	metrics.CompositeScore = weights.IntentWeight*metrics.IntentAccuracy +
		weights.EntityWeight*metrics.EntityExtractionRate +
		weights.ParameterWeight*metrics.ParameterCompleteness +
		weights.RelevanceWeight*metrics.ResponseRelevance +
		weights.SatisfactionWeight*metrics.UserSatisfaction

	return metrics
}

func (ua *UnifiedRequestAnalyzer) calculateProcessingEfficiency(latency time.Duration, complexity float64) float64 {
	// Higher efficiency for faster processing of complex requests
	expectedLatency := time.Duration(complexity * float64(time.Second*5))
	if latency < expectedLatency {
		return 1.0 - (float64(latency)/float64(expectedLatency))*0.5
	}
	return 0.5
}

func (ua *UnifiedRequestAnalyzer) calculateResourceUtilization() map[string]float64 {
	return map[string]float64{
		"cpu":     0.3,
		"memory":  0.4,
		"io":      0.2,
		"network": 0.1,
	}
}

func (ua *UnifiedRequestAnalyzer) generateOptimizationSuggestions(result *EnhancedAnalysisResult) []OptimizationSuggestion {
	suggestions := []OptimizationSuggestion{}

	if result.ComplexityScore > 0.8 {
		suggestions = append(suggestions, OptimizationSuggestion{
			Category:     "complexity_reduction",
			Suggestion:   "Consider breaking down the request into simpler components",
			ExpectedGain: 0.3,
			Effort:       "medium",
			Priority:     "high",
		})
	}

	if result.PerformanceMetrics.AnalysisLatency > time.Second*3 {
		suggestions = append(suggestions, OptimizationSuggestion{
			Category:     "performance",
			Suggestion:   "Optimize processing pipeline for faster response times",
			ExpectedGain: 0.4,
			Effort:       "high",
			Priority:     "medium",
		})
	}

	return suggestions
}

func (ua *UnifiedRequestAnalyzer) calculateConfidenceEnhanced(message string, intent string, params map[string]interface{}, contextInsights *ContextInsights, qualityMetrics *QualityMetrics) float64 {
	// Start with basic confidence
	baseConfidence := ua.calculateConfidence(message, intent, params)

	// Enhance with context intelligence
	contextBonus := contextInsights.UserIntentClarity * 0.2

	// Quality metrics influence
	qualityBonus := qualityMetrics.CompositeScore * 0.15

	// Combine factors
	enhancedConfidence := baseConfidence + contextBonus + qualityBonus

	return math.Min(1.0, math.Max(0.0, enhancedConfidence))
}

func (ua *UnifiedRequestAnalyzer) generateReasoningEnhanced(intent, entityType string, paramCount int, complexityScore float64, contextInsights *ContextInsights) string {
	// Start with basic reasoning
	baseReasoning := ua.generateReasoning(intent, entityType, paramCount)

	// Add enhanced reasoning components
	complexityReasoning := ""
	if complexityScore > 0.7 {
		complexityReasoning = " Request complexity is high, requiring advanced processing."
	} else if complexityScore < 0.3 {
		complexityReasoning = " Request is straightforward with standard processing."
	}

	contextReasoning := ""
	if contextInsights.UserIntentClarity > 0.8 {
		contextReasoning = " User intent is very clear based on context analysis."
	} else if contextInsights.UserIntentClarity < 0.5 {
		contextReasoning = " User intent requires clarification based on context."
	}

	return baseReasoning + complexityReasoning + contextReasoning
}

func (ua *UnifiedRequestAnalyzer) trackAnalysisPerformance(requestID string, latency time.Duration, success bool, confidence float64) {
	// Track performance metrics for continuous improvement
	if metrics, exists := ua.performanceTracker.metrics[requestID]; exists {
		metrics.TotalRequests++
		if success {
			metrics.SuccessfulAnalysis++
		}
		metrics.AverageLatency = time.Duration((int64(metrics.AverageLatency) + int64(latency)) / 2)
		metrics.AccuracyRate = (metrics.AccuracyRate + confidence) / 2
		metrics.LastUpdated = time.Now()
	} else {
		ua.performanceTracker.metrics[requestID] = &AnalysisMetrics{
			TotalRequests: 1,
			SuccessfulAnalysis: func() int64 {
				if success {
					return 1
				}
				return 0
			}(),
			AverageLatency: latency,
			AccuracyRate:   confidence,
			ThroughputRate: 1.0,
			LastUpdated:    time.Now(),
		}
	}
}
