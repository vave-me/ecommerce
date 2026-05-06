# Manager Service Autonomy Plan

## Overview
This document outlines the transformation of the managers service from a reactive system to a self-conscious, event-driven service that can autonomously make decisions based on platform events. The implementation will directly copy working patterns from the assistants service.

## Implementation Strategy

### Core Principle: Copy from Assistants
All tool execution, processors, and application patterns will be directly copied from the assistants service where they are known to be working correctly.

### Use Shared AI Infrastructure
The managers service will leverage the shared AI infrastructure in `/internal/ai` which provides:
- Multiple AI provider support (OpenAI, Anthropic, DeepSeek)
- Unified interfaces for completions, tools, and structured output
- Client factory and manager for provider selection
- Multimodal capabilities (vision, audio, image generation)
- Automatic fallback and best provider selection

## Files to Copy from Assistants to Managers

### 1. Application Layer Structure
**Copy these files directly from assistants to managers:**
- `/internal/application/application.go` → Update to add consciousness components
- `/internal/application/simplified_executors.go` → Copy as-is
- `/internal/application/repository_translator.go` → Copy as-is
- `/internal/application/security_validator.go` → Copy as-is
- `/internal/application/performance_optimizer.go` → Copy as-is
- `/internal/application/config.go` → Copy as-is

### 2. Tool Implementation
**Remove all existing tool files in managers and copy from assistants:**
```bash
# Remove managers tool files
rm -rf <workspace-root>/classified/managers/internal/application/tools/*

# Copy all tool files from assistants
cp -r <workspace-root>/classified/assistants/internal/application/tools/* \
      <workspace-root>/classified/managers/internal/application/tools/
```

**Key files to ensure are copied:**
- `tool_executor.go`
- `tool_registry.go`
- `tool_handlers.go`
- `tool_definitions.go`
- `tool_types.go`
- `type_converter.go`
- `openai_tool_converter.go`
- `comprehensive_tool_registry.go`
- All `*_tools.go` files
- All `handlers_*.go` files

### 3. Processors
**Copy all processor implementations:**
```bash
# Copy processors
cp -r <workspace-root>/classified/assistants/internal/application/processor/* \
      <workspace-root>/classified/managers/internal/application/processor/
```

### 4. Services
**Copy service implementations:**
```bash
# Copy services
cp <workspace-root>/classified/assistants/internal/application/services/ai_repository_language_service.go \
   <workspace-root>/classified/managers/internal/application/services/
   
cp <workspace-root>/classified/assistants/internal/application/services/interfaces.go \
   <workspace-root>/classified/managers/internal/application/services/
   
cp <workspace-root>/classified/assistants/internal/application/services/types.go \
   <workspace-root>/classified/managers/internal/application/services/
```

## Enhanced Architecture for Self-Consciousness

### 1. Event Processing Flow
```
Integration Event 
    ↓
ConsciousnessManager (NEW)
    ↓
Pattern Detection → Decision Making
    ↓
AutonomousActionExecutor (NEW)
    ↓
Tool Execution (COPIED FROM ASSISTANTS)
```

### 2. New Components to Add

#### A. ConsciousnessManager
```go
// /internal/consciousness/consciousness_manager.go
type ConsciousnessManager struct {
    app              application.App  // Uses copied application structure
    memoryCore       *MemoryCore
    patternDetector  *PatternDetector
    decisionMaker    *DecisionMaker
    actionExecutor   *AutonomousActionExecutor
    learningEngine   *LearningEngine
    aiManager        ai.AIClientManager  // Shared AI infrastructure
}

func (cm *ConsciousnessManager) ProcessEvent(ctx context.Context, event ddd.Event) error {
    // 1. Store in memory
    cm.memoryCore.StoreEvent(ctx, event)
    
    // 2. Detect patterns
    patterns := cm.patternDetector.DetectPatterns(ctx, event)
    
    // 3. Make decisions
    decisions := cm.decisionMaker.MakeDecisions(ctx, patterns)
    
    // 4. Execute actions using copied tool system
    for _, decision := range decisions {
        cm.actionExecutor.ExecuteDecision(ctx, decision)
    }
    
    // 5. Learn from outcomes
    cm.learningEngine.RecordOutcome(ctx, decisions)
    
    return nil
}
```

#### B. AutonomousActionExecutor
```go
// /internal/consciousness/autonomous_action_executor.go
type AutonomousActionExecutor struct {
    app          application.App  // Uses copied application interface
    toolExecutor tools.ToolExecutor // Uses copied tool executor
    aiClient     ai.EnhancedAIService // For intelligent decision making
}

func (ae *AutonomousActionExecutor) ExecuteDecision(ctx context.Context, decision Decision) error {
    switch decision.Type {
    case "send_notification":
        return ae.executeNotification(ctx, decision)
    case "create_offer":
        return ae.executeCreateOffer(ctx, decision)
    case "escalate_ticket":
        return ae.executeEscalateTicket(ctx, decision)
    // ... more action types
    }
}

// Uses the copied tool system
func (ae *AutonomousActionExecutor) executeNotification(ctx context.Context, decision Decision) error {
    toolCall := externalai.ToolCall{
        Name: "send_notification",
        Arguments: decision.Parameters,
    }
    
    results, err := ae.app.ExecuteTools(ctx, []externalai.ToolCall{toolCall}, &tools.ToolExecutionContext{
        UserID: "system",
        Role:   "consciousness",
    })
    
    return err
}
```

