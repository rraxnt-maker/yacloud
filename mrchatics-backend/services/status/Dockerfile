FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o status-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/status-service .

EXPOSE 8081

CMD ["./status-service"]
