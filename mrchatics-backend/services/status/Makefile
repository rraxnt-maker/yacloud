.PHONY: build run test docker-build docker-run docker-stop clean

# Проверка версии Go
check-version:
	@go version | grep -q "go1.26.1" || (echo "Need Go 1.26.1" && exit 1)

# Сборка бинарника
build: check-version
	go build -o bin/status-service ./cmd/server

# Локальный запуск
run: check-version
	go run ./cmd/server

# Запуск тестов
test: check-version
	go test -v ./...

# Сборка Docker образа
docker-build:
	docker build -t status-service:latest .

# Запуск через Docker Compose
docker-run:
	docker-compose up --build

# Остановка Docker Compose
docker-stop:
	docker-compose down

# Очистка
clean:
	rm -rf bin/
	go clean -cache

# Показать логи
logs:
	docker-compose logs -f

# Перезапуск
restart:
	docker-compose restart

# Статус контейнеров
status:
	docker-compose ps

# Подключение к БД
db-shell:
	docker exec -it status_postgres_1 psql -U postgres -d status_service

# Инициализация проекта
init:
	go mod download
	go mod tidy

# Все команды
all: clean init build
