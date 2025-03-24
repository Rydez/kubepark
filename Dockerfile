# Build stage
FROM golang:1.24.1-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/park ./park
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/carousel ./attractions/carousel
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/restroom ./attractions/restroom
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/guest ./guest

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -u 10001 kubepark

# Create internal directory for guest binary
RUN mkdir -p /opt/kubepark/internal && chown -R kubepark:kubepark /opt/kubepark

# Copy binaries from builder
COPY --from=builder /app/bin/park /usr/local/bin/
COPY --from=builder /app/bin/carousel /usr/local/bin/
COPY --from=builder /app/bin/restroom /usr/local/bin/
COPY --from=builder /app/bin/guest /opt/kubepark/internal/

# Set ownership of guest binary
RUN chown kubepark:kubepark /opt/kubepark/internal/guest

# Use non-root user
USER kubepark

# Set kubepark as the default command
CMD ["park"] 