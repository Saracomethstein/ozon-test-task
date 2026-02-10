.PHONY: mocks

mocks:
	@mockery
	@go mod tidy

build:
	@docker build -t ozon-test-task:latest .

run: build
	@docker run -p 8080:8080 ozon-test-task:latest