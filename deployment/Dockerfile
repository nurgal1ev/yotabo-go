FROM golang:1.25.1 AS base

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc g++ libwebp-dev git ca-certificates \
    && rm -rf /var/lib/apt/lists/*

FROM base AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@v1.63.0

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /build/app ./cmd/app

FROM base AS runner_prod

WORKDIR /

COPY --from=builder /build/app /build/app

ENTRYPOINT ["/build/app"]

FROM base AS runner_dev

WORKDIR /app

COPY . .

COPY --from=builder /build/app /app
COPY --from=builder /go/bin/air /usr/local/bin/air

CMD ["air"]
