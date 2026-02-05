GOLANGCI_LINT := $(HOME)/go/bin/golangci-lint

.PHONY: lint test run

lint:
	$(GOLANGCI_LINT) run

test:
	go test ./...

run:
	go run ./cmd/wallet-service/main.go
