GOLANGCI_LINT := $(HOME)/go/bin/golangci-lint

.PHONY: lint test run up down logs

lint:
	$(GOLANGCI_LINT) run

test:
	go test ./...

run:
	go run ./cmd/wallet-service/main.go

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f postgres