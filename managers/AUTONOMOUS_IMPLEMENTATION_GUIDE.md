# Autonomous Implementation Guide for Manager Service

This guide provides step-by-step instructions that can be executed without user interaction over 4-5 hours.

## Hour 1: Copy Core Components from Assistants

### Step 1.1: Backup Current Tools (5 minutes)
```bash
cd <workspace-root>/classified/managers
mkdir -p backup/tools
cp -r internal/application/tools/* backup/tools/
```

### Step 1.2: Copy Tool Implementation (10 minutes)
```bash
# Remove existing tools
rm -rf internal/application/tools/*

# Copy all tool files from assistants
cp -r ../assistants/internal/application/tools/* internal/application/tools/

# Update imports in all copied files
find internal/application/tools -name "*.go" -type f -exec sed -i 's/middleman\/assistants/middleman\/managers/g' {} +
```

### Step 1.3: Copy Processors (10 minutes)
```bash
# Backup existing processors
mkdir -p backup/processor
cp -r internal/application/processor/* backup/processor/

# Copy processors from assistants
cp -r ../assistants/internal/application/processor/* internal/application/processor/

# Update imports
find internal/application/processor -name "*.go" -type f -exec sed -i 's/middleman\/assistants/middleman\/managers/g' {} +
```

### Step 1.4: Copy Application Files (10 minutes)
```bash
# Copy specific application files
cp ../assistants/internal/application/simplified_executors.go internal/application/
cp ../assistants/internal/application/repository_translator.go internal/application/
cp ../assistants/internal/application/security_validator.go internal/application/

# Update imports
sed -i 's/middleman\/assistants/middleman\/managers/g' internal/application/simplified_executors.go
sed -i 's/middleman\/assistants/middleman\/managers/g' internal/application/repository_translator.go
sed -i 's/middleman\/assistants/middleman\/managers/g' internal/application/security_validator.go
```

### Step 1.5: Copy Service Interfaces (10 minutes)
```bash
# Copy service files
cp ../assistants/internal/application/services/ai_repository_language_service.go internal/application/services/
cp ../assistants/internal/application/services/interfaces.go internal/application/services/
cp ../assistants/internal/application/services/types.go internal/application/services/

# Update imports
find internal/application/services -name "*.go" -type f -exec sed -i 's/middleman\/assistants/middleman\/managers/g' {} +
```

### Step 1.6: Verify Compilation (15 minutes)
```bash
cd <workspace-root>/classified/managers
go mod tidy
go build ./...
```

## Hour 2: Create Consciousness Infrastructure

### Step 2.1: Create Consciousness Manager (20 minutes)
Create file: `internal/consciousness/consciousness_manager.go`
```go
package consciousness

import (
    "context"
    "fmt"
    "time"
    
    "github.com/rs/zerolog"
    "middleman/internal/ai"
    "middleman/internal/ddd"
    "middleman/managers/internal/application"
)

type ConsciousnessManager struct {
    app              application.App
    memoryCore       *MemoryCore
    patternDetector  *PatternDetector
    decisionMaker    *DecisionMaker
    actionExecutor   *AutonomousActionExecutor
    learningEngine   *LearningEngine
    aiManager        ai.AIClientManager
    logger           zerolog.Logger
}

func NewConsciousnessManager(
    app application.App,
    memoryCore *MemoryCore,
    patternDetector *PatternDetector,
    aiManager ai.AIClientManager,
    logger zerolog.Logger,
) *ConsciousnessManager {
    cm := &ConsciousnessManager{
        app:             app,
        memoryCore:      memoryCore,
        patternDetector: patternDetector,
        aiManager:       aiManager,
        logger:          logger,
    }
    
    // Initialize other components
    cm.decisionMaker = NewDecisionMaker(aiManager, logger)
    cm.actionExecutor = NewAutonomousActionExecutor(app, aiManager, logger)
    cm.learningEngine = NewLearningEngine(memoryCore, logger)
    
    return cm
}

func (cm *ConsciousnessManager) ProcessEvent(ctx context.Context, event ddd.Event) error {
    cm.logger.Info().
        Str("event_type", event.EventName()).
        Str("event_id", event.ID()).
        Msg("ConsciousnessManager processing event")
    
    // 1. Store in memory
    if err := cm.memoryCore.StoreEvent(ctx, event); err != nil {
        cm.logger.Error().Err(err).Msg("Failed to store event in memory")
    }
    
    // 2. Detect patterns
    pattern := cm.patternDetector.DetectPattern(ctx, event)
    if pattern == nil {
        return nil // No pattern detected, no action needed
    }
    
    cm.logger.Info().
        Str("pattern_type", pattern.Type).
        Float64("confidence", pattern.Confidence).
        Msg("Pattern detected")
    
    // 3. Make decision
    decision, err := cm.decisionMaker.MakeDecision(ctx, pattern)
    if err != nil {
        cm.logger.Error().Err(err).Msg("Failed to make decision")
        return err
    }
    
    if decision == nil {
        return nil // No decision made
    }
    
    // 4. Execute action
    if err := cm.actionExecutor.ExecuteDecision(ctx, decision); err != nil {
        cm.logger.Error().Err(err).Msg("Failed to execute decision")
        return err
    }
    
    // 5. Learn from outcome
    cm.learningEngine.RecordOutcome(ctx, decision, err)
    
    return nil
}
```

