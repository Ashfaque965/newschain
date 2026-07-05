

```dockerfile
# ==============================================================================
# NewsChain Dockerfile
# Multi-stage production build
# ==============================================================================

# ------------------------------------------------------------------------------
# Build Stage
# ------------------------------------------------------------------------------

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -o newschain \
    ./cmd/server

# ------------------------------------------------------------------------------
# Runtime Stage
# ------------------------------------------------------------------------------

FROM alpine:3.22

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata && \
    addgroup -S newschain && \
    adduser -S newschain -G newschain

WORKDIR /app

# Copy application binary
COPY --from=builder /app/newschain .

# Copy smart contracts (optional)
COPY contracts ./contracts

# Create writable data directory
RUN mkdir -p /app/data && \
    chown -R newschain:newschain /app

USER newschain

EXPOSE 8080

ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s \
            --timeout=5s \
            --start-period=10s \
            --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./newschain"]
```
