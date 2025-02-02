# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o octopus-cache ./cmd/octopus-cache/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/octopus-cache .
COPY --from=builder /app/config.yaml ./config/

VOLUME /data
EXPOSE 8080
CMD ["./octopus-cache"]