### Step 2.2: Create Decision Maker (20 minutes)
Create file: `internal/consciousness/decision_maker.go`
```go
package consciousness

import (
    "context"
    "fmt"
    
    "github.com/rs/zerolog"
    "middleman/internal/ai"
)

type DecisionMaker struct {
    aiManager ai.AIClientManager
    logger    zerolog.Logger
    rules     []DecisionRule
}

type DecisionRule struct {
    Name       string
    PatternType string
    Confidence float64
    Action     string
    Priority   int
}

func NewDecisionMaker(aiManager ai.AIClientManager, logger zerolog.Logger) *DecisionMaker {
    return &DecisionMaker{
        aiManager: aiManager,
        logger:    logger,
        rules:     loadDecisionRules(),
    }
}

func (dm *DecisionMaker) MakeDecision(ctx context.Context, pattern *Pattern) (*Decision, error) {
    // First try rule-based decision
    if decision := dm.checkRules(pattern); decision != nil {
        return decision, nil
    }
    
    // If no rule matches, use AI for decision
    return dm.analyzeWithAI(ctx, pattern)
}

func (dm *DecisionMaker) checkRules(pattern *Pattern) *Decision {
    for _, rule := range dm.rules {
        if rule.PatternType == pattern.Type && pattern.Confidence >= rule.Confidence {
            return &Decision{
                ID:       generateID(),
                Type:     rule.Action,
                Priority: fmt.Sprintf("%d", rule.Priority),
                Actions: []Action{
                    {
                        Type:       rule.Action,
                        Parameters: pattern.Data,
                    },
                },
            }
        }
    }
    return nil
}

func (dm *DecisionMaker) analyzeWithAI(ctx context.Context, pattern *Pattern) (*Decision, error) {
    client, err := dm.aiManager.GetDefaultClient()
    if err != nil {
        return nil, err
    }
    
    messages := []ai.Message{
        {
            Role: ai.RoleSystem,
            Content: "You are an autonomous e-commerce platform manager. Analyze patterns and suggest actions using available tools.",
        },
        {
            Role: ai.RoleUser,
            Content: fmt.Sprintf("Pattern detected: %s\nConfidence: %.2f\nData: %+v\n\nWhat action should be taken?",
                pattern.Type, pattern.Confidence, pattern.Data),
        },
    }
    
    tools := dm.getAvailableTools()
    
    response, err := client.ExecuteWithTools(ctx, messages, tools)
    if err != nil {
        return nil, err
    }
    
    return dm.parseAIResponse(response)
}

func loadDecisionRules() []DecisionRule {
    // TODO: Load from configuration file
    return []DecisionRule{
        {
            Name:       "abandoned_cart_recovery",
            PatternType: "cart_abandonment",
            Confidence: 0.7,
            Action:     "send_reminder_notification",
            Priority:   2,
        },
        {
            Name:       "fraud_alert",
            PatternType: "fraud_risk",
            Confidence: 0.8,
            Action:     "flag_for_review",
            Priority:   1,
        },
    }
}
```

