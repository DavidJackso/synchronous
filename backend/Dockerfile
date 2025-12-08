# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app ./cmd/app

# Final stage - minimal alpine image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    && update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/app .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Create configs directory (will be mounted as volume, but ensure it exists)
RUN mkdir -p ./configs

# Copy migrations and scripts
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts

# Copy configs if they exist (optional - config can be provided via volume or env vars)
# Note: configs directory will be mounted as volume in docker-compose

# Copy swagger.yaml if it exists
COPY --from=builder /app/swagger.yaml* ./

# Make scripts executable
RUN chmod +x /usr/local/bin/goose && \
    if [ -f ./scripts/migrate.sh ]; then chmod +x ./scripts/migrate.sh; fi

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

CMD ["./app"]

