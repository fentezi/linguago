all: lint static test migrate run

test:
	go test -v ./...

lint:
	golangci-lint run ./...	

static:
	staticcheck ./...

migrate:
	migrate -database 'postgresql://postgres:postgres@localhost:5432/dictionary' -path ./db/migrations up
run:
	go run cmd/translator/main.go --config='./config.yml'

.PHONY: lint static test run migrate
