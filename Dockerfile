# Build stage
FROM golang:1.21-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o techircd ./cmd/techircd

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/techircd .

# Copy configuration
COPY --from=builder /app/configs/config.json ./config.json

# Create directory for logs
RUN mkdir -p /var/log/techircd

# Expose IRC ports
EXPOSE 6667 6697

# Run the binary
CMD ["./techircd", "-config", "config.json"]
