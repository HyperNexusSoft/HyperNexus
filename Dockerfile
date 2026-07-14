# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source code
COPY go/ .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -buildvcs=false -o /tormentnexus ./cmd/tormentnexus

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

# Copy binary
COPY --from=builder /tormentnexus .

# Create data directory
RUN mkdir -p /root/.tormentnexus/memory

# Expose ports
EXPOSE 7778 7779

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:7778/health || exit 1

# Run
CMD ["./tormentnexus", "serve"]
