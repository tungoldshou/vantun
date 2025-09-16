# VANTUN Docker Configuration

# Official VANTUN Docker image
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build VANTUN
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o vantun cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates iptables ip6tables

# Create non-root user
RUN addgroup -g 1000 -S vantun && \
    adduser -u 1000 -S vantun -G vantun

# Create necessary directories
RUN mkdir -p /etc/vantun /var/log/vantun /var/lib/vantun && \
    chown -R vantun:vantun /etc/vantun /var/log/vantun /var/lib/vantun

# Copy binary from builder
COPY --from=builder /build/vantun /usr/local/bin/vantun

# Copy entrypoint script
COPY scripts/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Switch to non-root user
USER vantun

# Expose default ports
EXPOSE 4242/tcp 4242/udp 8080/tcp

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 4242 || exit 1

# Set entrypoint
ENTRYPOINT ["docker-entrypoint.sh"]

# Default command
CMD ["server"]