all: test migrate run

test:
	go test -v ./...

migrate:
	migrate -database 'postgresql://postgres:postgres@localhost:5432/dictionary' -path ./db/migrations up
run:
	go run cmd/translator/main.go --config='./config.yml'

.PHONY: test run migrate
