# Users Service Protocol Buffers

This directory contains the Protocol Buffer definitions for the Users service.

## Contents

- `api.proto` - Service API definitions
- `events.proto` - Event definitions for domain events
- Generated code files (*.pb.go, *.pb.gw.go)

## Recent Changes

Recently added token management functionality:
- Enhanced token refresh with token rotation
- Explicit token invalidation
- Improved token security and tracking

## Generating Code

After making changes to `.proto` files, you'll need to regenerate the Go code:

```bash
# Install protoc and required plugins if not already installed
# For Ubuntu/Debian:
# apt-get install -y protobuf-compiler
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate code
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    userspb/api.proto userspb/events.proto

# If you have gateway code generation:
protoc --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
    userspb/api.proto
```

## Implementation Notes

### Token Management

The token management system now supports:

1. **Token Rotation** - When refresh tokens are used, they are invalidated and a new refresh token is issued
2. **Token Invalidation** - Explicit invalidation of tokens for security purposes
3. **Token Tracking** - Better tracking and auditing of token lifecycle

### Domain Events

New domain events have been added:
- `UserTokenInvalidated` - When tokens are explicitly invalidated
- `UserTokenRefreshed` - When tokens are refreshed/rotated

These events can be used by event handlers to implement additional security features, such as token blacklisting or security notifications. 