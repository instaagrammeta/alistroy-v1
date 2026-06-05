# =========================================================
# AliStroy backend (Go) — multi-stage Dockerfile
# =========================================================

# ---- build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git tzdata

# Cache module downloads
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source
COPY . .

# Resolve any missing modules in case go.sum was empty
RUN go mod tidy

# Build a statically linked binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /out/alistroy ./cmd/api

# ---- runtime stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/alistroy /app/alistroy

# Pre-create the uploads directory (will be overlaid by docker volume)
RUN mkdir -p /app/uploads && chown -R app:app /app

ENV API_HOST=0.0.0.0 \
    API_PORT=8080 \
    API_ENV=production

EXPOSE 8080

USER app

CMD ["/app/alistroy"]