### 3. Integration with Existing Event Handler

**Update `/internal/handlers/integration_events.go`:**
```go
func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
    // Instead of calling app.ProcessPlatformEvent
    // Route to consciousness manager
    return h.consciousnessManager.ProcessEvent(ctx, event)
}
```

### 4. AI-Powered Decision Making

The consciousness system will use the shared AI infrastructure for intelligent decision making:

```go
// Example: Using AI to analyze patterns and make decisions
func (dm *DecisionMaker) analyzeWithAI(ctx context.Context, pattern Pattern) (*Decision, error) {
    // Create prompt for AI analysis
    messages := []ai.Message{
        {
            Role: ai.RoleSystem,
            Content: "You are an autonomous e-commerce platform manager. Analyze patterns and suggest actions.",
        },
        {
            Role: ai.RoleUser,
            Content: fmt.Sprintf("Pattern detected: %s\nConfidence: %.2f\nData: %+v\n\nWhat action should be taken?",
                pattern.Type, pattern.Confidence, pattern.Data),
        },
    }
    
    // Define tools the AI can suggest
    tools := []ai.Tool{
        {
            Type: ai.ToolTypeFunction,
            Function: ai.FunctionDef{
                Name: "send_notification",
                Description: "Send a notification to a user",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "user_id": map[string]string{"type": "string"},
                        "template": map[string]string{"type": "string"},
                    },
                },
            },
        },
        // Add more tools as needed
    }
    
    // Get AI recommendation
    response, err := dm.aiClient.ExecuteWithTools(ctx, messages, tools)
    if err != nil {
        return nil, err
    }
    
    // Convert AI response to Decision
    return dm.parseAIResponse(response)
}
```

### 4. Decision Rules Configuration

Create `/internal/consciousness/rules/rules.yaml`:
```yaml
decision_rules:
  - name: "abandoned_cart_recovery"
    trigger:
      event_type: "BasketAbandonedEvent"
      conditions:
        - field: "total_value"
          operator: "gt"
          value: 50
    actions:
      - type: "send_notification"
        tool: "notification_tools.SendNotification"
        delay: "30m"
        parameters:
          template: "abandoned_cart_reminder"
          
  - name: "fraud_detection"
    trigger:
      pattern: "suspicious_payment_pattern"
      confidence: 0.8
    actions:
      - type: "flag_order"
        tool: "order_tools.FlagOrder"
        immediate: true
        
  - name: "support_escalation"
    trigger:
      event_type: "TicketCreatedEvent"
      conditions:
        - field: "priority"
          operator: "eq"
          value: "urgent"
    actions:
      - type: "escalate_ticket"
        tool: "support_tools.EscalateTicket"
        immediate: true
```

## Implementation Steps

### Phase 1: Copy Core Components (Day 1-2)
1. Copy all tool-related files from assistants
2. Copy processors and services
3. Copy application structure
4. Update imports to use managers package

### Phase 2: Add Consciousness Layer (Day 3-4)
1. Create ConsciousnessManager
2. Create AutonomousActionExecutor
3. Create DecisionMaker with rules engine
4. Update event handlers to use consciousness

### Phase 3: Configure Decision Rules (Day 5)
1. Create rule configuration system
2. Define initial rule set
3. Add rule validation
4. Create rule testing framework

### Phase 4: Testing & Validation (Day 6-7)
1. Test copied tool execution
2. Test consciousness decisions
3. Test end-to-end flows
4. Performance testing

## Key Differences from Original Plan

1. **Direct Copy Strategy**: Instead of reimplementing, copy working code from assistants
2. **Minimal Modifications**: Only add consciousness layer on top
3. **Proven Foundation**: Use battle-tested tool execution from assistants
4. **Faster Implementation**: Copying reduces development time
5. **Consistent Behavior**: Ensures managers tools work exactly like assistants

## Configuration Updates

### Environment Variables
```bash
# Consciousness settings
MANAGER_CONSCIOUSNESS_ENABLED=true
MANAGER_DECISION_DELAY_MS=1000
MANAGER_CONFIDENCE_THRESHOLD=0.8
MANAGER_MAX_ACTIONS_PER_MINUTE=10

# Use same tool settings as assistants
MANAGER_TOOL_TIMEOUT=30s
MANAGER_TOOL_RETRY_COUNT=3

# AI Provider settings (using shared infrastructure)
AI_PROVIDER_DEFAULT=deepseek  # Cost-effective for autonomous operations
AI_PROVIDER_OPENAI_API_KEY=REDACTED
AI_PROVIDER_ANTHROPIC_API_KEY=REDACTED
AI_PROVIDER_DEEPSEEK_API_KEY=REDACTED
AI_PROVIDER_FALLBACK_ENABLED=true
```

### AI Integration Configuration
```go
// Initialize AI manager in application
func NewApplication(...) *Application {
    // Create AI client manager
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
    aiManager.SetDefaultProvider(ai.ProviderDeepSeek) // Cost-effective
    
    // ... rest of initialization
}
```

## Success Criteria

1. All assistants tools working in managers
2. Autonomous decisions based on events
3. Successful action execution via tools
4. No regression in existing functionality
5. Performance within acceptable limits

## Safety Measures

1. Feature flag for consciousness
2. Dry-run mode for testing
3. Action rate limiting
4. Audit trail for all decisions
5. Manual override capability