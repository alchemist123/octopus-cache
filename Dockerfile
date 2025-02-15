# Build stage
FROM golang:1.22.3-alpine AS builder
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git build-base

# Copy go modules separately for efficient caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o octopus-cache ./cmd/octopus-server

# Runtime stage
FROM alpine:latest
WORKDIR /app

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/octopus-cache .

# Create a volume mount point
VOLUME /data

# Expose the application port
EXPOSE 8080

# Set entrypoint
CMD ["/app/octopus-cache"]
