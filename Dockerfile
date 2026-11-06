# ===========================
# Stage 1 - Build
# ===========================
FROM golang:1.25-bookworm AS builder

WORKDIR /src

# Зависимости кешируются отдельно от исходников.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/daily-team-bot \
    ./cmd/bot

# ===========================
# Stage 2 - Runtime
# ===========================
FROM debian:bookworm-slim

WORKDIR /app

# Сертификаты нужны Telegram API, tzdata — для часовых поясов IANA.
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/daily-team-bot /app/daily-team-bot

ENV TZ=Europe/Moscow

ENTRYPOINT ["/app/daily-team-bot"]
