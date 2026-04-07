FROM golang:1.21-alpine AS builder

# Отключаем автоматическую загрузку тулчейна
ENV GOTOOLCHAIN=local

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходники
COPY . .

# Собираем
RUN CGO_ENABLED=0 go build -o user-profile-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/user-profile-service .
EXPOSE 8083
CMD ["./user-profile-service"]
