DEFAULT_GOAL := ci

.PHONY: ci
ci: lint test

.PHONY: lint
lint:
	@golangci-lint run ./...

.PHONY: test
test:
	@go test -failfast -race -covermode=atomic ./...