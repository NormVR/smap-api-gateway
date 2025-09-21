FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    build-base \
    musl-dev \
    pkgconfig

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 \
    CC=gcc \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-linkmode external -extldflags '-static' -w -s" \
    -tags musl \
    -trimpath \
    -o api-gateway cmd/api/main.go

FROM alpine:3.20
WORKDIR /root/

RUN apk add --no-cache \
    ca-certificates \
    tzdata

COPY --from=builder /app/api-gateway .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

EXPOSE 8080
CMD ["./api-gateway"]