.PHONY: mocks

mocks:
	@mockery
	@go mod tidy