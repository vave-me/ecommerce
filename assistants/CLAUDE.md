# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

The assistants service is a production-ready microservice that provides AI-powered assistant functionality within the e-commerce platform. It manages:

- **AI Assistant Management**: Creating, activating, and configuring different types of assistants (user, admin, support, business, scheduler)
- **Conversation Management**: Creating and managing conversations between users and assistants
- **Multi-modal Input Processing**: Handling text, speech, image, and document inputs
- **Direct Repository Access**: Each repository is exposed as an LLM tool for maximum flexibility
- **Multi-Provider LLM Integration**: Supporting OpenAI, Anthropic, and DeepSeek with automatic fallback

## Service Structure

The assistants service follows the standard microservice pattern:
```
assistants/
├── cmd/service/main.go       # Entry point
├── internal/                 # Private business logic
│   ├── application/         # Commands, queries, and services
│   │   ├── commands/       # Command handlers
│   │   ├── queries/        # Query handlers
│   │   ├── processor/      # LLM and input processors
│   │   ├── services/       # Business services
│   │   └── tools/          # Tool integrations
│   ├── domain/             # Domain models and repositories
│   ├── grpc/               # gRPC server and client repositories
│   ├── handlers/           # Event handlers
│   └── postgres/           # PostgreSQL repositories
├── assistantspb/           # Protocol buffers
├── assistantsclient/       # Generated client code
└── migrations/             # Database migrations
```

## Common Development Commands

### Build Commands
```bash
# Generate protobuf and mock code
make generate

# Build the assistants service
make build-assistants

# Build without cache (for troubleshooting)
make build-assistants-no-cache

# Build all microservices (includes assistants)
make build-micro
```

### Running the Service
```bash
# Run all microservices (includes assistants)
make run-micro

# Stop all microservices
make down-micro
```

### Testing
```bash
# Run all tests for the assistants service
cd assistants && go test ./...

# Run tests with coverage
cd assistants && go test -coverprofile=coverage.out ./...

# Run a specific test
cd assistants && go test -run TestLLMProcessor ./internal/application/processor/...

# Run tests with verbose output
cd assistants && go test -v ./...
```

## Key Components (Updated Architecture)

### AI Client Management
The service uses an AIClientProvider that manages multiple AI clients with circuit breaker patterns and automatic fallback:
- Primary: OpenAI
- Fallback order: OpenAI → Anthropic → DeepSeek
- Circuit breaker configuration: 5 max failures, 2-minute reset timeout
- Automatic retry with exponential backoff

### Simplified Tool System
The service now uses a direct repository access pattern where each repository is exposed as an LLM tool:

**Tool Registry** (`internal/application/tools/tool_registry.go`)
- Direct repository method exposure as OpenAI-compliant tools
- Parallel tool execution support
- Clean parameter extraction and validation
- No complex abstractions or blocking code

**Available Repository Tools**:
- **Activity**: Like/dislike tracking, interaction analytics
- **Product**: Search, create, update, inventory management
- **Order**: Create, track, update status, history
- **User**: Profile management, preferences, search
- **Payment**: Process payments, refunds, history
- **Shipping**: Calculate rates, create labels, track
- **Basket**: Add/remove items, checkout preparation
- **Notification**: Send notifications, manage preferences
- **Support**: Ticket creation and management
- Plus 15+ other repositories, all directly accessible

### LLM Processor
Simplified processor (`internal/application/processor/llm_processor.go`):
- Direct tool execution without complex layers
- Parallel tool call support
- Conversation history management
- Streaming response capability
- Clean error handling

### Assistant Types
Different assistant types are configured in `config/assistant_types.yaml`:
- **User Assistant**: General user support and shopping assistance
- **Admin Assistant**: Administrative tasks and system management
- **Support Assistant**: Customer support and issue resolution
- **Business Assistant**: Business analytics and reporting
- **Scheduler Assistant**: Automated scheduling and task management

### Event-Driven Architecture
The service publishes and subscribes to domain events:
- Assistant events: Created, Activated, Deactivated, ConfigurationUpdated
- Conversation events: Created, MessageAdded, ContextUpdated, Archived
- Integration with NATS for async messaging

## Important Development Notes

### Production Readiness Checklist
- ✅ Direct repository access for LLM tools
- ✅ Parallel tool execution support
- ✅ Circuit breaker for AI providers
- ✅ Automatic fallback between providers
- ✅ Clean error handling and sanitization
- ✅ Metrics collection for tool execution
- ✅ Streaming support for long operations
- ✅ Event sourcing for audit trail

### Security Considerations
- All tool executions validate parameters
- Auth context propagation to downstream services
- Sensitive error information is sanitized
- Repository methods respect user permissions

### Performance Optimizations
- Parallel tool execution (up to 20 concurrent)
- Streaming responses for real-time feedback
- Connection pooling for gRPC clients
- Efficient event sourcing with snapshots
- Tool execution timeout (10 minutes default)

### Database Schema
The service uses PostgreSQL with event sourcing:
- Event store for assistant and conversation aggregates
- Read models for query optimization
- Outbox pattern for reliable event publishing
- Migrations run automatically on startup

### Testing Approach
- Unit tests for business logic
- Integration tests for repository interactions
- Mock generation using mockery for interfaces
- Test files alongside implementation (`*_test.go`)

### Common Patterns
- Command/Query separation (CQRS)
- Event sourcing for aggregates
- Repository pattern for data access
- Dependency injection using container
- Transactional outbox for event publishing

## Production Deployment

### Environment Variables
```bash
# AI Provider Configuration
OPENAI_API_KEY=REDACTED
ANTHROPIC_API_KEY=REDACTED
DEEPSEEK_API_KEY=REDACTED

# Service Configuration
NATS_URL=nats://localhost:4222
DATABASE_URL=postgres://user:pass@localhost/assistants
RPC_PORT=50051
HTTP_PORT=8080

# Feature Flags
ENABLE_STREAMING=true
MAX_CONCURRENT_TOOLS=20
TOOL_EXECUTION_TIMEOUT=600s
```

### Monitoring
- Prometheus metrics exposed on `/metrics`
- OpenTelemetry tracing enabled
- Tool execution metrics tracked
- AI provider usage monitored

### Scaling Considerations
- Service is stateless and horizontally scalable
- Event sourcing allows read model scaling
- gRPC connections are pooled and reused
- Circuit breakers prevent cascade failures

## Recent Refactoring (Production Ready)

The service was recently refactored to remove complexity and enable direct LLM access to all repositories:

1. **Removed DAEMON subsystem**: All daemon-related code was removed as it added unnecessary complexity
2. **Simplified Tool Registry**: Direct repository method exposure without abstractions
3. **Clean LLM Processor**: Removed blocking code and complex routing
4. **Parallel Execution**: Tools can now execute in parallel for better performance
5. **OpenAI Compliance**: All tools follow OpenAI function calling standards

This refactoring makes the service production-ready with:
- Clear, maintainable code
- Direct repository access for flexibility
- Excellent performance with parallel execution
- Robust error handling and monitoring
- Easy to extend with new repositories