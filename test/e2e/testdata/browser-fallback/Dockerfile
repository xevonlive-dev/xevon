# Lightweight test image for browser fallback verification.
# Tests that xevon spidering finds system-installed chromium
# when the embedded browser binary is unavailable (e.g., ARM64 Linux).

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w" \
    -o /out/xevon ./cmd/xevon

# Runtime: Debian with system chromium (no embedded browser)
FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        chromium \
        ca-certificates \
        fonts-liberation \
        dumb-init \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/xevon /usr/local/bin/xevon

WORKDIR /root
ENTRYPOINT ["dumb-init", "--"]
