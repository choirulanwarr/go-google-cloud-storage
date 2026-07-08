# Stage 1: Build the Go binary
FROM golang:1.23.7-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy dependency files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /app

# Install ca-certificates for HTTPS calls to Google Cloud Storage
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder stage
COPY --from=builder /app/server .

# Create logs directory
RUN mkdir -p /app/logs

# Expose the application port
EXPOSE 4000

# Run the binary
CMD ["./server"]
