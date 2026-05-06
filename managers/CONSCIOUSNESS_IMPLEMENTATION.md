# Consciousness Implementation Documentation

## Overview

The managers service has been transformed into a self-conscious, event-driven system that autonomously makes decisions based on platform events. This implementation integrates AI-powered decision making with rule-based automation.

## Architecture

### Core Components

1. **ConsciousnessManager** (`internal/consciousness/consciousness_manager.go`)
   - Orchestrates all consciousness components
   - Routes events through pattern detection → decision making → action execution
   - Integrates with existing application components

2. **DecisionMaker** (`internal/consciousness/decision_maker.go`)
   - Combines rule-based and AI-powered decision making
   - Loads rules from YAML configuration
   - Falls back to AI when no rules match

3. **AutonomousActionExecutor** (`internal/consciousness/autonomous_action_executor.go`)
   - Executes decisions using the application's tool system
   - Maps decisions to tool calls
   - Provides action verification for critical operations

4. **LearningEngine** (`internal/consciousness/learning_engine.go`)
   - Tracks decision outcomes
   - Calculates success rates
   - Provides feedback for continuous improvement

### Event Flow

```
Integration Event (NATS)
    ↓
Integration Event Handler
    ↓
ConsciousnessManager.ProcessEvent()
    ↓
Pattern Detection (existing)
    ↓
Decision Making (rules → AI)
    ↓
Action Execution (tools)
    ↓
Learning & Feedback
```

## Configuration

### Environment Variables

Key environment variables (see `.env.example`):

```bash
# Core Consciousness Settings
MANAGER_CONSCIOUSNESS_ENABLED=true
MANAGER_CONFIDENCE_THRESHOLD=0.8
MANAGER_MAX_ACTIONS_PER_MINUTE=10

# AI Provider Configuration
AI_PROVIDER_DEFAULT=deepseek
AI_PROVIDER_OPENAI_API_KEY=REDACTED
AI_PROVIDER_ANTHROPIC_API_KEY=REDACTED
AI_PROVIDER_DEEPSEEK_API_KEY=REDACTED

# Feature Flags
FEATURE_AI_DECISIONS=true
FEATURE_AUTONOMOUS_ACTIONS=true
FEATURE_LEARNING_MODE=true
```

### Decision Rules

Rules are defined in `internal/consciousness/rules/rules.yaml`:

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
```

## Implementation Details

### 1. Tool System Integration

The implementation reuses the working tool system from the assistants service:
- All tool files copied from assistants
- Import paths updated to managers package
- Tool execution remains identical

### 2. AI Integration

The system uses the shared AI infrastructure (`internal/ai`):
- Multiple provider support (OpenAI, Anthropic, DeepSeek)
- Automatic fallback between providers
- Cost-optimized default (DeepSeek)

### 3. Pattern Detection

Existing pattern detectors are used:
- ActivitySurgeDetector
- AbandonmentRiskDetector
- FraudRiskDetector
- SupportCrisisDetector
- UserSurgeDetector
- CancellationSpikeDetector

### 4. Decision Types

Supported autonomous actions:
- **send_notification** - Send email/SMS notifications
- **create_offer** - Create promotional offers
- **escalate_ticket** - Escalate support tickets
- **flag_order** - Flag orders for review

## Usage

### Enabling Consciousness

1. Set environment variable:
   ```bash
   export MANAGER_CONSCIOUSNESS_ENABLED=true
   ```

2. Start the service:
   ```bash
   make run-managers
   ```

3. The service will automatically:
   - Process incoming events
   - Detect patterns
   - Make decisions
   - Execute actions

### Monitoring

Monitor consciousness activity through logs:

```
INFO ConsciousnessManager processing event event_type=OrderCreated
INFO Pattern detected pattern_type=fraud_risk confidence=0.85
INFO Decision made decision_id=xxx action_count=1
INFO Action executed successfully action_type=flag_order
```

### Testing

1. **Dry Run Mode**: Set `FEATURE_DRY_RUN_MODE=true` to log actions without execution

2. **Manual Testing**: Send test events through NATS:
   ```bash
   nats pub orders.events '{"type":"OrderCreated","data":{...}}'
   ```

3. **Unit Tests**: Run consciousness tests:
   ```bash
   go test ./internal/consciousness/...
   ```

## Extending the System

### Adding New Rules

1. Edit `internal/consciousness/rules/rules.yaml`
2. Add new rule with pattern type, conditions, and actions
3. Restart service or wait for reload interval

### Adding New Actions

1. Implement tool in `internal/application/tools/`
2. Register tool in tool registry
3. Map action type to tool in `AutonomousActionExecutor`

### Custom Pattern Detectors

1. Implement the `PatternDetectorFunc` interface
2. Add to pattern detector in application initialization
3. Create corresponding decision rules

## Safety Measures

1. **Rate Limiting**: Max actions per minute configured
2. **Confidence Threshold**: Actions only taken above threshold
3. **Circuit Breaker**: Prevents cascading failures
4. **Audit Trail**: All decisions and actions logged
5. **Manual Override**: Can disable consciousness at runtime

## Troubleshooting

### Common Issues

1. **No actions being taken**
   - Check if consciousness is enabled
   - Verify pattern confidence meets threshold
   - Check decision rules configuration

2. **AI decisions failing**
   - Verify AI provider credentials
   - Check network connectivity
   - Review AI provider quotas

3. **Tools not executing**
   - Ensure tool registry is initialized
   - Check tool permissions
   - Verify tool parameters

### Debug Mode

Enable detailed logging:
```bash
export LOG_LEVEL=debug
```

## Performance Considerations

- Event processing timeout: 5s (configurable)
- Concurrent tool execution: 10 (configurable)
- Pattern detection uses goroutines for efficiency
- AI calls have circuit breaker protection

## Future Enhancements

1. **Machine Learning Integration**
   - Train custom models on decision outcomes
   - Improve pattern detection accuracy
   - Personalized decision making

2. **Advanced Rules Engine**
   - Complex condition evaluation
   - Time-based rules
   - Multi-action workflows

3. **Dashboard & Analytics**
   - Real-time consciousness metrics
   - Decision effectiveness tracking
   - Pattern visualization

## Conclusion

The consciousness implementation transforms the managers service into an intelligent, autonomous system that can:
- React to platform events in real-time
- Make informed decisions based on patterns
- Execute actions to optimize platform performance
- Learn from outcomes to improve over time

This creates a self-managing e-commerce platform that can handle many operational tasks without human intervention.