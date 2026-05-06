# -------------------------------
# 1) BUILD STAGE
# -------------------------------
FROM golang:1.23.4-alpine AS builder

ARG service
WORKDIR /middleman

# Copy go.mod, go.sum first to leverage Docker cache
COPY go.* ./
RUN go mod download

# Copy the rest of the code
COPY .. ./

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /service \
    ./${service}/cmd/service

# -------------------------------
# 2) RUNTIME STAGE
# -------------------------------
FROM alpine:3 AS runtime

# (Optional) Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the wait-for script
COPY --from=builder /middleman/docker/wait-for /usr/local/bin/wait-for
RUN chmod +x /usr/local/bin/wait-for

# Copy the service binary from the build stage
COPY --from=builder /service /middleman/service

# Switch to the non-root user (optional but recommended)
USER appuser

# If you prefer CMD:
CMD ["/middleman/service"]

# Or an ENTRYPOINT if you want to pass additional args:
# ENTRYPOINT ["/middleman/service"]
