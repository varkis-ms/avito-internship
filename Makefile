# Commands
env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.sample .env

build:
	go build -o ./cmd/app

run:
	go run ./cmd/app

compose-up:
	docker-compose -f docker-compose.yml up -d --remove-orphans

compose-down:
	docker-compose down --remove-orphans

unit-test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

linter:
	golangci-lint run

swagger:
	swag init -g internal/app/app.go --parseInternal --parseDependency

.PHONY: env build run compose-up compose-down unit-test cover linter swagger