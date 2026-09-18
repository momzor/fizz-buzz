.PHONY: run build test test-cover lint tidy mocks docker-build docker-up docker-down swagger

run:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/fizzbuzz-api ./cmd/api

test:
	go test ./... -race

test-cover:
	go test ./internal/... -race -covermode=atomic -coverprofile=coverage.out
	go tool cover -func=coverage.out

lint:
	go vet ./...

tidy:
	go mod tidy

mocks:
	mockery --name StatsRepository --dir internal/usecase --output internal/usecase/mocks --outpkg mocks --with-expecter
	mockery --name FizzBuzzExecutor --dir internal/adapter/http --output internal/adapter/http/mocks --outpkg mocks --with-expecter
	mockery --name StatsQuery --dir internal/adapter/http --output internal/adapter/http/mocks --outpkg mocks --with-expecter

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs --pd

build:
	docker build -t fizzbuzz-api:local .

run:
	docker compose up --build

stop:
	docker compose down -v
