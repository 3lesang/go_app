# 1. Build stage
FROM golang:1.25.3-alpine AS builder

# Install git and ca-certificates if needed
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build Go binary
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server cmd/server/main.go

# 2. Final image
FROM alpine:3.21

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary and docs
COPY --from=builder /app/server /app/server
COPY --from=builder /app/docs /app/docs

# Set default command
CMD ["/app/server"]
