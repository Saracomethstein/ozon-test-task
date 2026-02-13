BINARY_NAME=build/ozon-test-task

.PHONY: all docker-build docker-up docker-down deps mocks build clean

all: docker-up

deps:
	@echo "==> Installing dependencies..."
	go mod tidy

build: clean deps
	@echo "==> Building the application..."
	mkdir build
	go build -o $(BINARY_NAME) cmd/service/main.go

clean:
	@echo "==> Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
	rm -rf build

docker-build:
	@echo "==> Building Docker containers..."
	docker compose build --no-cache

docker-up:
	@echo "==> Starting Docker containers..."
	docker compose up --build

docker-down:
	@echo "==> Stopping Docker containers..."
	docker compose down

mocks:
	@echo "==> Generating mocks..."
	@mockery
	@go mod tidy