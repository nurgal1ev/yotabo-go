# ---------- BUILD STAGE ----------
FROM golang:1.25.1 AS builder

WORKDIR /app

# Сначала зависимости (для кеша)
COPY go.mod go.sum ./
RUN go mod download

# Потом весь проект
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app


# ---------- RUN STAGE ----------
FROM debian:bookworm-slim

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/app .

# Запуск
CMD ["./app"]