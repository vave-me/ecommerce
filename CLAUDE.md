# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based e-commerce platform written in Go. The main services are:

- **cosec**: Security and authentication service
- **payments**: Payment processing and Stripe integration
- **ordering**: Order management and fulfillment
- **baskets**: Shopping cart functionality
- **shipping**: Shipping and logistics management

Additional services include: users, search, products, categories, notifications, messages, comments, activity, offers, wishlists, support, newsletters, media, and more.

## Service Structure

Each service follows this consistent pattern:
```
service/
├── cmd/service/main.go       # Entry point
├── internal/                 # Private business logic
│   ├── application/         # Commands and queries
│   ├── domain/             # Domain models
│   ├── handlers/           # Event/command handlers
│   └── grpc/               # gRPC server
├── migrations/              # Database migrations
├── module.go               # Service configuration
├── [service]pb/            # Protocol buffers
└── [service]client/        # Generated clients
```

## Common Development Commands

### Build Commands
```bash
# Install required tools
make install-tools

# Generate code (protobuf, mocks, etc.)
make generate

# Build specific service
make build-cosec
make build-payments
make build-ordering
make build-baskets
make build-shipping

# Build all microservices
make build-micro
```

### Running Services
```bash
# Run all microservices
make run-micro

# Run specific service groups
make run-payments
make run-frontend

# Stop services
make down-micro
```

### Testing
```bash
# Run tests for a specific service
cd cosec && go test ./...
cd payments && go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run specific test
go test -run TestName ./internal/...

# Run integration tests (if available)
go test -tags=integration ./...
```

### Database Migrations
Each service manages its own PostgreSQL schema. Migrations are in the `migrations/` directory and run automatically on service startup.

## Key Technologies

- **gRPC**: Primary service communication protocol
- **NATS**: Message broker for async events
- **PostgreSQL**: Primary database (each service has its own schema)
- **Protocol Buffers**: API definitions in `*pb/` directories
- **OpenTelemetry**: Distributed tracing and metrics
- **Docker Compose**: Local development environment

## Code Generation

The project uses several code generation tools:
- **buf**: Protocol buffer management
- **mockery**: Mock generation for testing
- **go-swagger**: OpenAPI spec generation

Run `make generate` after modifying:
- `.proto` files
- Interface definitions (for mocks)
- OpenAPI specifications

## Service Communication

Services communicate via:
1. **gRPC**: Direct synchronous calls
2. **NATS**: Asynchronous events (e.g., OrderCreated, PaymentProcessed)
3. **REST Gateway**: External API access through gRPC-gateway

## Testing Approach

- Unit tests: Alongside code files (`*_test.go`)
- Integration tests: Test database interactions
- Contract tests: Using Pact for service contracts
- Mock generation: Use `mockery` for interface mocks

## Important Notes

- Each service has public methods defined in `cmd/service/main.go` for auth bypass
- Services use separate PostgreSQL schemas, not separate databases
- Docker builds use a shared `Dockerfile.microservices` with service name as build arg
- The project uses a monorepo structure with shared dependencies in the root `go.mod`

## Manager Service Autonomous Implementation

The managers service is being transformed into a self-conscious, event-driven service. Follow these steps:

### Implementation Instructions

1. **Copy Working Components from Assistants** (CRITICAL - DO THIS FIRST):
   ```bash
   # Copy all tool implementations
   rm -rf managers/internal/application/tools/*
   cp -r assistants/internal/application/tools/* managers/internal/application/tools/
   
   # Copy processors
   cp -r assistants/internal/application/processor/* managers/internal/application/processor/
   
   # Copy key application files
   cp assistants/internal/application/simplified_executors.go managers/internal/application/
   cp assistants/internal/application/repository_translator.go managers/internal/application/
   cp assistants/internal/application/security_validator.go managers/internal/application/
   
   # Copy service interfaces
   cp assistants/internal/application/services/ai_repository_language_service.go managers/internal/application/services/
   cp assistants/internal/application/services/interfaces.go managers/internal/application/services/
   cp assistants/internal/application/services/types.go managers/internal/application/services/
   ```

2. **Configure AI Infrastructure**: The managers service uses the shared AI infrastructure from `/internal/ai`:
   - Supports multiple providers (OpenAI, Anthropic, DeepSeek)
   - Automatic fallback between providers
   - Cost-effective provider selection
   - Used for intelligent pattern analysis and decision making

3. **Update Imports**: After copying, update all imports from `middleman/assistants` to `middleman/managers`

4. **Create Consciousness Manager**:
   - Create `/internal/consciousness/consciousness_manager.go`
   - Integrate with existing consciousness components (MemoryCore, PatternDetector, etc.)
   - Route events through consciousness for autonomous decisions

5. **Create Autonomous Action Executor**:
   - Create `/internal/consciousness/autonomous_action_executor.go`
   - Use copied tool system for action execution
   - Integrate with AI client for intelligent decisions
   - Implement decision-to-action mapping

6. **Define Decision Rules**:
   - Create `/internal/consciousness/rules/rules.yaml`
   - Define triggers, conditions, and actions
   - Map actions to tool executions

7. **Update Event Handler**:
   - Modify `/internal/handlers/integration_events.go`
   - Route events to ConsciousnessManager instead of direct processing

### Key Implementation Pattern

```go
// Event flow:
Integration Event → ConsciousnessManager → Pattern Detection → Decision Making → Tool Execution

// Example autonomous action:
func (ae *AutonomousActionExecutor) executeAction(ctx context.Context, decision Decision) error {
    // Use the copied tool system from assistants
    toolCall := externalai.ToolCall{
        Name: decision.ToolName,
        Arguments: decision.Parameters,
    }
    
    results, err := ae.app.ExecuteTools(ctx, []externalai.ToolCall{toolCall}, 
        &tools.ToolExecutionContext{
            UserID: "system-consciousness",
            Role:   "manager",
        })
    
    return err
}
```

### Testing the Implementation

1. **Verify Tool Copying**: Ensure all tools work exactly as in assistants
2. **Test Event Processing**: Send test events and verify consciousness processing
3. **Validate Decisions**: Check that decisions are made based on rules
4. **Confirm Actions**: Verify autonomous actions execute correctly

### Environment Configuration

```bash
# Add to managers service environment
MANAGER_CONSCIOUSNESS_ENABLED=true
MANAGER_DECISION_DELAY_MS=1000
MANAGER_CONFIDENCE_THRESHOLD=0.8
MANAGER_MAX_ACTIONS_PER_MINUTE=10

# AI Provider configuration
AI_PROVIDER_DEFAULT=deepseek
AI_PROVIDER_OPENAI_API_KEY=REDACTED
AI_PROVIDER_ANTHROPIC_API_KEY=REDACTED
AI_PROVIDER_DEEPSEEK_API_KEY=REDACTED
AI_PROVIDER_FALLBACK_ENABLED=true
```

### Reference Documentation

See `<workspace-root>/classified/managers/MANAGER_AUTONOMY_PLAN.md` for detailed implementation plan.