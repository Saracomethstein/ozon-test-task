BINARY_NAME=build/ozon-test-task

.PHONY: all docker-build docker-up docker-down deps mocks build clean

all: docker-up

deps:
	@echo "==> Installing dependencies..."
	go mod tidy

build: deps
	@echo "==> Building the application..."
	mkdir build
	go build -o $(BINARY_NAME) cmd/service/main.go

clean:
	@echo "==> Cleaning up..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -rf build
	@rm coverage.*

test:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@open ./coverage.html    

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