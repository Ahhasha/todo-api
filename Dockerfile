# Этап 1: Сборка
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd

# Этап 2: Минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/service .

# Открываем порт
EXPOSE 8080

# Запускаем сервис
CMD ["./service"]