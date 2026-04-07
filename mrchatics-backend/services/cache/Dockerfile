FROM golang:1.21-alpine

WORKDIR /app

COPY main.go .

RUN go build -o cache-service main.go

EXPOSE 8082

CMD ["./cache-service"]