### Step 2.3: Create Autonomous Action Executor (20 minutes)
Create file: `internal/consciousness/autonomous_action_executor.go`
```go
package consciousness

import (
    "context"
    "fmt"
    
    "github.com/rs/zerolog"
    "middleman/internal/ai"
    "middleman/managers/internal/application"
    "middleman/managers/internal/application/tools"
)

type AutonomousActionExecutor struct {
    app       application.App
    aiClient  ai.EnhancedAIService
    logger    zerolog.Logger
}

func NewAutonomousActionExecutor(app application.App, aiManager ai.AIClientManager, logger zerolog.Logger) *AutonomousActionExecutor {
    client, _ := aiManager.GetDefaultClient()
    return &AutonomousActionExecutor{
        app:      app,
        aiClient: client,
        logger:   logger,
    }
}

func (ae *AutonomousActionExecutor) ExecuteDecision(ctx context.Context, decision *Decision) error {
    ae.logger.Info().
        Str("decision_id", decision.ID).
        Str("decision_type", decision.Type).
        Msg("Executing autonomous decision")
    
    for _, action := range decision.Actions {
        if err := ae.executeAction(ctx, action); err != nil {
            ae.logger.Error().
                Err(err).
                Str("action_type", action.Type).
                Msg("Failed to execute action")
            return err
        }
    }
    
    return nil
}

func (ae *AutonomousActionExecutor) executeAction(ctx context.Context, action Action) error {
    // Map action types to tool calls
    toolCall := ae.mapActionToTool(action)
    
    execCtx := &tools.ToolExecutionContext{
        UserID: "system-consciousness",
        Role:   "manager",
    }
    
    results, err := ae.app.ExecuteTools(ctx, []ai.ToolCall{toolCall}, execCtx)
    if err != nil {
        return fmt.Errorf("tool execution failed: %w", err)
    }
    
    if len(results) > 0 && results[0].Error != nil {
        return fmt.Errorf("tool returned error: %w", results[0].Error)
    }
    
    ae.logger.Info().
        Str("action_type", action.Type).
        Msg("Action executed successfully")
    
    return nil
}

func (ae *AutonomousActionExecutor) mapActionToTool(action Action) ai.ToolCall {
    // Map internal action types to tool calls
    switch action.Type {
    case "send_reminder_notification":
        return ai.ToolCall{
            ID:   generateID(),
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionCall{
                Name:      "send_notification",
                Arguments: marshalParameters(action.Parameters),
            },
        }
    case "flag_for_review":
        return ai.ToolCall{
            ID:   generateID(),
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionCall{
                Name:      "flag_order",
                Arguments: marshalParameters(action.Parameters),
            },
        }
    default:
        return ai.ToolCall{
            ID:   generateID(),
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionCall{
                Name:      action.Type,
                Arguments: marshalParameters(action.Parameters),
            },
        }
    }
}
```

## Hour 3: Update Integration Points

### Step 3.1: Update Event Handler (15 minutes)
Update file: `internal/handlers/integration_events.go`

Add consciousness manager field:
```go
type integrationHandlers[T ddd.Event] struct {
    app                  application.App
    consciousnessManager *consciousness.ConsciousnessManager
}
```

Update HandleEvent method:
```go
func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
    // Route to consciousness manager
    if h.consciousnessManager != nil {
        return h.consciousnessManager.ProcessEvent(ctx, event)
    }
    
    // Fallback to direct processing
    return h.app.ProcessPlatformEvent(ctx, event)
}
```

### Step 3.2: Update Application Initialization (20 minutes)
Update file: `internal/application/application.go`

Add AI manager initialization:
```go
func New(
    // ... existing parameters
) *Application {
    // ... existing code
    
    // Initialize AI infrastructure
    factory := ai.NewClientFactory()
    
    // Register providers
    factory.RegisterProvider(ai.ProviderOpenAI, ai.ProviderConfig{
        APIKey:       os.Getenv("AI_PROVIDER_OPENAI_API_KEY"),
        DefaultModel: ai.ModelGPT4oMini,
        Enabled:      true,
    })
    
    factory.RegisterProvider(ai.ProviderDeepSeek, ai.ProviderConfig{
        APIKey:       os.Getenv("AI_PROVIDER_DEEPSEEK_API_KEY"),
        DefaultModel: ai.ModelDeepSeekV3,
        Enabled:      true,
    })
    
    aiManager := ai.NewClientManager(factory)
    aiManager.SetDefaultProvider(ai.ProviderDeepSeek)
    
    // Store AI manager in application
    app.aiManager = aiManager
    
    return app
}
```

