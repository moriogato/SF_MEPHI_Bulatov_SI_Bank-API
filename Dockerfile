# Этап 1: Сборка
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bank-api ./cmd/api/main.go

# Этап 2: Финальный образ
FROM alpine:latest

# Устанавливаем CA-сертификаты для HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем собранный бинарник из этапа сборки
COPY --from=builder /app/bank-api .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Копируем .env (опционально, лучше передавать через переменные окружения)
COPY --from=builder /app/.env .env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./bank-api"]
