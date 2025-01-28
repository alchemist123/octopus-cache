FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ttldb ./cmd/octopus-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/ttldb /ttldb
EXPOSE 8080
CMD ["/ttldb"]