### Step 3.3: Create Helper Functions (25 minutes)
Create file: `internal/consciousness/helpers.go`
```go
package consciousness

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "middleman/internal/ai"
)

func generateID() string {
    return uuid.New().String()
}

func marshalParameters(params map[string]interface{}) string {
    data, _ := json.Marshal(params)
    return string(data)
}

func (dm *DecisionMaker) getAvailableTools() []ai.ToolDefinition {
    return []ai.ToolDefinition{
        {
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionDef{
                Name:        "send_notification",
                Description: "Send a notification to a user",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "user_id":  map[string]string{"type": "string", "description": "User ID"},
                        "template": map[string]string{"type": "string", "description": "Notification template"},
                        "data":     map[string]string{"type": "object", "description": "Template data"},
                    },
                    "required": []string{"user_id", "template"},
                },
            },
        },
        {
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionDef{
                Name:        "create_offer",
                Description: "Create a promotional offer",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "product_id": map[string]string{"type": "string"},
                        "discount":   map[string]string{"type": "number"},
                        "expires_at": map[string]string{"type": "string"},
                    },
                    "required": []string{"product_id", "discount"},
                },
            },
        },
    }
}

func (dm *DecisionMaker) parseAIResponse(response *ai.CompletionResponse) (*Decision, error) {
    if len(response.Choices) == 0 {
        return nil, fmt.Errorf("no response from AI")
    }
    
    choice := response.Choices[0]
    if len(choice.Message.ToolCalls) == 0 {
        return nil, fmt.Errorf("no tool calls in AI response")
    }
    
    actions := make([]Action, 0, len(choice.Message.ToolCalls))
    for _, toolCall := range choice.Message.ToolCalls {
        var params map[string]interface{}
        if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
            continue
        }
        
        actions = append(actions, Action{
            Type:       toolCall.Function.Name,
            Parameters: params,
        })
    }
    
    return &Decision{
        ID:       generateID(),
        Type:     "ai_suggested",
        Priority: "medium",
        Actions:  actions,
    }, nil
}
```

## Hour 4: Create Configuration and Testing

### Step 4.1: Create Decision Rules Configuration (15 minutes)
Create directory and file: `internal/consciousness/rules/rules.yaml`
```yaml
decision_rules:
  - name: "abandoned_cart_recovery"
    pattern_type: "cart_abandonment"
    confidence_threshold: 0.7
    conditions:
      - field: "cart_value"
        operator: "gt"
        value: 50
    actions:
      - type: "send_notification"
        delay: "30m"
        parameters:
          template: "abandoned_cart_reminder"
          
  - name: "high_value_order_fraud_check"
    pattern_type: "high_value_order"
    confidence_threshold: 0.8
    conditions:
      - field: "order_amount"
        operator: "gt"
        value: 1000
    actions:
      - type: "flag_for_review"
        immediate: true
        parameters:
          reason: "high_value_order"
          priority: "high"
          
  - name: "support_ticket_escalation"
    pattern_type: "urgent_support"
    confidence_threshold: 0.9
    actions:
      - type: "escalate_ticket"
        immediate: true
        parameters:
          escalation_level: 2
          
  - name: "inventory_alert"
    pattern_type: "low_stock"
    confidence_threshold: 0.85
    conditions:
      - field: "stock_level"
        operator: "lt"
        value: 10
    actions:
      - type: "send_notification"
        parameters:
          template: "low_stock_alert"
          recipient_role: "seller"
```

### Step 4.2: Create Learning Engine (20 minutes)
Create file: `internal/consciousness/learning_engine.go`
```go
package consciousness

import (
    "context"
    "time"
    
    "github.com/rs/zerolog"
)

type LearningEngine struct {
    memoryCore *MemoryCore
    logger     zerolog.Logger
    outcomes   []OutcomeRecord
}

type OutcomeRecord struct {
    DecisionID string
    Decision   *Decision
    Success    bool
    Error      error
    Timestamp  time.Time
}

func NewLearningEngine(memoryCore *MemoryCore, logger zerolog.Logger) *LearningEngine {
    return &LearningEngine{
        memoryCore: memoryCore,
        logger:     logger,
        outcomes:   make([]OutcomeRecord, 0),
    }
}

func (le *LearningEngine) RecordOutcome(ctx context.Context, decision *Decision, err error) {
    outcome := OutcomeRecord{
        DecisionID: decision.ID,
        Decision:   decision,
        Success:    err == nil,
        Error:      err,
        Timestamp:  time.Now(),
    }
    
    le.outcomes = append(le.outcomes, outcome)
    
    // Analyze patterns in outcomes
    le.analyzeOutcomes()
}

func (le *LearningEngine) analyzeOutcomes() {
    // Calculate success rates
    if len(le.outcomes) < 10 {
        return
    }
    
    successCount := 0
    for _, outcome := range le.outcomes[len(le.outcomes)-10:] {
        if outcome.Success {
            successCount++
        }
    }
    
    successRate := float64(successCount) / 10.0
    
    if successRate < 0.5 {
        le.logger.Warn().
            Float64("success_rate", successRate).
            Msg("Low success rate detected in recent decisions")
    }
}

func (le *LearningEngine) GetSuccessRate(decisionType string) float64 {
    total := 0
    successful := 0
    
    for _, outcome := range le.outcomes {
        if outcome.Decision.Type == decisionType {
            total++
            if outcome.Success {
                successful++
            }
        }
    }
    
    if total == 0 {
        return 0.0
    }
    
    return float64(successful) / float64(total)
}
```

### Step 4.3: Create Module Integration (25 minutes)
Update file: `module.go`

Add consciousness initialization:
```go
func initConsciousness(app application.App, aiManager ai.AIClientManager) *consciousness.ConsciousnessManager {
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
    
    // Use existing components
    memoryCore := app.GetMemoryCore()
    patternDetector := app.GetPatternDetector()
    
    // Create consciousness manager
    consciousnessManager := consciousness.NewConsciousnessManager(
        app,
        memoryCore,
        patternDetector,
        aiManager,
        logger,
    )
    
    return consciousnessManager
}
```

## Hour 5: Testing and Verification

### Step 5.1: Create Unit Tests (30 minutes)
Create test files for each component:
- `internal/consciousness/consciousness_manager_test.go`
- `internal/consciousness/decision_maker_test.go`
- `internal/consciousness/autonomous_action_executor_test.go`

### Step 5.2: Integration Testing (20 minutes)
Create file: `internal/consciousness/integration_test.go`
```go
// +build integration

package consciousness_test

import (
    "context"
    "testing"
    
    "middleman/internal/ddd"
    "middleman/managers/internal/consciousness"
)

func TestEventProcessingFlow(t *testing.T) {
    // Test complete flow from event to action
}
```

### Step 5.3: Final Verification (10 minutes)
```bash
# Run all tests
cd <workspace-root>/classified/managers
go test ./...

# Build the service
make build-managers

# Check for any compilation errors
go vet ./...
```

## Environment Setup

Create `.env` file for managers service:
```bash
# Consciousness settings
MANAGER_CONSCIOUSNESS_ENABLED=true
MANAGER_DECISION_DELAY_MS=1000
MANAGER_CONFIDENCE_THRESHOLD=0.8
MANAGER_MAX_ACTIONS_PER_MINUTE=10

# AI Provider settings
AI_PROVIDER_DEFAULT=deepseek
AI_PROVIDER_OPENAI_API_KEY=REDACTED
AI_PROVIDER_ANTHROPIC_API_KEY=REDACTED
AI_PROVIDER_DEEPSEEK_API_KEY=REDACTED
AI_PROVIDER_FALLBACK_ENABLED=true

# Tool settings
MANAGER_TOOL_TIMEOUT=30s
MANAGER_TOOL_RETRY_COUNT=3
```

## Completion Checklist

- [ ] All tool files copied from assistants
- [ ] All imports updated to managers package
- [ ] Consciousness manager created
- [ ] Decision maker implemented
- [ ] Autonomous action executor created
- [ ] Learning engine implemented
- [ ] Event handler updated
- [ ] AI infrastructure integrated
- [ ] Decision rules configured
- [ ] Tests written
- [ ] Service compiles successfully
- [ ] Environment variables configured

## Notes for Autonomous Execution

1. Each step includes specific commands that can be run
2. File contents are provided in full
3. Import updates use sed commands for automation
4. Test steps verify each phase
5. Rollback is possible using backup directories

This guide can be executed step-by-step over 4-5 hours without user interaction.